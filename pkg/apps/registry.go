package apps

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrAppNotFound    = errors.New("application not found")
	ErrInvalidAppFile = errors.New("invalid app file format")
)

// Registry manages application loading and registration
type Registry struct {
	*AppRegistry
	dataDir string
}

// NewRegistry creates a new registry with default hardcoded apps
func NewRegistry(dataDir string) *Registry {
	registry := &Registry{
		AppRegistry: NewAppRegistry(),
		dataDir:     dataDir,
	}

	// Load hardcoded apps as fallback
	registry.loadHardcodedApps()

	return registry
}

// LoadApps loads applications from configuration
func (r *Registry) LoadApps(appNames []string) error {
	for _, name := range appNames {
		if err := r.LoadApp(name); err != nil {
			// If loading from file fails, use hardcoded data
			continue
		}
	}
	return nil
}

// LoadApp loads a single application from file or hardcoded data
func (r *Registry) LoadApp(name string) error {
	// Try to load from file first
	if r.dataDir != "" {
		appPath := filepath.Join(r.dataDir, name+".yaml")
		if app, err := r.loadAppFromFile(appPath); err == nil {
			r.Register(app)
			return nil
		}
	}

	// If file loading fails, app should already be loaded from hardcoded data
	if _, exists := r.Get(name); exists {
		return nil
	}

	return ErrAppNotFound
}

// loadAppFromFile loads an app definition from a YAML file
func (r *Registry) loadAppFromFile(path string) (*App, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var app App
	if err := yaml.Unmarshal(data, &app); err != nil {
		return nil, ErrInvalidAppFile
	}

	return &app, nil
}

// loadHardcodedApps loads the original hardcoded application data
func (r *Registry) loadHardcodedApps() {
	// Convert original table data to structured format
	shortcuts := []struct {
		key     string
		vim     string
		zsh     string
		dwm     string
		st      string
		lf      string
		zathura string
	}{
		{"h", "← move", "back char", "focus left", "← move", "left", "scroll ←"},
		{"l", "→ move", "forward char", "focus right", "→ move", "right", "scroll →"},
		{"j", "↓ move", "down history", "focus down", "↓ scroll", "down", "scroll ↓"},
		{"k", "↑ move", "up history", "focus up", "↑ scroll", "up", "scroll ↑"},
		{"gg", "top", "-", "-", "-", "top", "-"},
		{"G", "bottom", "-", "-", "-", "bottom", "-"},
		{"/", "search", "search history", "-", "search", "search", "search"},
		{":", "command", "prompt", "command", "-", "command", "-"},
		{"q", "quit", "exit", "close win", "exit", "quit", "quit"},
	}

	// Create apps from hardcoded data
	apps := map[string]*App{
		"vim": {
			Name:        "vim",
			Description: "Vi IMproved text editor",
			Version:     "1.0",
			Categories:  []string{"editor"},
			Shortcuts:   []Shortcut{},
		},
		"zsh": {
			Name:        "zsh",
			Description: "Z Shell",
			Version:     "1.0",
			Categories:  []string{"shell"},
			Shortcuts:   []Shortcut{},
		},
		"dwm": {
			Name:        "dwm",
			Description: "Dynamic window manager",
			Version:     "1.0",
			Categories:  []string{"wm"},
			Shortcuts:   []Shortcut{},
		},
		"st": {
			Name:        "st",
			Description: "Simple terminal",
			Version:     "1.0",
			Categories:  []string{"terminal"},
			Shortcuts:   []Shortcut{},
		},
		"lf": {
			Name:        "lf",
			Description: "Terminal file manager",
			Version:     "1.0",
			Categories:  []string{"file-manager"},
			Shortcuts:   []Shortcut{},
		},
		"zathura": {
			Name:        "zathura",
			Description: "Document viewer",
			Version:     "1.0",
			Categories:  []string{"viewer"},
			Shortcuts:   []Shortcut{},
		},
	}

	// Populate shortcuts for each app
	for _, shortcut := range shortcuts {
		descriptions := []string{shortcut.vim, shortcut.zsh, shortcut.dwm, shortcut.st, shortcut.lf, shortcut.zathura}
		appNames := []string{"vim", "zsh", "dwm", "st", "lf", "zathura"}

		for i, desc := range descriptions {
			if desc != "-" && desc != "" {
				apps[appNames[i]].Shortcuts = append(apps[appNames[i]].Shortcuts, Shortcut{
					Keys:        shortcut.key,
					Description: desc,
					Category:    "general",
				})
			}
		}
	}

	// Register all apps
	for _, app := range apps {
		r.Register(app)
	}
}

