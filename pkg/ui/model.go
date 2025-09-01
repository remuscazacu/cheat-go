package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/cache"
	"cheat-go/pkg/config"
	"cheat-go/pkg/notes"
	"cheat-go/pkg/online"
	"cheat-go/pkg/plugins"
	"cheat-go/pkg/sync"
)

// View modes for Phase 4 features
type ViewMode int

const (
	ViewMain ViewMode = iota
	ViewNotes
	ViewPlugins
	ViewOnline
	ViewSync
	ViewHelp
)

// Aliases for backward compatibility with lowercase names used in main.go
const (
	viewMain    = ViewMain
	viewNotes   = ViewNotes
	viewPlugins = ViewPlugins
	viewOnline  = ViewOnline
	viewSync    = ViewSync
	viewHelp    = ViewHelp
)

type Model struct {
	// Original fields
	Registry     *apps.Registry
	Config       *config.Config
	Renderer     *TableRenderer
	Rows         [][]string
	FilteredRows [][]string
	CursorX      int
	CursorY      int
	SearchMode   bool
	SearchQuery  string
	LastSearch   string
	AllRows      [][]string
	FilterMode   bool
	FilteredApps []string
	AllApps      []string
	HelpMode     bool

	// Phase 4 fields
	ViewMode     ViewMode
	Cache        cache.Cache
	NotesManager notes.Manager
	PluginLoader *plugins.Loader
	OnlineClient online.Client
	SyncManager  *sync.Manager

	// View-specific state
	NotesList   []*notes.Note
	PluginsList []*plugins.LoadedPlugin
	ReposList   []online.Repository
	CheatSheets []online.CheatSheet
	SyncStatus  sync.SyncStatus

	// UI state for Phase 4 views
	NoteCursor    int
	PluginCursor  int
	RepoCursor    int
	SheetCursor   int
	StatusMessage string
	Loading       bool
}

func NewModel() Model {
	return Model{
		ViewMode: ViewMain,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.ViewMode {
		case ViewMain:
			if m.SearchMode {
				return m.HandleSearchInput(msg)
			}
			if m.FilterMode {
				return m.HandleFilterInput(msg)
			}
			if m.HelpMode {
				return m.HandleHelpInput(msg)
			}
			return m.HandleMainInput(msg)
		case ViewNotes:
			return m.HandleNotesInput(msg)
		case ViewPlugins:
			return m.HandlePluginsInput(msg)
		case ViewOnline:
			return m.HandleOnlineInput(msg)
		case ViewSync:
			return m.HandleSyncInput(msg)
		case ViewHelp:
			return m.HandleHelpInput(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.ViewMode {
	case ViewNotes:
		return m.ViewNotes()
	case ViewPlugins:
		return m.ViewPlugins()
	case ViewOnline:
		return m.ViewOnline()
	case ViewSync:
		return m.ViewSync()
	case ViewHelp:
		return m.ViewHelp()
	default:
		return m.ViewMain()
	}
}
