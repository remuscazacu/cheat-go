package notes

import (
	"bytes"
	"cheat-go/pkg/apps"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFileManager_CreateNote(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	note := &Note{
		Title:    "Test Note",
		Content:  "Test content",
		AppName:  "vim",
		Category: "editor",
		Tags:     []string{"test", "sample"},
	}

	err = manager.CreateNote(note)
	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}

	if note.ID == "" {
		t.Error("Note ID should be generated")
	}

	if note.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	// Test duplicate ID
	existingNote := &Note{
		ID:      note.ID,
		Title:   "Duplicate",
		Content: "Duplicate content",
	}

	err = manager.CreateNote(existingNote)
	if err != ErrNoteExists {
		t.Errorf("Expected ErrNoteExists, got %v", err)
	}
}

func TestFileManager_GetNote(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:   "Get Test",
		Content: "Get test content",
	}

	manager.CreateNote(note)

	retrieved, err := manager.GetNote(note.ID)
	if err != nil {
		t.Fatalf("GetNote() error = %v", err)
	}

	if retrieved.Title != note.Title {
		t.Errorf("GetNote() title = %v, want %v", retrieved.Title, note.Title)
	}

	// Test non-existent note
	_, err = manager.GetNote("nonexistent")
	if err != ErrNoteNotFound {
		t.Errorf("Expected ErrNoteNotFound, got %v", err)
	}
}

func TestFileManager_UpdateNote(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:   "Original",
		Content: "Original content",
	}

	manager.CreateNote(note)
	originalCreatedAt := note.CreatedAt

	time.Sleep(10 * time.Millisecond)

	updated := &Note{
		Title:   "Updated",
		Content: "Updated content",
	}

	err = manager.UpdateNote(note.ID, updated)
	if err != nil {
		t.Fatalf("UpdateNote() error = %v", err)
	}

	retrieved, _ := manager.GetNote(note.ID)

	if retrieved.Title != "Updated" {
		t.Errorf("UpdateNote() title = %v, want Updated", retrieved.Title)
	}

	if retrieved.CreatedAt != originalCreatedAt {
		t.Error("CreatedAt should not change on update")
	}

	if !retrieved.UpdatedAt.After(originalCreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}

func TestFileManager_DeleteNote(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:   "To Delete",
		Content: "Will be deleted",
	}

	manager.CreateNote(note)

	err = manager.DeleteNote(note.ID)
	if err != nil {
		t.Fatalf("DeleteNote() error = %v", err)
	}

	_, err = manager.GetNote(note.ID)
	if err != ErrNoteNotFound {
		t.Errorf("Expected ErrNoteNotFound after delete, got %v", err)
	}

	// Test deleting non-existent note
	err = manager.DeleteNote("nonexistent")
	if err != ErrNoteNotFound {
		t.Errorf("Expected ErrNoteNotFound, got %v", err)
	}
}

func TestFileManager_SearchNotes(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test notes
	notes := []*Note{
		{
			Title:      "Vim Shortcuts",
			Content:    "Vim specific content",
			AppName:    "vim",
			Category:   "editor",
			Tags:       []string{"vim", "editor"},
			IsFavorite: true,
		},
		{
			Title:    "Git Commands",
			Content:  "Git specific content",
			AppName:  "git",
			Category: "vcs",
			Tags:     []string{"git", "version-control"},
		},
		{
			Title:    "Vim Advanced",
			Content:  "Advanced vim techniques",
			AppName:  "vim",
			Category: "editor",
			Tags:     []string{"vim", "advanced"},
		},
	}

	for _, note := range notes {
		manager.CreateNote(note)
	}

	tests := []struct {
		name     string
		opts     SearchOptions
		expected int
	}{
		{
			name:     "Search by query",
			opts:     SearchOptions{Query: "vim"},
			expected: 2,
		},
		{
			name:     "Search by app name",
			opts:     SearchOptions{AppName: "git"},
			expected: 1,
		},
		{
			name:     "Search by category",
			opts:     SearchOptions{Category: "editor"},
			expected: 2,
		},
		{
			name:     "Search favorites only",
			opts:     SearchOptions{OnlyFavorites: true},
			expected: 1,
		},
		{
			name:     "Search by tags",
			opts:     SearchOptions{Tags: []string{"advanced"}},
			expected: 1,
		},
		{
			name:     "Combined search",
			opts:     SearchOptions{AppName: "vim", Tags: []string{"editor"}},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := manager.SearchNotes(tt.opts)
			if err != nil {
				t.Fatalf("SearchNotes() error = %v", err)
			}
			if len(results) != tt.expected {
				t.Errorf("SearchNotes() returned %d notes, want %d", len(results), tt.expected)
			}
		})
	}
}

