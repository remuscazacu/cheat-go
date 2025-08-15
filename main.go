package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/cache"
	"cheat-go/pkg/config"
	"cheat-go/pkg/notes"
	"cheat-go/pkg/online"
	"cheat-go/pkg/plugins"
	"cheat-go/pkg/sync"
	"cheat-go/pkg/ui"
)

const (
	version = "v1.0.0-phase4"
	appName = "cheat-go"
)

// View modes for Phase 4 features
type viewMode int

const (
	viewMain viewMode = iota
	viewNotes
	viewPlugins
	viewOnline
	viewSync
	viewHelp
)

type model struct {
	// Original fields
	registry     *apps.Registry
	config       *config.Config
	renderer     *ui.TableRenderer
	rows         [][]string
	filteredRows [][]string
	cursorX      int
	cursorY      int
	searchMode   bool
	searchQuery  string
	lastSearch   string
	allRows      [][]string
	filterMode   bool
	filteredApps []string
	allApps      []string
	helpMode     bool

	// Phase 4 fields
	viewMode     viewMode
	cache        cache.Cache
	notesManager notes.Manager
	pluginLoader *plugins.Loader
	onlineClient online.Client
	syncManager  *sync.Manager

	// View-specific state
	notesList   []*notes.Note
	pluginsList []*plugins.LoadedPlugin
	reposList   []online.Repository
	cheatSheets []online.CheatSheet
	syncStatus  sync.SyncStatus

	// UI state for Phase 4 views
	noteCursor    int
	pluginCursor  int
	repoCursor    int
	sheetCursor   int
	statusMessage string
	loading       bool
}

type cliOptions struct {
	showHelp    bool
	showVersion bool
	theme       string
	tableStyle  string
	configFile  string
}

func printHelp() {
	fmt.Printf(`%s %s - Interactive terminal cheat sheet viewer with cloud features

USAGE:
    %s [OPTIONS]

DESCRIPTION:
    A fast, interactive terminal application for displaying keyboard shortcuts 
    and command cheat sheets. Navigate through shortcuts for popular applications
    with plugin support, personal notes, online repositories, and cloud sync.

OPTIONS:
    -h, --help              Show this help message and exit
    -v, --version           Show version information and exit
    -t, --theme THEME       Set the display theme
                            Options: default, dark, light, minimal
                            Default: default
    -s, --style STYLE       Set the table style  
                            Options: simple, rounded, bold, minimal
                            Default: simple
    -c, --config FILE       Use custom configuration file
                            Default: ~/.config/cheat-go/config.yaml

NAVIGATION:
    Arrow Keys / hjkl       Navigate through the table
    /                       Search mode
    f                       Filter apps
    n                       Open notes manager
    p                       Open plugin manager
    o                       Browse online repositories
    s                       Show sync status
    Ctrl+S                  Force sync
    ?                       Show help
    q / Ctrl+C              Quit the application

PHASE 4 FEATURES:
    Notes Manager (n)       Create and manage personal notes
    Plugin Manager (p)      Load and manage plugins
    Online Browser (o)      Browse community cheat sheets
    Sync Status (s)         View and manage cloud sync
    
THEMES:
    default                 Balanced colors for general use
    dark                    High-contrast for dark terminals
    light                   Clean appearance for light terminals  
    minimal                 Reduced visual elements

TABLE STYLES:
    simple                  Clean borders with lines
    rounded                 Elegant rounded corners
    bold                    Thick borders for visibility
    minimal                 Spacing-based separation

EXAMPLES:
    %s                      # Start with default settings
    %s --theme dark         # Use dark theme
    %s --style rounded      # Use rounded table borders
    %s -t dark -s bold      # Dark theme with bold borders
    %s --config my.yaml     # Use custom config file

For more information, visit: https://github.com/remuscazacu/cheat-go
`, appName, version, appName, appName, appName, appName, appName, appName)
}

func printVersion() {
	fmt.Printf("%s %s\n", appName, version)
}

