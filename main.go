package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/config"
	"cheat-go/pkg/ui"
)

type model struct {
	registry     *apps.Registry
	config       *config.Config
	renderer     *ui.TableRenderer
	rows         [][]string
	filteredRows [][]string
	cursorX      int
	cursorY      int
	searchMode   bool
	searchQuery  string
	lastSearch   string     // Last applied search term for highlighting
	allRows      [][]string // Original unfiltered data
	filterMode   bool
	filteredApps []string // Apps to show (empty = show all)
	allApps      []string // All available apps
	helpMode     bool     // Show help screen
}

func initialModel() model {
	// Load configuration
	loader := config.NewLoader("")
	cfg, err := loader.Load()
	if err != nil {
		// Log error but continue with defaults
		fmt.Printf("Warning: Could not load config (%v), using defaults\n", err)
		cfg = config.DefaultConfig()
	}

	// Initialize app registry
	registry := apps.NewRegistry(cfg.DataDir)
	if err := registry.LoadApps(cfg.Apps); err != nil {
		// Log warning but continue with hardcoded data
		fmt.Printf("Warning: Could not load some apps (%v), using defaults\n", err)
	}

	// Create theme and renderer
	theme := ui.GetTheme(cfg.Theme)
	renderer := ui.NewTableRenderer(theme)

	// Generate table data
	rows := registry.GetTableData(cfg.Apps)

	return model{
		registry:     registry,
		config:       cfg,
		renderer:     renderer,
		rows:         rows,
		filteredRows: rows, // Initially, filtered rows same as all rows
		cursorX:      0,
		cursorY:      1,
		searchMode:   false,
		searchQuery:  "",
		lastSearch:   "",
		allRows:      rows, // Store original data
		filterMode:   false,
		filteredApps: make([]string, 0),
		allApps:      cfg.Apps,
		helpMode:     false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}
		if m.filterMode {
			return m.handleFilterInput(msg)
		}
		if m.helpMode {
			return m.handleHelpInput(msg)
		}
		return m.handleNormalInput(msg)
	}
	return m, nil
}

func (m model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "ctrl+[":
		// Exit search mode
		m.searchMode = false
		m.searchQuery = ""
		m.rows = m.allRows
		m.lastSearch = ""
		m.cursorY = 1
		return m, nil
	case "ctrl+u":
		// Clear search query
		m.searchQuery = ""
		return m, nil
	case "enter":
		// Confirm search and exit search mode
		m.searchMode = false
		if m.searchQuery == "" {
			// If empty search, show all results
			m.rows = m.allRows
			m.lastSearch = ""
		} else {
			// Apply search filter using registry
			m.rows = m.registry.SearchTableData(m.config.Apps, m.searchQuery)
			m.lastSearch = m.searchQuery
		}
		m.cursorY = 1
		return m, nil
	case "backspace":
		// Remove last character
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		return m, nil
	default:
		// Add character to search query
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
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "ctrl+[", "?":
		// Exit help mode
		m.helpMode = false
		return m, nil
	}
	return m, nil
}

func (m model) handleNormalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "/":
		// Enter search mode
		m.searchMode = true
		m.searchQuery = ""
		return m, nil
	case "f":
		// Enter filter mode
		m.filterMode = true
		return m, nil
	case "?":
		// Show help
		m.helpMode = true
		return m, nil
	case "esc", "ctrl+[":
		// Clear search filter if any
		if len(m.rows) != len(m.allRows) {
			m.rows = m.allRows
			m.lastSearch = ""
			m.cursorY = 1
		}
		return m, nil
	case "ctrl+f":
		// Alternative way to enter filter mode
		m.filterMode = true
		return m, nil
	case "ctrl+r":
		// Refresh/reload data
		m.rows = m.registry.GetTableData(m.config.Apps)
		m.allRows = m.rows
		m.cursorY = 1
		return m, nil
	case "home", "ctrl+a":
		// Go to first data row
		m.cursorY = 1
		return m, nil
	case "end", "ctrl+e":
		// Go to last row
		if len(m.rows) > 1 {
			m.cursorY = len(m.rows) - 1
		}
		return m, nil
	case "up", "k":
		if m.cursorY > 1 {
			m.cursorY--
		}
	case "down", "j":
		if m.cursorY < len(m.rows)-1 {
			m.cursorY++
		}
	case "left", "h":
		if m.cursorX > 0 {
			m.cursorX--
		}
	case "right", "l":
		if m.cursorX < len(m.rows[0])-1 {
			m.cursorX++
		}
	}
	return m, nil
}

