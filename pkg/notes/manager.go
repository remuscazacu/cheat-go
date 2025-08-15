package notes

import (
	"cheat-go/pkg/apps"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoteNotFound  = errors.New("note not found")
	ErrInvalidFormat = errors.New("invalid format")
	ErrNoteExists    = errors.New("note already exists")
)

type FileManager struct {
	dataDir string
	mu      sync.RWMutex
	notes   map[string]*Note
}

func NewFileManager(dataDir string) (*FileManager, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create notes directory: %w", err)
	}

	fm := &FileManager{
		dataDir: dataDir,
		notes:   make(map[string]*Note),
	}

	if err := fm.loadNotes(); err != nil {
		return nil, fmt.Errorf("failed to load notes: %w", err)
	}

	return fm, nil
}

func (fm *FileManager) loadNotes() error {
	notesFile := filepath.Join(fm.dataDir, "notes.json")

	if _, err := os.Stat(notesFile); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(notesFile)
	if err != nil {
		return fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []*Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return fmt.Errorf("failed to unmarshal notes: %w", err)
	}

	for _, note := range notes {
		fm.notes[note.ID] = note
	}

	return nil
}

func (fm *FileManager) saveNotes() error {
	notesFile := filepath.Join(fm.dataDir, "notes.json")

	notes := make([]*Note, 0, len(fm.notes))
	for _, note := range fm.notes {
		notes = append(notes, note)
	}

	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %w", err)
	}

	if err := os.WriteFile(notesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %w", err)
	}

	return nil
}

func (fm *FileManager) CreateNote(note *Note) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if note.ID == "" {
		note.ID = generateID()
	}

	if _, exists := fm.notes[note.ID]; exists {
		return ErrNoteExists
	}

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	fm.notes[note.ID] = note
	return fm.saveNotes()
}

func (fm *FileManager) GetNote(id string) (*Note, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	note, exists := fm.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}

	return note, nil
}

func (fm *FileManager) UpdateNote(id string, updatedNote *Note) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	note, exists := fm.notes[id]
	if !exists {
		return ErrNoteNotFound
	}

	updatedNote.ID = id
	updatedNote.CreatedAt = note.CreatedAt
	updatedNote.UpdatedAt = time.Now()

	fm.notes[id] = updatedNote
	return fm.saveNotes()
}

func (fm *FileManager) DeleteNote(id string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, exists := fm.notes[id]; !exists {
		return ErrNoteNotFound
	}

	delete(fm.notes, id)
	return fm.saveNotes()
}

func (fm *FileManager) SearchNotes(opts SearchOptions) ([]*Note, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	results := []*Note{}

	for _, note := range fm.notes {
		if !matchesSearchOptions(note, opts) {
			continue
		}
		results = append(results, note)
	}

	sortNotes(results, opts.SortBy)

	if opts.Limit > 0 && len(results) > opts.Limit {
		results = results[opts.Offset:min(opts.Offset+opts.Limit, len(results))]
	}

	return results, nil
}

func (fm *FileManager) ListNotes() ([]*Note, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	notes := make([]*Note, 0, len(fm.notes))
	for _, note := range fm.notes {
		notes = append(notes, note)
	}

	sortNotes(notes, "updated_at")
	return notes, nil
}

func (fm *FileManager) AddShortcutToNote(noteID string, shortcut apps.Shortcut) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	note, exists := fm.notes[noteID]
	if !exists {
		return ErrNoteNotFound
	}

	note.Shortcuts = append(note.Shortcuts, shortcut)
	note.UpdatedAt = time.Now()

	return fm.saveNotes()
}

func (fm *FileManager) RemoveShortcutFromNote(noteID string, shortcutIndex int) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	note, exists := fm.notes[noteID]
	if !exists {
		return ErrNoteNotFound
	}

	if shortcutIndex < 0 || shortcutIndex >= len(note.Shortcuts) {
		return fmt.Errorf("invalid shortcut index")
	}

	note.Shortcuts = append(note.Shortcuts[:shortcutIndex], note.Shortcuts[shortcutIndex+1:]...)
	note.UpdatedAt = time.Now()

	return fm.saveNotes()
}