func parseFlags() cliOptions {
	var opts cliOptions

	flag.BoolVar(&opts.showHelp, "h", false, "Show help message")
	flag.BoolVar(&opts.showHelp, "help", false, "Show help message")
	flag.BoolVar(&opts.showVersion, "v", false, "Show version")
	flag.BoolVar(&opts.showVersion, "version", false, "Show version")
	flag.StringVar(&opts.theme, "t", "", "Theme")
	flag.StringVar(&opts.theme, "theme", "", "Theme")
	flag.StringVar(&opts.tableStyle, "s", "", "Table style")
	flag.StringVar(&opts.tableStyle, "style", "", "Table style")
	flag.StringVar(&opts.configFile, "c", "", "Configuration file path")
	flag.StringVar(&opts.configFile, "config", "", "Configuration file path")

	flag.Parse()

	// Validate theme option
	if opts.theme != "" {
		validThemes := []string{"default", "dark", "light", "minimal"}
		valid := false
		for _, t := range validThemes {
			if opts.theme == t {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Error: Invalid theme '%s'. Valid options: %s\n",
				opts.theme, strings.Join(validThemes, ", "))
			os.Exit(1)
		}
	}

	// Validate table style option
	if opts.tableStyle != "" {
		validStyles := []string{"simple", "rounded", "bold", "minimal"}
		valid := false
		for _, s := range validStyles {
			if opts.tableStyle == s {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Error: Invalid table style '%s'. Valid options: %s\n",
				opts.tableStyle, strings.Join(validStyles, ", "))
			os.Exit(1)
		}
	}

	return opts
}

func initialModel(opts cliOptions) model {
	// Load configuration
	loader := config.NewLoader(opts.configFile)
	cfg, err := loader.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load config (%v), using defaults\n", err)
		cfg = config.DefaultConfig()
	}

	// Override config with CLI options
	if opts.theme != "" {
		cfg.Theme = opts.theme
	}
	if opts.tableStyle != "" {
		cfg.Layout.TableStyle = opts.tableStyle
	}

	// Initialize app registry
	registry := apps.NewRegistry(cfg.DataDir)
	if err := registry.LoadApps(cfg.Apps); err != nil {
		fmt.Printf("Warning: Could not load some apps (%v), using defaults\n", err)
	}

	// Create theme and renderer
	theme := ui.GetTheme(cfg.Theme)
	renderer := ui.NewTableRenderer(theme)
	renderer.SetTableStyle(cfg.Layout.TableStyle)
	renderer.SetMaxWidth(cfg.Layout.MaxWidth)

	// Generate table data
	rows := registry.GetTableData(cfg.Apps)

	// Initialize Phase 4 components
	m := model{
		registry:     registry,
		config:       cfg,
		renderer:     renderer,
		rows:         rows,
		filteredRows: rows,
		cursorX:      0,
		cursorY:      1,
		searchMode:   false,
		searchQuery:  "",
		lastSearch:   "",
		allRows:      rows,
		filterMode:   false,
		filteredApps: make([]string, 0),
		allApps:      cfg.Apps,
		helpMode:     false,
		viewMode:     viewMain,
	}

	// Initialize cache
	m.cache = cache.NewLRUCache(10*1024*1024, 1000) // 10MB, 1000 items

	// Initialize notes manager
	notesDir := os.ExpandEnv("$HOME/.config/cheat-go/notes")
	if cfg.DataDir != "" {
		notesDir = cfg.DataDir + "/notes"
	}
	m.notesManager, _ = notes.NewFileManager(notesDir)

	// Initialize plugin loader
	pluginDirs := []string{
		os.ExpandEnv("$HOME/.config/cheat-go/plugins"),
		"/usr/local/share/cheat-go/plugins",
	}
	if cfg.DataDir != "" {
		pluginDirs = append([]string{cfg.DataDir + "/plugins"}, pluginDirs...)
	}
	m.pluginLoader = plugins.NewLoader(pluginDirs...)
	m.pluginLoader.LoadAll()

	// Initialize online client (mock for now)
	m.onlineClient = online.NewMockClient()

	// Initialize sync manager (disabled by default)
	// m.syncManager would be initialized if sync is enabled in config

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle different view modes
		switch m.viewMode {
		case viewMain:
			if m.searchMode {
				return m.handleSearchInput(msg)
			}
			if m.filterMode {
				return m.handleFilterInput(msg)
			}
			if m.helpMode {
				return m.handleHelpInput(msg)
			}
			return m.handleMainInput(msg)
		case viewNotes:
			return m.handleNotesInput(msg)
		case viewPlugins:
			return m.handlePluginsInput(msg)
		case viewOnline:
			return m.handleOnlineInput(msg)
		case viewSync:
			return m.handleSyncInput(msg)
		case viewHelp:
			return m.handleHelpInput(msg)
		}
	}
	return m, nil
}