// filterRowsBySearch filters table rows based on search query
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
			if containsIgnoreCase(cell, query) {
				filtered = append(filtered, row)
				break // Found match, add row and move to next
			}
		}
	}

	return filtered
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// isAppSelected checks if an app is in the filtered apps list
func (m model) isAppSelected(appName string) bool {
	for _, name := range m.filteredApps {
		if name == appName {
			return true
		}
	}
	return false
}

// getInstructions returns the appropriate instructions for the current state
func (m model) getInstructions() string {
	return "Arrow keys/hjkl: move • /: search • f: filter • Ctrl+R: refresh • ?: help • q: quit"
}

// renderHelp returns the help screen content
func (m model) renderHelp() string {
	theme := m.renderer.GetTheme()
	title := theme.HeaderStyle.Render("CHEAT-GO KEYBOARD SHORTCUTS")

	helpText := `
NAVIGATION:
  ↑/k           Move cursor up
  ↓/j           Move cursor down  
  ←/h           Move cursor left
  →/l           Move cursor right
  Home/Ctrl+A   Go to first row
  End/Ctrl+E    Go to last row

SEARCH:
  /             Enter search mode
  Esc           Exit search mode / clear filters
  Enter         Confirm search
  Backspace     Delete character in search
  Ctrl+U        Clear search query

FILTERING:
  f/Ctrl+F      Enter filter mode
  1-9           Toggle app selection
  a             Select all apps
  c             Clear all selections
  Ctrl+U        Clear all selections
  Enter         Apply filter
  Esc           Cancel filter

GENERAL:
  Ctrl+R        Refresh data
  ?             Show/hide this help
  q/Ctrl+C      Quit application

Press ? or Esc to close this help screen`

	return title + helpText
}

func (m model) View() string {
	var table string
	// Use highlighting if we have an active search term
	if m.lastSearch != "" && len(m.rows) != len(m.allRows) {
		table = m.renderer.RenderWithHighlighting(m.rows, m.cursorX, m.cursorY, m.lastSearch)
	} else {
		table = m.renderer.Render(m.rows, m.cursorX, m.cursorY)
	}

	if m.searchMode {
		// Style the search UI
		theme := m.renderer.GetTheme()
		searchPrompt := theme.SearchStyle.Render("Search: ")
		searchInput := theme.SearchInputStyle.Render(m.searchQuery + "_")
		searchLine := searchPrompt + searchInput

		instructions := "Type to search, Enter to confirm, Esc to cancel"
		return table + "\n" + searchLine + "\n" + instructions
	}

	if m.filterMode {
		// Show filter UI
		theme := m.renderer.GetTheme()
		filterPrompt := theme.SearchStyle.Render("Filter Apps: ")

		var appList strings.Builder
		for i, app := range m.allApps {
			selected := m.isAppSelected(app)
			marker := " "
			if selected {
				marker = "✓"
			}
			appList.WriteString(fmt.Sprintf(" [%d] %s%s", i+1, marker, app))
		}

		instructions := "1-9: toggle apps, a: all, c: clear, Enter: apply, Esc: cancel"
		return table + "\n" + filterPrompt + appList.String() + "\n" + instructions
	}

	if m.helpMode {
		return m.renderHelp()
	}

	instructions := m.getInstructions()
	if len(m.rows) != len(m.allRows) {
		// Show that results are filtered
		resultCount := len(m.rows) - 1 // Exclude header
		totalCount := len(m.allRows) - 1
		filterInfo := fmt.Sprintf(" (%d/%d results)", resultCount, totalCount)
		instructions = "Filtered results" + filterInfo + " - Press / to search, f to filter, Esc to clear, ? for help, q to quit."
	}

	return table + "\n" + instructions
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