func TestFileManager_Shortcuts(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:   "Shortcut Test",
		Content: "Testing shortcuts",
	}

	manager.CreateNote(note)

	shortcut := apps.Shortcut{
		Keys:        "gg",
		Description: "Go to top",
		Category:    "navigation",
	}

	// Add shortcut
	err = manager.AddShortcutToNote(note.ID, shortcut)
	if err != nil {
		t.Fatalf("AddShortcutToNote() error = %v", err)
	}

	retrieved, _ := manager.GetNote(note.ID)
	if len(retrieved.Shortcuts) != 1 {
		t.Errorf("Expected 1 shortcut, got %d", len(retrieved.Shortcuts))
	}

	// Remove shortcut
	err = manager.RemoveShortcutFromNote(note.ID, 0)
	if err != nil {
		t.Fatalf("RemoveShortcutFromNote() error = %v", err)
	}

	retrieved, _ = manager.GetNote(note.ID)
	if len(retrieved.Shortcuts) != 0 {
		t.Errorf("Expected 0 shortcuts after removal, got %d", len(retrieved.Shortcuts))
	}

	// Test invalid index
	err = manager.RemoveShortcutFromNote(note.ID, 10)
	if err == nil {
		t.Error("Expected error for invalid shortcut index")
	}
}

func TestFileManager_ToggleFavorite(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:      "Favorite Test",
		Content:    "Testing favorites",
		IsFavorite: false,
	}

	manager.CreateNote(note)

	err = manager.ToggleFavorite(note.ID)
	if err != nil {
		t.Fatalf("ToggleFavorite() error = %v", err)
	}

	retrieved, _ := manager.GetNote(note.ID)
	if !retrieved.IsFavorite {
		t.Error("Note should be favorite after toggle")
	}

	err = manager.ToggleFavorite(note.ID)
	if err != nil {
		t.Fatalf("ToggleFavorite() error = %v", err)
	}

	retrieved, _ = manager.GetNote(note.ID)
	if retrieved.IsFavorite {
		t.Error("Note should not be favorite after second toggle")
	}
}

func TestFileManager_ExportImport(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test notes
	notes := []*Note{
		{
			Title:   "Export Test 1",
			Content: "Content 1",
			Tags:    []string{"export", "test"},
			Shortcuts: []apps.Shortcut{
				{Keys: "ctrl+c", Description: "Copy"},
			},
		},
		{
			Title:      "Export Test 2",
			Content:    "Content 2",
			IsFavorite: true,
		},
	}

	for _, note := range notes {
		manager.CreateNote(note)
	}

	// Test JSON export/import
	jsonData, err := manager.ExportNotes("json")
	if err != nil {
		t.Fatalf("ExportNotes(json) error = %v", err)
	}

	var exportedNotes []*Note
	if err := json.Unmarshal(jsonData, &exportedNotes); err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	if len(exportedNotes) != 2 {
		t.Errorf("Expected 2 exported notes, got %d", len(exportedNotes))
	}

	// Test YAML export
	yamlData, err := manager.ExportNotes("yaml")
	if err != nil {
		t.Fatalf("ExportNotes(yaml) error = %v", err)
	}

	if len(yamlData) == 0 {
		t.Error("YAML export should not be empty")
	}

	// Test Markdown export
	mdData, err := manager.ExportNotes("markdown")
	if err != nil {
		t.Fatalf("ExportNotes(markdown) error = %v", err)
	}

	mdString := string(mdData)
	if !strings.Contains(mdString, "# Personal Notes") {
		t.Error("Markdown export should contain header")
	}
	if !strings.Contains(mdString, "Export Test 1") {
		t.Error("Markdown export should contain note titles")
	}
	if !strings.Contains(mdString, "⭐") {
		t.Error("Markdown export should show favorite indicator")
	}

	// Test import
	tempDir2 := t.TempDir()
	manager2, _ := NewFileManager(tempDir2)

	err = manager2.ImportNotes(jsonData, "json")
	if err != nil {
		t.Fatalf("ImportNotes() error = %v", err)
	}

	importedNotes, _ := manager2.ListNotes()
	if len(importedNotes) != 2 {
		t.Errorf("Expected 2 imported notes, got %d", len(importedNotes))
	}

	// Test invalid format
	_, err = manager.ExportNotes("invalid")
	if err != ErrInvalidFormat {
		t.Errorf("Expected ErrInvalidFormat, got %v", err)
	}
}

