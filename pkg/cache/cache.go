package cache

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrExpired   = errors.New("cache entry expired")
)

type Entry struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	ExpireAt   time.Time   `json:"expire_at"`
	AccessedAt time.Time   `json:"accessed_at"`
	Size       int64       `json:"size"`
}

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Stats() CacheStats
}

type CacheStats struct {
	Hits      int64     `json:"hits"`
	Misses    int64     `json:"misses"`
	Evictions int64     `json:"evictions"`
	Size      int64     `json:"size"`
	Items     int       `json:"items"`
	LastClean time.Time `json:"last_clean"`
}

// LRUCache implements a Least Recently Used cache with TTL support
type LRUCache struct {
	maxSize         int64
	maxItems        int
	size            int64
	items           map[string]*list.Element
	evictList       *list.List
	mu              sync.RWMutex
	stats           CacheStats
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

func NewLRUCache(maxSize int64, maxItems int) *LRUCache {
	cache := &LRUCache{
		maxSize:         maxSize,
		maxItems:        maxItems,
		items:           make(map[string]*list.Element),
		evictList:       list.New(),
		cleanupInterval: 5 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}

	go cache.cleanupLoop()

	return cache
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, exists := c.items[key]
	if !exists {
		c.stats.Misses++
		return nil, ErrCacheMiss
	}

	// Get the value as interface{} first
	val := interface{}(element.Value)
	entry := val.(*Entry)

	if time.Now().After(entry.ExpireAt) {
		c.removeElement(element)
		c.stats.Misses++
		return nil, ErrExpired
	}

	entry.AccessedAt = time.Now()
	c.evictList.MoveToFront(element)
	c.stats.Hits++

	return entry.Value, nil
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	size := c.estimateSize(value)

	if size > c.maxSize {
		return fmt.Errorf("value too large for cache")
	}

	// Check if key already exists
	if element, exists := c.items[key]; exists {
		val := interface{}(element.Value)
		entry := val.(*Entry)
		c.size -= entry.Size
		entry.Value = value
		entry.Size = size
		entry.ExpireAt = time.Now().Add(ttl)
		entry.AccessedAt = time.Now()
		c.size += size
		c.evictList.MoveToFront(element)
		return nil
	}

	// Evict items if necessary
	for c.size+size > c.maxSize || len(c.items) >= c.maxItems {
		c.evictOldest()
	}

	entry := &Entry{
		Key:        key,
		Value:      value,
		ExpireAt:   time.Now().Add(ttl),
		AccessedAt: time.Now(),
		Size:       size,
	}

	elem := c.evictList.PushFront(entry)
	c.items[key] = elem
	c.size += size

	return nil
}

func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.items[key]
	if !exists {
		return ErrCacheMiss
	}

	c.removeElement(elem)
	return nil
}

func (c *LRUCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.evictList.Init()
	c.size = 0

	return nil
}

func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.stats.Size = c.size
	c.stats.Items = len(c.items)

	return c.stats
}

func (c *LRUCache) Stop() {
	close(c.stopCleanup)
}

func (c *LRUCache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *LRUCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for elem := c.evictList.Back(); elem != nil; {
		entry := interface{}(elem.Value).(*Entry)
		if now.After(entry.ExpireAt) {
			prev := elem.Prev()
			c.removeElement(elem)
			elem = prev
		} else {
			break
		}
	}

	c.stats.LastClean = now
}

func (c *LRUCache) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
		c.stats.Evictions++
	}
}

func (c *LRUCache) removeElement(elem *list.Element) {
	entry := interface{}(elem.Value).(*Entry)
	delete(c.items, entry.Key)
	c.evictList.Remove(elem)
	c.size -= entry.Size
}

func (c *LRUCache) estimateSize(value interface{}) int64 {
	// Simple size estimation based on JSON encoding
	data, err := json.Marshal(value)
	if err != nil {
		return 1024 // Default size if marshaling fails
	}
	return int64(len(data))
}

// FileCache implements a file-based cache for persistent storage
type FileCache struct {
	cacheDir string
	ttl      time.Duration
	mu       sync.RWMutex
	stats    CacheStats
}