func (m model) handleMainInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "?":
		m.helpMode = true
		m.viewMode = viewHelp
		return m, nil
	case "/":
		m.searchMode = true
		return m, nil
	case "f", "ctrl+f":
		m.filterMode = true
		return m, nil
	case "n":
		// Open notes manager
		m.viewMode = viewNotes
		m.loadNotes()
		return m, nil
	case "p":
		// Open plugin manager
		m.viewMode = viewPlugins
		m.loadPlugins()
		return m, nil
	case "o":
		// Browse online repositories
		m.viewMode = viewOnline
		m.loadRepositories()
		return m, nil
	case "s":
		// Show sync status
		m.viewMode = viewSync
		m.loadSyncStatus()
		return m, nil
	case "ctrl+s":
		// Force sync
		m.statusMessage = "Syncing..."
		// Trigger sync in background
		return m, nil
	case "up", "k":
		if m.cursorY > 1 {
			m.cursorY--
		}
		return m, nil
	case "down", "j":
		if m.cursorY < len(m.rows)-1 {
			m.cursorY++
		}
		return m, nil
	case "left", "h":
		if m.cursorX > 0 {
			m.cursorX--
		}
		return m, nil
	case "right", "l":
		if m.cursorX < len(m.rows[0])-1 {
			m.cursorX++
		}
		return m, nil
	case "ctrl+a", "home":
		m.cursorY = 1
		return m, nil
	case "ctrl+e", "end":
		m.cursorY = len(m.rows) - 1
		return m, nil
	case "ctrl+r":
		// Refresh data
		m.rows = m.registry.GetTableData(m.config.Apps)
		m.allRows = m.rows
		return m, nil
	case "esc", "ctrl+[":
		// Clear search/filter
		m.searchMode = false
		m.searchQuery = ""
		m.lastSearch = ""
		m.rows = m.allRows
		m.cursorY = 1
		return m, nil
	}
	return m, nil
}

func (m model) handleNotesInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewMode = viewMain
		return m, nil
	case "up", "k":
		if m.noteCursor > 0 {
			m.noteCursor--
		}
		return m, nil
	case "down", "j":
		if m.noteCursor < len(m.notesList)-1 {
			m.noteCursor++
		}
		return m, nil
	case "n":
		// Create new note with simple hardcoded example
		newNote := &notes.Note{
			Title:    fmt.Sprintf("New Note %d", time.Now().Unix()),
			Content:  "Enter your note content here",
			Category: "general",
			Tags:     []string{"new"},
		}
		err := m.notesManager.CreateNote(newNote)
		if err != nil {
			m.statusMessage = fmt.Sprintf("Error creating note: %v", err)
		} else {
			m.loadNotes()
			m.statusMessage = "Note created successfully"
		}
		return m, nil
	case "e":
		// Edit selected note in default editor
		if m.noteCursor < len(m.notesList) {
			note := m.notesList[m.noteCursor]
			updatedNote, err := m.openEditorForNote(note)
			if err != nil {
				m.statusMessage = fmt.Sprintf("Error opening editor: %v", err)
			} else if updatedNote != nil {
				err := m.notesManager.UpdateNote(note.ID, updatedNote)
				if err != nil {
					m.statusMessage = fmt.Sprintf("Error updating note: %v", err)
				} else {
					m.loadNotes()
					m.statusMessage = fmt.Sprintf("Note '%s' updated", updatedNote.Title)
				}
			}
		}
		return m, nil
	case "d":
		// Delete selected note
		if m.noteCursor < len(m.notesList) {
			noteID := m.notesList[m.noteCursor].ID
			m.notesManager.DeleteNote(noteID)
			m.loadNotes()
			m.statusMessage = "Note deleted"
		}
		return m, nil
	case "f":
		// Toggle favorite
		if m.noteCursor < len(m.notesList) {
			noteID := m.notesList[m.noteCursor].ID
			m.notesManager.ToggleFavorite(noteID)
			m.loadNotes()
		}
		return m, nil
	}
	return m, nil
}

