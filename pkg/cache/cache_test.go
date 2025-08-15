package cache

import (
	"testing"
	"time"
)

func TestNewLRUCache(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100) // 1MB, 100 items
	if cache == nil {
		t.Fatal("NewLRUCache should return non-nil cache")
	}
	if cache.maxSize != 1024*1024 {
		t.Error("cache maxSize not set correctly")
	}
	if cache.maxItems != 100 {
		t.Error("cache maxItems not set correctly")
	}
}

func TestLRUCache_SetAndGet(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Test Set and Get
	err := cache.Set("key1", "value1", 1*time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	// Test cache miss
	_, err = cache.Get("nonexistent")
	if err != ErrCacheMiss {
		t.Error("Expected ErrCacheMiss for nonexistent key")
	}
}

func TestLRUCache_Expiration(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Set with very short TTL
	err := cache.Set("key1", "value1", 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(2 * time.Millisecond)

	// Should return expired error
	_, err = cache.Get("key1")
	if err != ErrExpired {
		t.Error("Expected ErrExpired for expired key")
	}
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Set a value
	cache.Set("key1", "value1", 1*time.Hour)

	// Delete it
	err := cache.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Should not exist anymore
	_, err = cache.Get("key1")
	if err != ErrCacheMiss {
		t.Error("Expected ErrCacheMiss after delete")
	}

	// Delete non-existent key should return ErrCacheMiss
	err = cache.Delete("nonexistent")
	if err != ErrCacheMiss {
		t.Error("Delete of non-existent key should return ErrCacheMiss")
	}
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Add multiple items
	cache.Set("key1", "value1", 1*time.Hour)
	cache.Set("key2", "value2", 1*time.Hour)
	cache.Set("key3", "value3", 1*time.Hour)

	// Clear cache
	err := cache.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// All items should be gone
	_, err = cache.Get("key1")
	if err != ErrCacheMiss {
		t.Error("Cache should be empty after Clear")
	}

	stats := cache.Stats()
	if stats.Items != 0 {
		t.Error("Stats should show 0 items after Clear")
	}
}

func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Initial stats
	stats := cache.Stats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Error("Initial stats should be zero")
	}

	// Add item and get it
	cache.Set("key1", "value1", 1*time.Hour)
	cache.Get("key1") // Hit
	cache.Get("key2") // Miss

	stats = cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.Items != 1 {
		t.Errorf("Expected 1 item, got %d", stats.Items)
	}
}

func TestLRUCache_MaxSize(t *testing.T) {
	// Small cache that can only hold a few items
	cache := NewLRUCache(100, 10)

	// Try to add a value that's too large
	largeValue := make([]byte, 200)
	err := cache.Set("large", largeValue, 1*time.Hour)
	if err == nil {
		t.Error("Should error when value exceeds max size")
	}

	// Add values that will trigger eviction
	for i := 0; i < 5; i++ {
		key := string(rune('a' + i))
		cache.Set(key, "value", 1*time.Hour)
	}

	stats := cache.Stats()
	if stats.Items > 10 {
		t.Error("Cache should not exceed max items")
	}
}

func TestLRUCache_MaxItems(t *testing.T) {
	cache := NewLRUCache(1024*1024, 3) // Only 3 items max

	// Add 4 items, oldest should be evicted
	cache.Set("key1", "value1", 1*time.Hour)
	cache.Set("key2", "value2", 1*time.Hour)
	cache.Set("key3", "value3", 1*time.Hour)
	cache.Set("key4", "value4", 1*time.Hour)

	// key1 should be evicted
	_, err := cache.Get("key1")
	if err != ErrCacheMiss {
		t.Error("Oldest item should be evicted when max items exceeded")
	}

	// Others should still exist
	_, err = cache.Get("key4")
	if err != nil {
		t.Error("Newest item should still exist")
	}
}

func TestLRUCache_UpdateExisting(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Set initial value
	cache.Set("key1", "value1", 1*time.Hour)

	// Update with new value
	cache.Set("key1", "updated", 1*time.Hour)

	val, err := cache.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "updated" {
		t.Errorf("Expected updated value, got %v", val)
	}
}

func TestLRUCache_LRUOrder(t *testing.T) {
	cache := NewLRUCache(1024*1024, 3)

	// Add 3 items
	cache.Set("key1", "value1", 1*time.Hour)
	cache.Set("key2", "value2", 1*time.Hour)
	cache.Set("key3", "value3", 1*time.Hour)

	// Access key1 to make it recently used
	cache.Get("key1")

	// Add key4, should evict key2 (key1 was accessed, key3 is newer)
	cache.Set("key4", "value4", 1*time.Hour)

	// key2 should be evicted
	_, err := cache.Get("key2")
	if err != ErrCacheMiss {
		t.Error("key2 should be evicted as least recently used")
	}

	// key1 should still exist (was accessed)
	_, err = cache.Get("key1")
	if err != nil {
		t.Error("key1 should still exist after access")
	}
}