func (fm *FileManager) ToggleFavorite(id string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	note, exists := fm.notes[id]
	if !exists {
		return ErrNoteNotFound
	}

	note.IsFavorite = !note.IsFavorite
	note.UpdatedAt = time.Now()

	return fm.saveNotes()
}

func (fm *FileManager) ExportNotes(format string) ([]byte, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	notes := make([]*Note, 0, len(fm.notes))
	for _, note := range fm.notes {
		notes = append(notes, note)
	}

	switch format {
	case "json":
		return json.MarshalIndent(notes, "", "  ")
	case "yaml", "yml":
		return yaml.Marshal(notes)
	case "markdown", "md":
		return exportToMarkdown(notes), nil
	default:
		return nil, ErrInvalidFormat
	}
}

func (fm *FileManager) ImportNotes(data []byte, format string) error {
	var notes []*Note

	switch format {
	case "json":
		if err := json.Unmarshal(data, &notes); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &notes); err != nil {
			return fmt.Errorf("failed to unmarshal YAML: %w", err)
		}
	default:
		return ErrInvalidFormat
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, note := range notes {
		if note.ID == "" {
			note.ID = generateID()
		}
		if _, exists := fm.notes[note.ID]; !exists {
			fm.notes[note.ID] = note
		}
	}

	return fm.saveNotes()
}

func matchesSearchOptions(note *Note, opts SearchOptions) bool {
	if opts.Query != "" {
		query := strings.ToLower(opts.Query)
		if !strings.Contains(strings.ToLower(note.Title), query) &&
			!strings.Contains(strings.ToLower(note.Content), query) {
			return false
		}
	}

	if opts.AppName != "" && note.AppName != opts.AppName {
		return false
	}

	if opts.Category != "" && note.Category != opts.Category {
		return false
	}

	if opts.OnlyFavorites && !note.IsFavorite {
		return false
	}

	if len(opts.Tags) > 0 {
		hasTag := false
		for _, searchTag := range opts.Tags {
			for _, noteTag := range note.Tags {
				if searchTag == noteTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	return true
}

func sortNotes(notes []*Note, sortBy string) {
	switch sortBy {
	case "title":
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].Title < notes[j].Title
		})
	case "created_at":
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].CreatedAt.After(notes[j].CreatedAt)
		})
	case "updated_at", "":
		sort.Slice(notes, func(i, j int) bool {
			return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
		})
	}
}

func exportToMarkdown(notes []*Note) []byte {
	var sb strings.Builder

	sb.WriteString("# Personal Notes\n\n")

	for _, note := range notes {
		if note.IsFavorite {
			sb.WriteString("⭐ ")
		}
		sb.WriteString(fmt.Sprintf("## %s\n\n", note.Title))

		if note.AppName != "" {
			sb.WriteString(fmt.Sprintf("**App:** %s\n", note.AppName))
		}
		if note.Category != "" {
			sb.WriteString(fmt.Sprintf("**Category:** %s\n", note.Category))
		}
		if len(note.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("**Tags:** %s\n", strings.Join(note.Tags, ", ")))
		}

		sb.WriteString(fmt.Sprintf("\n%s\n", note.Content))

		if len(note.Shortcuts) > 0 {
			sb.WriteString("\n### Shortcuts\n\n")
			for _, shortcut := range note.Shortcuts {
				sb.WriteString(fmt.Sprintf("- `%s`: %s\n", shortcut.Keys, shortcut.Description))
			}
		}

		sb.WriteString(fmt.Sprintf("\n*Created: %s | Updated: %s*\n\n---\n\n",
			note.CreatedAt.Format("2006-01-02"),
			note.UpdatedAt.Format("2006-01-02")))
	}

	return []byte(sb.String())
}