func NewFileCache(cacheDir string, ttl time.Duration) (*FileCache, error) {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cache := &FileCache{
		cacheDir: cacheDir,
		ttl:      ttl,
	}

	go cache.cleanupLoop()

	return cache, nil
}

func (f *FileCache) Get(key string) (interface{}, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	filePath := f.getFilePath(key)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			f.stats.Misses++
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	if time.Since(info.ModTime()) > f.ttl {
		f.stats.Misses++
		os.Remove(filePath)
		return nil, ErrExpired
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	f.stats.Hits++
	return value, nil
}

func (f *FileCache) Set(key string, value interface{}, ttl time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	filePath := f.getFilePath(key)
	return os.WriteFile(filePath, data, 0644)
}

func (f *FileCache) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := f.getFilePath(key)

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return ErrCacheMiss
		}
		return err
	}

	return nil
}

func (f *FileCache) Clear() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	entries, err := os.ReadDir(f.cacheDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".cache" {
			os.Remove(filepath.Join(f.cacheDir, entry.Name()))
		}
	}

	return nil
}

func (f *FileCache) Stats() CacheStats {
	f.mu.RLock()
	defer f.mu.RUnlock()

	entries, _ := os.ReadDir(f.cacheDir)
	f.stats.Items = 0
	f.stats.Size = 0

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".cache" {
			f.stats.Items++
			if info, err := entry.Info(); err == nil {
				f.stats.Size += info.Size()
			}
		}
	}

	return f.stats
}

func (f *FileCache) getFilePath(key string) string {
	// Use a simple hash for the filename
	return filepath.Join(f.cacheDir, fmt.Sprintf("%x.cache", key))
}

func (f *FileCache) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		f.cleanup()
	}
}

func (f *FileCache) cleanup() {
	f.mu.Lock()
	defer f.mu.Unlock()

	entries, err := os.ReadDir(f.cacheDir)
	if err != nil {
		return
	}

	now := time.Now()

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".cache" {
			filePath := filepath.Join(f.cacheDir, entry.Name())
			if info, err := entry.Info(); err == nil {
				if now.Sub(info.ModTime()) > f.ttl {
					os.Remove(filePath)
					f.stats.Evictions++
				}
			}
		}
	}

	f.stats.LastClean = now
}

// MultiLevelCache combines memory and file caches for optimal performance
type MultiLevelCache struct {
	memory *LRUCache
	file   *FileCache
	mu     sync.RWMutex
}

func NewMultiLevelCache(memSize int64, memItems int, cacheDir string, ttl time.Duration) (*MultiLevelCache, error) {
	fileCache, err := NewFileCache(cacheDir, ttl)
	if err != nil {
		return nil, err
	}

	return &MultiLevelCache{
		memory: NewLRUCache(memSize, memItems),
		file:   fileCache,
	}, nil
}

func (m *MultiLevelCache) Get(key string) (interface{}, error) {
	// Try memory cache first
	if value, err := m.memory.Get(key); err == nil {
		return value, nil
	}

	// Fall back to file cache
	value, err := m.file.Get(key)
	if err != nil {
		return nil, err
	}

	// Promote to memory cache
	m.memory.Set(key, value, 5*time.Minute)

	return value, nil
}

func (m *MultiLevelCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Write to both caches
	if err := m.memory.Set(key, value, ttl); err != nil {
		// Memory cache failure is not critical
		fmt.Printf("Warning: failed to set memory cache: %v\n", err)
	}

	return m.file.Set(key, value, ttl)
}

func (m *MultiLevelCache) Delete(key string) error {
	m.memory.Delete(key)
	return m.file.Delete(key)
}

func (m *MultiLevelCache) Clear() error {
	m.memory.Clear()
	return m.file.Clear()
}

func (m *MultiLevelCache) Stats() CacheStats {
	memStats := m.memory.Stats()
	fileStats := m.file.Stats()

	return CacheStats{
		Hits:      memStats.Hits + fileStats.Hits,
		Misses:    memStats.Misses + fileStats.Misses,
		Evictions: memStats.Evictions + fileStats.Evictions,
		Size:      memStats.Size + fileStats.Size,
		Items:     memStats.Items + fileStats.Items,
		LastClean: memStats.LastClean,
	}
}