// GetTableData returns data in the original table format for backward compatibility
func (r *Registry) GetTableData(appNames []string) [][]string {
	// Header row
	header := make([]string, len(appNames)+1)
	header[0] = "Shortcut"
	copy(header[1:], appNames)

	// Collect all unique shortcuts
	shortcutMap := make(map[string][]string)
	for i, appName := range appNames {
		if app, exists := r.Get(appName); exists {
			for _, shortcut := range app.Shortcuts {
				if _, exists := shortcutMap[shortcut.Keys]; !exists {
					shortcutMap[shortcut.Keys] = make([]string, len(appNames))
					for j := range shortcutMap[shortcut.Keys] {
						shortcutMap[shortcut.Keys][j] = "-"
					}
				}
				shortcutMap[shortcut.Keys][i] = shortcut.Description
			}
		}
	}

	// Convert to table format
	rows := [][]string{header}
	for keys, descriptions := range shortcutMap {
		row := make([]string, len(appNames)+1)
		row[0] = keys
		copy(row[1:], descriptions)
		rows = append(rows, row)
	}

	return rows
}

// SearchTableData returns filtered table data based on search query
func (r *Registry) SearchTableData(appNames []string, query string) [][]string {
	// If no query, return all data
	if query == "" {
		return r.GetTableData(appNames)
	}

	// Header row
	header := make([]string, len(appNames)+1)
	header[0] = "Shortcut"
	copy(header[1:], appNames)

	// Collect shortcuts that match the search query
	shortcutMap := make(map[string][]string)
	for i, appName := range appNames {
		if app, exists := r.Get(appName); exists {
			for _, shortcut := range app.Shortcuts {
				// Search in keys, description, and category
				if r.shortcutMatches(shortcut, query) {
					if _, exists := shortcutMap[shortcut.Keys]; !exists {
						shortcutMap[shortcut.Keys] = make([]string, len(appNames))
						for j := range shortcutMap[shortcut.Keys] {
							shortcutMap[shortcut.Keys][j] = "-"
						}
					}
					shortcutMap[shortcut.Keys][i] = shortcut.Description
				}
			}
		}
	}

	// Convert to table format
	rows := [][]string{header}
	for keys, descriptions := range shortcutMap {
		row := make([]string, len(appNames)+1)
		row[0] = keys
		copy(row[1:], descriptions)
		rows = append(rows, row)
	}

	return rows
}

// shortcutMatches checks if a shortcut matches the search query
func (r *Registry) shortcutMatches(shortcut Shortcut, query string) bool {
	queryLower := strings.ToLower(query)

	// Search in keys
	if strings.Contains(strings.ToLower(shortcut.Keys), queryLower) {
		return true
	}

	// Search in description
	if strings.Contains(strings.ToLower(shortcut.Description), queryLower) {
		return true
	}

	// Search in category
	if strings.Contains(strings.ToLower(shortcut.Category), queryLower) {
		return true
	}

	return false
}

// SearchShortcuts returns all shortcuts matching the query across all apps
func (r *Registry) SearchShortcuts(query string) []ShortcutResult {
	var results []ShortcutResult
	queryLower := strings.ToLower(query)

	for appName, app := range r.AppRegistry.apps {
		for _, shortcut := range app.Shortcuts {
			if r.shortcutMatches(shortcut, query) {
				results = append(results, ShortcutResult{
					AppName:  appName,
					Shortcut: shortcut,
					Matches:  r.getSearchMatches(shortcut, queryLower),
				})
			}
		}
	}

	return results
}

// getSearchMatches returns which fields matched the search query
func (r *Registry) getSearchMatches(shortcut Shortcut, queryLower string) []string {
	var matches []string

	if strings.Contains(strings.ToLower(shortcut.Keys), queryLower) {
		matches = append(matches, "keys")
	}
	if strings.Contains(strings.ToLower(shortcut.Description), queryLower) {
		matches = append(matches, "description")
	}
	if strings.Contains(strings.ToLower(shortcut.Category), queryLower) {
		matches = append(matches, "category")
	}

	return matches
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}