func TestLRUCache_Cleanup(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Set with short TTL
	cache.Set("expire1", "value1", 1*time.Millisecond)
	cache.Set("expire2", "value2", 1*time.Millisecond)
	cache.Set("keep", "value3", 1*time.Hour)

	// Wait for expiration
	time.Sleep(2 * time.Millisecond)

	// Manually trigger cleanup
	cache.cleanup()

	// Expired items should be gone
	_, err := cache.Get("expire1")
	if err != ErrCacheMiss && err != ErrExpired {
		t.Error("Expired item should be cleaned up")
	}

	// Non-expired should remain
	_, err = cache.Get("keep")
	if err != nil {
		t.Error("Non-expired item should remain after cleanup")
	}
}

func TestLRUCache_ConcurrentAccess(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)
	done := make(chan bool)

	// Concurrent writes
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(string(rune(i)), i, 1*time.Hour)
		}
		done <- true
	}()

	// Concurrent reads
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get(string(rune(i)))
		}
		done <- true
	}()

	// Concurrent deletes
	go func() {
		for i := 0; i < 50; i++ {
			cache.Delete(string(rune(i)))
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}

	// Should not panic or deadlock
}

func TestLRUCache_EstimateSize(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Test different types
	size1 := cache.estimateSize("string")
	if size1 == 0 {
		t.Error("String should have non-zero size")
	}

	size2 := cache.estimateSize(12345)
	if size2 == 0 {
		t.Error("Int should have non-zero size")
	}

	size3 := cache.estimateSize(map[string]string{"key": "value"})
	if size3 == 0 {
		t.Error("Map should have non-zero size")
	}

	size4 := cache.estimateSize([]int{1, 2, 3, 4, 5})
	if size4 == 0 {
		t.Error("Slice should have non-zero size")
	}

	// Test struct
	type TestStruct struct {
		Field1 string
		Field2 int
	}
	size5 := cache.estimateSize(TestStruct{Field1: "test", Field2: 42})
	if size5 == 0 {
		t.Error("Struct should have non-zero size")
	}

	// Test nil
	size6 := cache.estimateSize(nil)
	if size6 == 0 {
		t.Error("Nil should have some size")
	}
}

func TestLRUCache_RemoveElement(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Add and then cause removal through expiration
	cache.Set("key1", "value1", 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	// Access to trigger removal
	_, err := cache.Get("key1")
	if err != ErrExpired {
		t.Error("Should return ErrExpired")
	}

	// Stats should show removal
	stats := cache.Stats()
	if stats.Items != 0 {
		t.Error("Expired item should be removed")
	}
}

func TestLRUCache_EvictOldest(t *testing.T) {
	// Very small cache
	cache := NewLRUCache(50, 2)

	// Fill cache
	cache.Set("a", "val1", 1*time.Hour)
	cache.Set("b", "val2", 1*time.Hour)

	// This should evict 'a'
	cache.Set("c", "val3", 1*time.Hour)

	_, err := cache.Get("a")
	if err != ErrCacheMiss {
		t.Error("Oldest item should be evicted")
	}

	// b and c should exist
	_, err = cache.Get("b")
	if err != nil {
		t.Error("Item b should still exist")
	}
	_, err = cache.Get("c")
	if err != nil {
		t.Error("Item c should exist")
	}
}

func TestLRUCache_Stop(t *testing.T) {
	cache := NewLRUCache(1024*1024, 100)

	// Add some items
	cache.Set("key1", "value1", 1*time.Hour)

	// Stop the cache
	cache.Stop()

	// Cache should still be usable after Stop
	val, err := cache.Get("key1")
	if err != nil || val != "value1" {
		t.Error("Cache should still work after Stop")
	}
}

func TestLRUCache_CleanupLoop(t *testing.T) {
	// This test ensures cleanup loop doesn't panic
	cache := NewLRUCache(1024*1024, 100)

	// Add items with short TTL
	for i := 0; i < 10; i++ {
		cache.Set(string(rune('a'+i)), "value", 10*time.Millisecond)
	}

	// Let cleanup run
	time.Sleep(20 * time.Millisecond)

	// Force cleanup
	cache.cleanup()

	// All should be expired and cleaned
	stats := cache.Stats()
	if stats.Items > 0 {
		t.Error("All expired items should be cleaned")
	}
}