func generateID() string {
	return fmt.Sprintf("note-%d", time.Now().UnixNano())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type CloudSyncer struct {
	local     Manager
	remote    io.ReadWriter
	conflicts []Conflict
	mu        sync.RWMutex
}

func NewCloudSyncer(local Manager, remote io.ReadWriter) *CloudSyncer {
	return &CloudSyncer{
		local:     local,
		remote:    remote,
		conflicts: []Conflict{},
	}
}

func (cs *CloudSyncer) Sync() (*SyncStatus, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	localNotes, err := cs.local.ListNotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get local notes: %w", err)
	}

	remoteData, err := io.ReadAll(cs.remote)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote notes: %w", err)
	}

	var remoteNotes []*Note
	if len(remoteData) > 0 {
		if err := json.Unmarshal(remoteData, &remoteNotes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal remote notes: %w", err)
		}
	}

	cs.conflicts = detectConflicts(localNotes, remoteNotes)

	status := &SyncStatus{
		LastSync:     time.Now(),
		TotalNotes:   len(localNotes),
		SyncedNotes:  len(localNotes) - len(cs.conflicts),
		HasConflicts: len(cs.conflicts) > 0,
		Conflicts:    cs.conflicts,
	}

	return status, nil
}

func (cs *CloudSyncer) ResolveConflict(conflict Conflict, resolution Resolution) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	switch resolution {
	case KeepLocal:
		// Keep local version, no action needed
		return nil
	case KeepRemote:
		// Update local with remote version
		return cs.local.UpdateNote(conflict.NoteID, conflict.RemoteNote)
	case Merge:
		// Merge notes (simple strategy: combine content)
		merged := mergeNotes(conflict.LocalNote, conflict.RemoteNote)
		return cs.local.UpdateNote(conflict.NoteID, merged)
	default:
		return fmt.Errorf("invalid resolution type")
	}
}

func (cs *CloudSyncer) AutoResolve(conflicts []Conflict) error {
	for _, conflict := range conflicts {
		resolution := determineAutoResolution(conflict)
		if err := cs.ResolveConflict(conflict, resolution); err != nil {
			return fmt.Errorf("failed to resolve conflict for note %s: %w", conflict.NoteID, err)
		}
	}
	return nil
}

func detectConflicts(local, remote []*Note) []Conflict {
	conflicts := []Conflict{}
	remoteMap := make(map[string]*Note)

	for _, note := range remote {
		remoteMap[note.ID] = note
	}

	for _, localNote := range local {
		if remoteNote, exists := remoteMap[localNote.ID]; exists {
			if localNote.UpdatedAt != remoteNote.UpdatedAt {
				conflicts = append(conflicts, Conflict{
					NoteID:     localNote.ID,
					LocalNote:  localNote,
					RemoteNote: remoteNote,
					DetectedAt: time.Now(),
				})
			}
		}
	}

	return conflicts
}

func mergeNotes(local, remote *Note) *Note {
	merged := &Note{
		ID:         local.ID,
		Title:      local.Title,
		Content:    fmt.Sprintf("%s\n\n--- Remote Version ---\n\n%s", local.Content, remote.Content),
		AppName:    local.AppName,
		Category:   local.Category,
		Tags:       mergeTags(local.Tags, remote.Tags),
		CreatedAt:  local.CreatedAt,
		UpdatedAt:  time.Now(),
		IsFavorite: local.IsFavorite || remote.IsFavorite,
		Shortcuts:  mergeShortcuts(local.Shortcuts, remote.Shortcuts),
	}

	if remote.UpdatedAt.After(local.UpdatedAt) {
		merged.Title = remote.Title
		merged.AppName = remote.AppName
		merged.Category = remote.Category
	}

	return merged
}

func mergeTags(local, remote []string) []string {
	tagMap := make(map[string]bool)
	for _, tag := range local {
		tagMap[tag] = true
	}
	for _, tag := range remote {
		tagMap[tag] = true
	}

	merged := []string{}
	for tag := range tagMap {
		merged = append(merged, tag)
	}
	sort.Strings(merged)
	return merged
}

func mergeShortcuts(local, remote []apps.Shortcut) []apps.Shortcut {
	shortcutMap := make(map[string]apps.Shortcut)

	for _, s := range local {
		shortcutMap[s.Keys] = s
	}
	for _, s := range remote {
		if _, exists := shortcutMap[s.Keys]; !exists {
			shortcutMap[s.Keys] = s
		}
	}

	merged := []apps.Shortcut{}
	for _, s := range shortcutMap {
		merged = append(merged, s)
	}

	return merged
}

func determineAutoResolution(conflict Conflict) Resolution {
	if conflict.LocalNote.UpdatedAt.After(conflict.RemoteNote.UpdatedAt) {
		return KeepLocal
	}
	return KeepRemote
}