func TestFileManager_Persistence(t *testing.T) {
	tempDir := t.TempDir()

	// Create manager and add notes
	manager1, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	note := &Note{
		Title:   "Persistent Note",
		Content: "This should persist",
	}

	manager1.CreateNote(note)
	noteID := note.ID

	// Create new manager instance with same directory
	manager2, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := manager2.GetNote(noteID)
	if err != nil {
		t.Fatalf("Note should persist: %v", err)
	}

	if retrieved.Title != "Persistent Note" {
		t.Errorf("Persisted note title = %v, want Persistent Note", retrieved.Title)
	}
}

func TestCloudSyncer_Sync(t *testing.T) {
	tempDir := t.TempDir()
	manager, _ := NewFileManager(tempDir)

	// Create local notes
	localNote := &Note{
		ID:        "note1",
		Title:     "Local Note",
		Content:   "Local content",
		UpdatedAt: time.Now(),
	}
	manager.CreateNote(localNote)

	// Create remote data
	remoteNote := &Note{
		ID:        "note1",
		Title:     "Remote Note",
		Content:   "Remote content",
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	remoteData, _ := json.Marshal([]*Note{remoteNote})
	remote := bytes.NewBuffer(remoteData)

	syncer := NewCloudSyncer(manager, remote)

	status, err := syncer.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if !status.HasConflicts {
		t.Error("Should detect conflict between local and remote")
	}

	if len(status.Conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(status.Conflicts))
	}
}