func (m model) handlePluginsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewMode = viewMain
		return m, nil
	case "up", "k":
		if m.pluginCursor > 0 {
			m.pluginCursor--
		}
		return m, nil
	case "down", "j":
		if m.pluginCursor < len(m.pluginsList)-1 {
			m.pluginCursor++
		}
		return m, nil
	case "l":
		// Load plugin
		if m.pluginCursor < len(m.pluginsList) {
			plugin := m.pluginsList[m.pluginCursor]
			m.statusMessage = fmt.Sprintf("Loading plugin: %s", plugin.Metadata.Name)
		}
		return m, nil
	case "u":
		// Unload plugin
		if m.pluginCursor < len(m.pluginsList) {
			plugin := m.pluginsList[m.pluginCursor]
			m.pluginLoader.UnloadPlugin(plugin.Metadata.Name)
			m.loadPlugins()
			m.statusMessage = fmt.Sprintf("Unloaded plugin: %s", plugin.Metadata.Name)
		}
		return m, nil
	case "r":
		// Reload plugins
		m.pluginLoader.LoadAll()
		m.loadPlugins()
		m.statusMessage = "Plugins reloaded"
		return m, nil
	}
	return m, nil
}

func (m model) handleOnlineInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewMode = viewMain
		return m, nil
	case "up", "k":
		if m.repoCursor > 0 {
			m.repoCursor--
		}
		return m, nil
	case "down", "j":
		if m.repoCursor < len(m.reposList)-1 {
			m.repoCursor++
		}
		return m, nil
	case "enter":
		// Browse selected repository
		if m.repoCursor < len(m.reposList) {
			repo := m.reposList[m.repoCursor]
			m.loadCheatSheets(repo.URL)
		}
		return m, nil
	case "d":
		// Download selected cheat sheet
		if m.sheetCursor < len(m.cheatSheets) {
			sheet := m.cheatSheets[m.sheetCursor]
			m.statusMessage = fmt.Sprintf("Downloading: %s", sheet.Name)
			// Download logic here
		}
		return m, nil
	case "/":
		// Search online
		m.searchMode = true
		return m, nil
	}
	return m, nil
}

func (m model) handleSyncInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewMode = viewMain
		return m, nil
	case "s":
		// Trigger sync
		m.statusMessage = "Syncing..."
		if m.syncManager != nil {
			go m.syncManager.Sync()
		}
		return m, nil
	case "r":
		// Resolve conflicts
		m.statusMessage = "Resolving conflicts..."
		return m, nil
	case "a":
		// Toggle auto-sync
		if m.syncManager != nil {
			// Toggle auto-sync setting
			m.statusMessage = "Auto-sync toggled"
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "ctrl+[":
		m.searchMode = false
		m.searchQuery = ""
		m.rows = m.allRows
		m.lastSearch = ""
		m.cursorY = 1
		return m, nil
	case "ctrl+u":
		m.searchQuery = ""
		return m, nil
	case "enter":
		m.searchMode = false
		if m.searchQuery == "" {
			m.rows = m.allRows
			m.lastSearch = ""
		} else {
			m.rows = m.registry.SearchTableData(m.config.Apps, m.searchQuery)
			m.lastSearch = m.searchQuery
		}
		m.cursorY = 1
		return m, nil
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
		return m, nil
	}
}

