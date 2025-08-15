package notes

import (
	"cheat-go/pkg/apps"
	"time"
)

type Note struct {
	ID         string          `json:"id" yaml:"id"`
	Title      string          `json:"title" yaml:"title"`
	Content    string          `json:"content" yaml:"content"`
	AppName    string          `json:"app_name" yaml:"app_name"`
	Category   string          `json:"category" yaml:"category"`
	Tags       []string        `json:"tags" yaml:"tags"`
	CreatedAt  time.Time       `json:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" yaml:"updated_at"`
	IsFavorite bool            `json:"is_favorite" yaml:"is_favorite"`
	Shortcuts  []apps.Shortcut `json:"shortcuts,omitempty" yaml:"shortcuts,omitempty"`
}

type SearchOptions struct {
	Query         string   `json:"query" yaml:"query"`
	AppName       string   `json:"app_name" yaml:"app_name"`
	Category      string   `json:"category" yaml:"category"`
	Tags          []string `json:"tags" yaml:"tags"`
	OnlyFavorites bool     `json:"only_favorites" yaml:"only_favorites"`
	SortBy        string   `json:"sort_by" yaml:"sort_by"`
	Limit         int      `json:"limit" yaml:"limit"`
	Offset        int      `json:"offset" yaml:"offset"`
}

type Manager interface {
	CreateNote(note *Note) error
	GetNote(id string) (*Note, error)
	UpdateNote(id string, note *Note) error
	DeleteNote(id string) error
	SearchNotes(opts SearchOptions) ([]*Note, error)
	ListNotes() ([]*Note, error)
	AddShortcutToNote(noteID string, shortcut apps.Shortcut) error
	RemoveShortcutFromNote(noteID string, shortcutIndex int) error
	ToggleFavorite(id string) error
	ExportNotes(format string) ([]byte, error)
	ImportNotes(data []byte, format string) error
}

type SyncStatus struct {
	LastSync     time.Time  `json:"last_sync" yaml:"last_sync"`
	TotalNotes   int        `json:"total_notes" yaml:"total_notes"`
	SyncedNotes  int        `json:"synced_notes" yaml:"synced_notes"`
	HasConflicts bool       `json:"has_conflicts" yaml:"has_conflicts"`
	Conflicts    []Conflict `json:"conflicts,omitempty" yaml:"conflicts,omitempty"`
}

type Conflict struct {
	NoteID     string    `json:"note_id" yaml:"note_id"`
	LocalNote  *Note     `json:"local_note" yaml:"local_note"`
	RemoteNote *Note     `json:"remote_note" yaml:"remote_note"`
	DetectedAt time.Time `json:"detected_at" yaml:"detected_at"`
}

type Resolution int

const (
	KeepLocal Resolution = iota
	KeepRemote
	Merge
)

type ConflictResolver interface {
	ResolveConflict(conflict Conflict, resolution Resolution) error
	AutoResolve(conflicts []Conflict) error
}
