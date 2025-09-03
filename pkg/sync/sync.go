package sync

import (
	"bytes"
	"cheat-go/pkg/apps"
	"cheat-go/pkg/notes"
	"cheat-go/pkg/online"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrSyncInProgress = errors.New("sync already in progress")
	ErrSyncFailed     = errors.New("sync failed")
	ErrNoSyncService  = errors.New("no sync service configured")
)

type SyncService interface {
	Push(data SyncData) error
	Pull() (*SyncData, error)
	GetLastSync() (time.Time, error)
	ResolveConflict(item SyncItem, resolution ConflictResolution) error
}

type SyncData struct {
	Version     string              `json:"version"`
	Timestamp   time.Time           `json:"timestamp"`
	DeviceID    string              `json:"device_id"`
	Apps        []apps.App          `json:"apps,omitempty"`
	Notes       []*notes.Note       `json:"notes,omitempty"`
	CheatSheets []online.CheatSheet `json:"cheat_sheets,omitempty"`
	Checksum    string              `json:"checksum"`
}

type SyncItem struct {
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	Local     interface{} `json:"local"`
	Remote    interface{} `json:"remote"`
	Timestamp time.Time   `json:"timestamp"`
}

type ConflictResolution int

const (
	KeepLocal ConflictResolution = iota
	KeepRemote
	Merge
	Skip
)

type Manager struct {
	service      SyncService
	localDataDir string
	deviceID     string
	syncInterval time.Duration
	mu           sync.RWMutex
	isSyncing    bool
	lastSync     time.Time
	conflicts    []SyncItem
	stopChan     chan struct{}
}

func NewManager(service SyncService, localDataDir string) (*Manager, error) {
	deviceID, err := getOrCreateDeviceID(localDataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get device ID: %w", err)
	}

	return &Manager{
		service:      service,
		localDataDir: localDataDir,
		deviceID:     deviceID,
		syncInterval: 15 * time.Minute,
		stopChan:     make(chan struct{}),
	}, nil
}

func (m *Manager) StartAutoSync() error {
	go func() {
		ticker := time.NewTicker(m.syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.Sync(); err != nil {
					fmt.Printf("Auto-sync failed: %v\n", err)
				}
			case <-m.stopChan:
				return
			}
		}
	}()

	return nil
}

func (m *Manager) StopAutoSync() {
	close(m.stopChan)
}

func (m *Manager) Sync() error {
	m.mu.Lock()
	if m.isSyncing {
		m.mu.Unlock()
		return ErrSyncInProgress
	}
	m.isSyncing = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.isSyncing = false
		m.mu.Unlock()
	}()

	localData, err := m.gatherLocalData()
	if err != nil {
		return fmt.Errorf("failed to gather local data: %w", err)
	}

	remoteData, err := m.service.Pull()
	if err != nil {
		return fmt.Errorf("failed to pull remote data: %w", err)
	}

	conflicts := m.detectConflicts(localData, remoteData)
	if len(conflicts) > 0 {
		m.mu.Lock()
		m.conflicts = conflicts
		m.mu.Unlock()

		if err := m.autoResolveConflicts(conflicts); err != nil {
			return fmt.Errorf("failed to resolve conflicts: %w", err)
		}
	}

	mergedData := m.mergeData(localData, remoteData)

	if err := m.service.Push(*mergedData); err != nil {
		return fmt.Errorf("failed to push data: %w", err)
	}

	if err := m.saveLocalData(mergedData); err != nil {
		return fmt.Errorf("failed to save local data: %w", err)
	}

	m.mu.Lock()
	m.lastSync = time.Now()
	m.conflicts = nil
	m.mu.Unlock()

	return nil
}

func (m *Manager) GetSyncStatus() SyncStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return SyncStatus{
		LastSync:     m.lastSync,
		IsSyncing:    m.isSyncing,
		HasConflicts: len(m.conflicts) > 0,
		Conflicts:    m.conflicts,
		DeviceID:     m.deviceID,
	}
}