func (m model) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "ctrl+[":
		// Exit filter mode
		m.filterMode = false
		return m, nil
	case "ctrl+u":
		// Clear all selections
		m.filteredApps = make([]string, 0)
		return m, nil
	case "enter":
		// Apply filter and exit filter mode
		m.filterMode = false
		if len(m.filteredApps) == 0 {
			// If no apps selected, show all
			m.rows = m.registry.GetTableData(m.allApps)
		} else {
			// Show only selected apps
			m.rows = m.registry.GetTableData(m.filteredApps)
		}
		m.cursorY = 1
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Toggle app selection by number
		appIndex := int(msg.String()[0] - '1') // Convert '1' to 0, '2' to 1, etc.
		if appIndex < len(m.allApps) {
			appName := m.allApps[appIndex]
			// Toggle app in filteredApps
			found := false
			for i, name := range m.filteredApps {
				if name == appName {
					// Remove app
					m.filteredApps = append(m.filteredApps[:i], m.filteredApps[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				// Add app
				m.filteredApps = append(m.filteredApps, appName)
			}
		}
		return m, nil
	case "c":
		// Clear all filters
		m.filteredApps = make([]string, 0)
		return m, nil
	case "a":
		// Select all apps
		m.filteredApps = make([]string, len(m.allApps))
		copy(m.filteredApps, m.allApps)
		return m, nil
	}
	return m, nil
}

func (m model) handleHelpInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?", "esc", "q":
		m.helpMode = false
		m.viewMode = viewMain
		return m, nil
	}
	return m, nil
}

// Helper methods for loading data
func (m *model) loadNotes() {
	notes, _ := m.notesManager.ListNotes()
	m.notesList = notes
	m.noteCursor = 0
}

func (m *model) loadPlugins() {
	m.pluginsList = m.pluginLoader.ListPlugins()
	m.pluginCursor = 0
}

func (m *model) loadRepositories() {
	repos, _ := m.onlineClient.GetRepositories()
	m.reposList = repos
	m.repoCursor = 0
}

func (m *model) loadCheatSheets(repoURL string) {
	sheets, _ := m.onlineClient.SearchCheatSheets(online.SearchOptions{
		Repository: repoURL,
		Limit:      50,
	})
	m.cheatSheets = sheets
	m.sheetCursor = 0
}

func (m *model) loadSyncStatus() {
	if m.syncManager != nil {
		m.syncStatus = m.syncManager.GetSyncStatus()
	} else {
		m.syncStatus = sync.SyncStatus{
			LastSync:  time.Time{},
			IsSyncing: false,
			DeviceID:  "not-configured",
		}
	}
}

func (m model) View() string {
	switch m.viewMode {
	case viewNotes:
		return m.viewNotes()
	case viewPlugins:
		return m.viewPlugins()
	case viewOnline:
		return m.viewOnline()
	case viewSync:
		return m.viewSync()
	case viewHelp:
		return m.viewHelp()
	default:
		return m.viewMain()
	}
}

func (m model) viewMain() string {
	// Original main view implementation
	var output strings.Builder

	// Render the table with cursor position
	tableStr := m.renderer.RenderWithHighlighting(
		m.rows,
		m.cursorX,
		m.cursorY,
		m.lastSearch,
	)
	output.WriteString(tableStr)
	output.WriteString("\n")

	// Show status bar
	if m.searchMode {
		output.WriteString(fmt.Sprintf("\nSearch: %s_\nType to search, Enter to confirm, Esc to cancel\n", m.searchQuery))
	} else if m.filterMode {
		output.WriteString("\nFilter Apps: ")
		for i, app := range m.allApps {
			isSelected := false
			for _, selected := range m.filteredApps {
				if app == selected {
					isSelected = true
					break
				}
			}
			if isSelected {
				output.WriteString(fmt.Sprintf(" [%d] ✓%s", i+1, app))
			} else {
				output.WriteString(fmt.Sprintf(" [%d] %s", i+1, app))
			}
		}
		output.WriteString("\n1-9: toggle apps, a: all, c: clear, Enter: apply, Esc: cancel\n")
	} else {
		output.WriteString("\nArrow keys/hjkl: move • /: search • f: filter • n: notes • p: plugins • o: online • s: sync • ?: help • q: quit\n")
	}

	if m.statusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.statusMessage))
	}

	return output.String()
}

