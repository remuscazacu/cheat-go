package sync

import (
	"testing"
	"time"
)

// Mock sync service for testing
type mockSyncService struct {
	pushCalled   bool
	pullCalled   bool
	lastSyncTime time.Time
	returnError  bool
	returnData   *SyncData
}

func (m *mockSyncService) Push(data SyncData) error {
	m.pushCalled = true
	if m.returnError {
		return ErrSyncFailed
	}
	return nil
}

func (m *mockSyncService) Pull() (*SyncData, error) {
	m.pullCalled = true
	if m.returnError {
		return nil, ErrSyncFailed
	}
	if m.returnData != nil {
		return m.returnData, nil
	}
	return &SyncData{
		Version:   "1.0",
		Timestamp: time.Now(),
		DeviceID:  "remote-device",
	}, nil
}

func (m *mockSyncService) GetLastSync() (time.Time, error) {
	if m.returnError {
		return time.Time{}, ErrSyncFailed
	}
	return m.lastSyncTime, nil
}

func (m *mockSyncService) ResolveConflict(item SyncItem, resolution ConflictResolution) error {
	if m.returnError {
		return ErrSyncFailed
	}
	return nil
}

func TestNewManager(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if manager == nil {
		t.Fatal("Manager should not be nil")
	}

	if manager.localDataDir != tmpDir {
		t.Error("Local data dir not set correctly")
	}

	if manager.deviceID == "" {
		t.Error("Device ID should be generated")
	}
}

func TestManager_Sync(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Perform sync
	err = manager.Sync()
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	// Check that service methods were called
	if !service.pullCalled {
		t.Error("Pull should be called during sync")
	}

	if !service.pushCalled {
		t.Error("Push should be called during sync")
	}
}

func TestManager_SyncWithError(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{returnError: true}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Perform sync that should fail
	err = manager.Sync()
	if err == nil {
		t.Error("Sync should fail when service returns error")
	}
}

func TestManager_ConcurrentSync(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Start first sync
	done := make(chan bool)
	go func() {
		manager.Sync()
		done <- true
	}()

	// Small delay to ensure first sync starts
	time.Sleep(10 * time.Millisecond)

	// Try second sync - should return ErrSyncInProgress
	err = manager.Sync()
	// The error might be ErrSyncInProgress or nil if first sync finished quickly
	if err != nil && err != ErrSyncInProgress {
		t.Errorf("Unexpected error: %v", err)
	}

	// Wait for first sync to complete
	<-done
}

func TestManager_GetSyncStatus(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	status := manager.GetSyncStatus()

	if status.DeviceID != manager.deviceID {
		t.Error("Status should include correct device ID")
	}

	if status.IsSyncing {
		t.Error("Should not be syncing initially")
	}
}

func TestManager_StartAutoSync(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Set short interval for testing
	manager.syncInterval = 50 * time.Millisecond

	// Start auto-sync
	err = manager.StartAutoSync()
	if err != nil {
		t.Fatalf("StartAutoSync failed: %v", err)
	}

	// Wait for at least one sync
	time.Sleep(100 * time.Millisecond)

	// Check that sync was called
	if !service.pullCalled || !service.pushCalled {
		t.Error("Auto-sync should have triggered sync")
	}

	// Stop
	manager.StopAutoSync()
}

func TestManager_StopAutoSync(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Start auto-sync
	err = manager.StartAutoSync()
	if err != nil {
		t.Fatalf("StartAutoSync failed: %v", err)
	}

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	// Stop it
	manager.StopAutoSync()

	// Should still be able to use manager after stop
	status := manager.GetSyncStatus()
	if status.DeviceID == "" {
		t.Error("Should still be able to get status after stop")
	}
}

func TestManager_ResolveConflict(t *testing.T) {
	tmpDir := t.TempDir()
	service := &mockSyncService{}

	manager, err := NewManager(service, tmpDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	// Add a conflict
	manager.conflicts = []SyncItem{
		{
			Type:      "app",
			ID:        "vim",
			Local:     "local-data",
			Remote:    "remote-data",
			Timestamp: time.Now(),
		},
	}

	// Try to resolve with invalid ID
	err = manager.ResolveConflict("invalid-id", KeepLocal)
	if err == nil {
		t.Error("Should error with invalid conflict ID")
	}
}

func TestGetOrCreateDeviceID(t *testing.T) {
	tmpDir := t.TempDir()

	// First call should create ID
	id1, err := getOrCreateDeviceID(tmpDir)
	if err != nil {
		t.Fatalf("getOrCreateDeviceID failed: %v", err)
	}

	if id1 == "" {
		t.Error("Device ID should not be empty")
	}

	// Second call should return same ID
	id2, err := getOrCreateDeviceID(tmpDir)
	if err != nil {
		t.Fatalf("getOrCreateDeviceID failed: %v", err)
	}

	if id1 != id2 {
		t.Error("Device ID should be persistent")
	}
}

func TestConflictResolution(t *testing.T) {
	tests := []struct {
		resolution ConflictResolution
		value      int
	}{
		{KeepLocal, 0},
		{KeepRemote, 1},
		{Merge, 2},
		{Skip, 3},
	}

	for _, test := range tests {
		if int(test.resolution) != test.value {
			t.Errorf("Resolution %v should have value %d", test.resolution, test.value)
		}
	}
}

func TestSyncStatus(t *testing.T) {
	status := SyncStatus{
		LastSync:     time.Now(),
		IsSyncing:    true,
		DeviceID:     "test-device",
		HasConflicts: true,
		Conflicts: []SyncItem{
			{Type: "test", ID: "1"},
		},
	}

	if status.DeviceID != "test-device" {
		t.Error("DeviceID not set")
	}

	if !status.IsSyncing {
		t.Error("IsSyncing should be true")
	}

	if !status.HasConflicts {
		t.Error("HasConflicts should be true")
	}

	if len(status.Conflicts) != 1 {
		t.Error("Should have one conflict")
	}
}

func TestNewCloudSyncService(t *testing.T) {
	service := NewCloudSyncService("https://example.com", "test-key")

	if service == nil {
		t.Fatal("CloudSyncService should not be nil")
	}

	if service.endpoint != "https://example.com" {
		t.Error("Endpoint not set correctly")
	}

	if service.apiKey != "test-key" {
		t.Error("API key not set correctly")
	}
}