func (m *Manager) ResolveConflict(itemID string, resolution ConflictResolution) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, conflict := range m.conflicts {
		if conflict.ID == itemID {
			if err := m.service.ResolveConflict(conflict, resolution); err != nil {
				return err
			}

			m.conflicts = append(m.conflicts[:i], m.conflicts[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("conflict not found: %s", itemID)
}

func (m *Manager) gatherLocalData() (*SyncData, error) {
	data := &SyncData{
		Version:   "1.0",
		Timestamp: time.Now(),
		DeviceID:  m.deviceID,
	}

	appsFile := filepath.Join(m.localDataDir, "apps.json")
	if _, err := os.Stat(appsFile); err == nil {
		appsData, err := os.ReadFile(appsFile)
		if err == nil {
			json.Unmarshal(appsData, &data.Apps)
		}
	}

	notesFile := filepath.Join(m.localDataDir, "notes.json")
	if _, err := os.Stat(notesFile); err == nil {
		notesData, err := os.ReadFile(notesFile)
		if err == nil {
			json.Unmarshal(notesData, &data.Notes)
		}
	}

	data.Checksum = m.calculateChecksum(data)

	return data, nil
}

func (m *Manager) saveLocalData(data *SyncData) error {
	if len(data.Apps) > 0 {
		appsFile := filepath.Join(m.localDataDir, "apps.json")
		appsData, _ := json.MarshalIndent(data.Apps, "", "  ")
		if err := os.WriteFile(appsFile, appsData, 0644); err != nil {
			return err
		}
	}

	if len(data.Notes) > 0 {
		notesFile := filepath.Join(m.localDataDir, "notes.json")
		notesData, _ := json.MarshalIndent(data.Notes, "", "  ")
		if err := os.WriteFile(notesFile, notesData, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) detectConflicts(local, remote *SyncData) []SyncItem {
	conflicts := []SyncItem{}

	if local == nil || remote == nil {
		return conflicts
	}

	// Check for conflicting notes
	remoteNotesMap := make(map[string]*notes.Note)
	for _, note := range remote.Notes {
		remoteNotesMap[note.ID] = note
	}

	for _, localNote := range local.Notes {
		if remoteNote, exists := remoteNotesMap[localNote.ID]; exists {
			if localNote.UpdatedAt != remoteNote.UpdatedAt {
				conflicts = append(conflicts, SyncItem{
					Type:      "note",
					ID:        localNote.ID,
					Local:     localNote,
					Remote:    remoteNote,
					Timestamp: time.Now(),
				})
			}
		}
	}

	return conflicts
}

func (m *Manager) mergeData(local, remote *SyncData) *SyncData {
	merged := &SyncData{
		Version:   "1.0",
		Timestamp: time.Now(),
		DeviceID:  m.deviceID,
	}

	if local != nil {
		merged.Apps = local.Apps
		merged.Notes = local.Notes
		merged.CheatSheets = local.CheatSheets
	}

	if remote != nil && remote.Timestamp.After(local.Timestamp) {
		merged.Apps = remote.Apps
		merged.Notes = remote.Notes
		merged.CheatSheets = remote.CheatSheets
	}

	merged.Checksum = m.calculateChecksum(merged)

	return merged
}

func (m *Manager) autoResolveConflicts(conflicts []SyncItem) error {
	for _, conflict := range conflicts {
		resolution := m.determineResolution(conflict)
		if err := m.service.ResolveConflict(conflict, resolution); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) determineResolution(conflict SyncItem) ConflictResolution {
	switch conflict.Type {
	case "note":
		localNote := conflict.Local.(*notes.Note)
		remoteNote := conflict.Remote.(*notes.Note)

		if localNote.UpdatedAt.After(remoteNote.UpdatedAt) {
			return KeepLocal
		}
		return KeepRemote

	default:
		return KeepLocal
	}
}

func (m *Manager) calculateChecksum(data *SyncData) string {
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

func getOrCreateDeviceID(dataDir string) (string, error) {
	deviceFile := filepath.Join(dataDir, ".device_id")

	if data, err := os.ReadFile(deviceFile); err == nil {
		return string(data), nil
	}

	deviceID := generateDeviceID()

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(deviceFile, []byte(deviceID), 0644); err != nil {
		return "", err
	}

	return deviceID, nil
}

func generateDeviceID() string {
	hostname, _ := os.Hostname()
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s-%d", hostname, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

type SyncStatus struct {
	LastSync     time.Time  `json:"last_sync"`
	IsSyncing    bool       `json:"is_syncing"`
	HasConflicts bool       `json:"has_conflicts"`
	Conflicts    []SyncItem `json:"conflicts,omitempty"`
	DeviceID     string     `json:"device_id"`
}

// CloudSyncService implements sync with a cloud backend
type CloudSyncService struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

func NewCloudSyncService(endpoint, apiKey string) *CloudSyncService {
	return &CloudSyncService{
		endpoint: endpoint,
		apiKey:   apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CloudSyncService) Push(data SyncData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/push", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("push failed: %s", body)
	}

	return nil
}

func (c *CloudSyncService) Pull() (*SyncData, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/pull", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &SyncData{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pull failed: %s", body)
	}

	var data SyncData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *CloudSyncService) GetLastSync() (time.Time, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/last-sync", nil)
	if err != nil {
		return time.Time{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	var result struct {
		LastSync time.Time `json:"last_sync"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return time.Time{}, err
	}

	return result.LastSync, nil
}

func (c *CloudSyncService) ResolveConflict(item SyncItem, resolution ConflictResolution) error {
	data := map[string]interface{}{
		"item":       item,
		"resolution": resolution,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/resolve", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("conflict resolution failed: %s", body)
	}

	return nil
}