func (m model) viewNotes() string {
	var output strings.Builder

	output.WriteString("╭─ Personal Notes ─────────────────────────────────────────╮\n")

	if len(m.notesList) == 0 {
		output.WriteString("│  No notes found. Press 'n' to create a new note.        │\n")
	} else {
		for i, note := range m.notesList {
			if i > 10 {
				output.WriteString(fmt.Sprintf("│  ... and %d more notes                                   │\n", len(m.notesList)-10))
				break
			}

			cursor := "  "
			if i == m.noteCursor {
				cursor = "▶ "
			}

			favorite := " "
			if note.IsFavorite {
				favorite = "⭐"
			}

			line := fmt.Sprintf("%s%s %-30s %s", cursor, favorite, note.Title, note.AppName)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: n: new • e: edit • d: delete • f: favorite • esc: back\n")

	if m.statusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.statusMessage))
	}

	return output.String()
}

func (m model) viewPlugins() string {
	var output strings.Builder

	output.WriteString("╭─ Plugin Manager ─────────────────────────────────────────╮\n")

	if len(m.pluginsList) == 0 {
		output.WriteString("│  No plugins loaded. Place plugins in plugins directory. │\n")
	} else {
		for i, plugin := range m.pluginsList {
			if i > 10 {
				output.WriteString(fmt.Sprintf("│  ... and %d more plugins                                 │\n", len(m.pluginsList)-10))
				break
			}

			cursor := "  "
			if i == m.pluginCursor {
				cursor = "▶ "
			}

			line := fmt.Sprintf("%s%-20s v%-8s %s", cursor, plugin.Metadata.Name, plugin.Metadata.Version, plugin.Metadata.Author)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: l: load • u: unload • r: reload all • esc: back\n")

	if m.statusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.statusMessage))
	}

	return output.String()
}

func (m model) viewOnline() string {
	var output strings.Builder

	output.WriteString("╭─ Online Repositories ────────────────────────────────────╮\n")

	if len(m.reposList) == 0 {
		output.WriteString("│  Loading repositories...                                 │\n")
	} else {
		for i, repo := range m.reposList {
			cursor := "  "
			if i == m.repoCursor {
				cursor = "▶ "
			}

			line := fmt.Sprintf("%s%-30s ⭐%d", cursor, repo.Name, repo.Stars)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	if len(m.cheatSheets) > 0 {
		output.WriteString("│──────────────────────────────────────────────────────────│\n")
		output.WriteString("│ Cheat Sheets:                                            │\n")
		for i, sheet := range m.cheatSheets {
			if i > 5 {
				break
			}
			line := fmt.Sprintf("  %-25s ⬇%d ★%.1f", sheet.Name, sheet.Downloads, sheet.Rating)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: enter: browse • d: download • /: search • esc: back\n")

	if m.statusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.statusMessage))
	}

	return output.String()
}

func (m model) viewSync() string {
	var output strings.Builder

	output.WriteString("╭─ Sync Status ────────────────────────────────────────────╮\n")

	if m.syncManager == nil {
		output.WriteString("│  Sync is not configured.                                 │\n")
		output.WriteString("│  Configure sync in ~/.config/cheat-go/config.yaml        │\n")
	} else {
		status := "Idle"
		if m.syncStatus.IsSyncing {
			status = "Syncing..."
		}

		lastSync := "Never"
		if !m.syncStatus.LastSync.IsZero() {
			lastSync = m.syncStatus.LastSync.Format("2006-01-02 15:04:05")
		}

		output.WriteString(fmt.Sprintf("│  Status:     %-43s │\n", status))
		output.WriteString(fmt.Sprintf("│  Last Sync:  %-43s │\n", lastSync))
		output.WriteString(fmt.Sprintf("│  Device ID:  %-43s │\n", m.syncStatus.DeviceID[:16]+"..."))

		if m.syncStatus.HasConflicts {
			output.WriteString(fmt.Sprintf("│  ⚠ Conflicts: %-42d │\n", len(m.syncStatus.Conflicts)))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: s: sync now • r: resolve conflicts • a: auto-sync • esc: back\n")

	if m.statusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.statusMessage))
	}

	return output.String()
}