func TestCloudSyncer_ResolveConflict(t *testing.T) {
	tempDir := t.TempDir()
	manager, _ := NewFileManager(tempDir)

	localNote := &Note{
		ID:        "note1",
		Title:     "Local",
		Content:   "Local content",
		UpdatedAt: time.Now(),
	}
	manager.CreateNote(localNote)

	remoteNote := &Note{
		ID:        "note1",
		Title:     "Remote",
		Content:   "Remote content",
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	conflict := Conflict{
		NoteID:     "note1",
		LocalNote:  localNote,
		RemoteNote: remoteNote,
		DetectedAt: time.Now(),
	}

	remote := bytes.NewBuffer([]byte{})
	syncer := NewCloudSyncer(manager, remote)

	// Test KeepLocal resolution
	err := syncer.ResolveConflict(conflict, KeepLocal)
	if err != nil {
		t.Fatalf("ResolveConflict(KeepLocal) error = %v", err)
	}

	note, _ := manager.GetNote("note1")
	if note.Title != "Local" {
		t.Error("KeepLocal should keep local version")
	}

	// Test KeepRemote resolution
	err = syncer.ResolveConflict(conflict, KeepRemote)
	if err != nil {
		t.Fatalf("ResolveConflict(KeepRemote) error = %v", err)
	}

	note, _ = manager.GetNote("note1")
	if note.Title != "Remote" {
		t.Error("KeepRemote should use remote version")
	}

	// Test Merge resolution
	err = syncer.ResolveConflict(conflict, Merge)
	if err != nil {
		t.Fatalf("ResolveConflict(Merge) error = %v", err)
	}

	note, _ = manager.GetNote("note1")
	if !strings.Contains(note.Content, "Local content") || !strings.Contains(note.Content, "Remote content") {
		t.Error("Merge should combine both contents")
	}
}

func TestSortNotes(t *testing.T) {
	notes := []*Note{
		{Title: "B", CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now().Add(-1 * time.Hour)},
		{Title: "A", CreatedAt: time.Now().Add(-3 * time.Hour), UpdatedAt: time.Now()},
		{Title: "C", CreatedAt: time.Now().Add(-1 * time.Hour), UpdatedAt: time.Now().Add(-2 * time.Hour)},
	}

	// Test sort by title
	sortNotes(notes, "title")
	if notes[0].Title != "A" || notes[1].Title != "B" || notes[2].Title != "C" {
		t.Error("Notes should be sorted by title")
	}

	// Test sort by created_at
	sortNotes(notes, "created_at")
	if notes[0].Title != "C" {
		t.Error("Most recently created should be first")
	}

	// Test sort by updated_at
	sortNotes(notes, "updated_at")
	if notes[0].Title != "A" {
		t.Error("Most recently updated should be first")
	}
}

func TestFileManager_NewFileManagerError(t *testing.T) {
	// Test with invalid directory path
	manager, err := NewFileManager("/invalid/path/that/cannot/be/created")
	if err == nil {
		t.Error("Expected error for invalid directory path")
	}
	if manager != nil {
		t.Error("Manager should be nil on error")
	}
}

func TestFileManager_ImportNotes(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test JSON import
	notes := []*Note{
		{
			ID:       "import-1",
			Title:    "Imported Note 1",
			Content:  "Imported content 1",
			Category: "imported",
		},
		{
			ID:       "import-2",
			Title:    "Imported Note 2",
			Content:  "Imported content 2",
			Category: "imported",
		},
	}

	jsonData, _ := json.Marshal(notes)
	err = manager.ImportNotes(jsonData, "json")
	if err != nil {
		t.Fatalf("ImportNotes() error = %v", err)
	}

	// Verify imported notes
	importedNote, err := manager.GetNote("import-1")
	if err != nil {
		t.Errorf("Failed to get imported note: %v", err)
	}
	if importedNote.Title != "Imported Note 1" {
		t.Error("Imported note title mismatch")
	}

	// Test YAML import
	yamlData := `
- id: yaml-1
  title: YAML Note
  content: YAML content
  category: yaml
`
	err = manager.ImportNotes([]byte(yamlData), "yaml")
	if err != nil {
		t.Fatalf("ImportNotes YAML error = %v", err)
	}

	// Test invalid format
	err = manager.ImportNotes([]byte("invalid"), "invalid")
	if err != ErrInvalidFormat {
		t.Errorf("Expected ErrInvalidFormat, got %v", err)
	}

	// Test invalid JSON
	err = manager.ImportNotes([]byte("invalid json"), "json")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test invalid YAML
	err = manager.ImportNotes([]byte("invalid: yaml: ["), "yaml")
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestFileManager_AutoResolve(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	syncer := NewCloudSyncer(manager, &bytes.Buffer{})

	// Test empty conflicts
	err = syncer.AutoResolve([]Conflict{})
	if err != nil {
		t.Errorf("AutoResolve empty conflicts error = %v", err)
	}

	// Create a test note for conflicts
	localNote := &Note{
		ID:        "conflict-note",
		Title:     "Local Title",
		Content:   "Local content",
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	remoteNote := &Note{
		ID:        "conflict-note",
		Title:     "Remote Title",
		Content:   "Remote content",
		UpdatedAt: time.Now(),
	}

	manager.CreateNote(localNote)

	conflicts := []Conflict{
		{
			NoteID:     "conflict-note",
			LocalNote:  localNote,
			RemoteNote: remoteNote,
			DetectedAt: time.Now(),
		},
	}

	err = syncer.AutoResolve(conflicts)
	if err != nil {
		t.Errorf("AutoResolve error = %v", err)
	}
}

func TestFileManager_MinFunction(t *testing.T) {
	// Test the min helper function
	result := min(5, 3)
	if result != 3 {
		t.Errorf("min(5, 3) = %d, want 3", result)
	}

	result = min(2, 8)
	if result != 2 {
		t.Errorf("min(2, 8) = %d, want 2", result)
	}

	result = min(4, 4)
	if result != 4 {
		t.Errorf("min(4, 4) = %d, want 4", result)
	}
}

func TestFileManager_RemoveShortcutEdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewFileManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	note := &Note{
		ID:      "shortcut-test",
		Title:   "Shortcut Test",
		Content: "Test content",
		Shortcuts: []apps.Shortcut{
			{Keys: "ctrl+a", Description: "Select all"},
			{Keys: "ctrl+c", Description: "Copy"},
		},
	}

	err = manager.CreateNote(note)
	if err != nil {
		t.Fatalf("CreateNote error = %v", err)
	}

	// Test removing shortcut with invalid index (negative)
	err = manager.RemoveShortcutFromNote(note.ID, -1)
	if err == nil {
		t.Error("Expected error for negative shortcut index")
	}

	// Test removing shortcut with invalid index (too large)
	err = manager.RemoveShortcutFromNote(note.ID, 10)
	if err == nil {
		t.Error("Expected error for out-of-bounds shortcut index")
	}

	// Test removing shortcut from non-existent note
	err = manager.RemoveShortcutFromNote("non-existent", 0)
	if err != ErrNoteNotFound {
		t.Errorf("Expected ErrNoteNotFound, got %v", err)
	}
}

func TestFileManager_LoadNotesErrors(t *testing.T) {
	tempDir := t.TempDir()

	// Create invalid JSON file
	invalidFile := tempDir + "/notes.json"
	err := os.WriteFile(invalidFile, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	// This should handle the error gracefully
	manager, err := NewFileManager(tempDir)
	if err == nil {
		t.Error("Expected error for invalid JSON file")
	}
	if manager != nil {
		t.Error("Manager should be nil on error")
	}
}