func (m model) viewHelp() string {
	return `
╭─ Help ───────────────────────────────────────────────────╮
│                                                          │
│  NAVIGATION                                              │
│    ↑/k, ↓/j, ←/h, →/l  Navigate                        │
│    Home/Ctrl+A          Go to first row                 │
│    End/Ctrl+E           Go to last row                  │
│                                                          │
│  FEATURES                                                │
│    /                    Search mode                     │
│    f                    Filter apps                     │
│    n                    Notes manager                   │
│    p                    Plugin manager                  │
│    o                    Browse online                   │
│    s                    Sync status                     │
│    Ctrl+S               Force sync                      │
│    Ctrl+R               Refresh data                    │
│    ?                    This help screen                │
│    q/Ctrl+C             Quit                           │
│                                                          │
│  SEARCH MODE                                            │
│    Type to search                                       │
│    Enter                Confirm search                  │
│    Esc                  Cancel search                   │
│    Ctrl+U               Clear search                    │
│                                                          │
╰──────────────────────────────────────────────────────────╯

Press ? or Esc to close help
`
}

func (m model) filterRowsBySearch(query string) [][]string {
	if query == "" {
		return m.allRows
	}

	var filtered [][]string
	// Always include header row
	if len(m.allRows) > 0 {
		filtered = append(filtered, m.allRows[0])
	}

	// Filter data rows (skip header at index 0)
	for i := 1; i < len(m.allRows); i++ {
		row := m.allRows[i]
		// Search in all columns of the row
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), strings.ToLower(query)) {
				filtered = append(filtered, row)
				break // Found match, add row and move to next
			}
		}
	}

	return filtered
}

func (m model) isAppSelected(appName string) bool {
	for _, name := range m.filteredApps {
		if name == appName {
			return true
		}
	}
	return false
}

// openEditorForNote opens the default editor to edit a note's content
func (m model) openEditorForNote(note *notes.Note) (*notes.Note, error) {
	// Get the default editor from environment variable, default to nano
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	// Create a temporary file with the note content
	tmpFile, err := ioutil.TempFile("", "cheat-go-note-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write current note content to temp file
	content := fmt.Sprintf("# Title: %s\n# Category: %s\n# Tags: %s\n\n%s",
		note.Title, note.Category, strings.Join(note.Tags, ", "), note.Content)

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Open the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor exited with error: %v", err)
	}

	// Read the edited content back
	editedContent, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read edited content: %v", err)
	}

	// Parse the edited content to extract title, category, tags, and content
	lines := strings.Split(string(editedContent), "\n")
	updatedNote := &notes.Note{
		ID:       note.ID,
		Title:    note.Title,
		Category: note.Category,
		Tags:     note.Tags,
		Content:  note.Content,
	}

	var contentStart int
	for i, line := range lines {
		if strings.HasPrefix(line, "# Title: ") {
			updatedNote.Title = strings.TrimSpace(strings.TrimPrefix(line, "# Title: "))
		} else if strings.HasPrefix(line, "# Category: ") {
			updatedNote.Category = strings.TrimSpace(strings.TrimPrefix(line, "# Category: "))
		} else if strings.HasPrefix(line, "# Tags: ") {
			tagStr := strings.TrimSpace(strings.TrimPrefix(line, "# Tags: "))
			if tagStr != "" {
				updatedNote.Tags = strings.Split(tagStr, ", ")
				for j := range updatedNote.Tags {
					updatedNote.Tags[j] = strings.TrimSpace(updatedNote.Tags[j])
				}
			} else {
				updatedNote.Tags = []string{}
			}
		} else if line == "" && i > 0 && strings.HasPrefix(lines[i-1], "# Tags: ") {
			contentStart = i + 1
			break
		}
	}

	if contentStart < len(lines) {
		updatedNote.Content = strings.Join(lines[contentStart:], "\n")
	}

	return updatedNote, nil
}

func main() {
	opts := parseFlags()

	if opts.showHelp {
		printHelp()
		os.Exit(0)
	}

	if opts.showVersion {
		printVersion()
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel(opts))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
