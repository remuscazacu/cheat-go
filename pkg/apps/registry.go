package apps

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrAppNotFound    = errors.New("application not found")
	ErrInvalidAppFile = errors.New("invalid app file format")
	ErrAppValidation  = errors.New("app validation failed")
	ErrDirectoryRead  = errors.New("failed to read app directory")
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
			// Log error but continue with next app (backward compatibility)
			continue
		}
	}
	return nil
}

// LoadAllAppsFromDirectory scans the data directory and loads all available apps
func (r *Registry) LoadAllAppsFromDirectory() error {
	if r.dataDir == "" {
		return nil
	}

	expandedDir := expandPath(r.dataDir)
	entries, err := os.ReadDir(expandedDir)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDirectoryRead, expandedDir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		// Extract app name from filename
		appName := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
		if err := r.LoadApp(appName); err != nil {
			// Log but don't fail for individual app loading errors
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
		return nil, fmt.Errorf("%w: %v", ErrInvalidAppFile, err)
	}

	// Validate the loaded app
	if err := r.validateApp(&app); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppValidation, err)
	}

	return &app, nil
}

// validateApp validates an app definition
func (r *Registry) validateApp(app *App) error {
	var errors []error

	// Required fields
	if app.Name == "" {
		errors = append(errors, fmt.Errorf("app name is required"))
	}

	if app.Description == "" {
		errors = append(errors, fmt.Errorf("app description is required"))
	}

	// Validate shortcuts
	for i, shortcut := range app.Shortcuts {
		if shortcut.Keys == "" {
			errors = append(errors, fmt.Errorf("shortcut %d: keys field is required", i))
		}
		if shortcut.Description == "" {
			errors = append(errors, fmt.Errorf("shortcut %d: description field is required", i))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

// SaveApp saves an app definition to a YAML file
func (r *Registry) SaveApp(app *App) error {
	if r.dataDir == "" {
		return fmt.Errorf("data directory not configured")
	}

	// Validate before saving
	if err := r.validateApp(app); err != nil {
		return err
	}

	expandedDir := expandPath(r.dataDir)

	// Ensure directory exists
	if err := os.MkdirAll(expandedDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	appPath := filepath.Join(expandedDir, app.Name+".yaml")
	data, err := yaml.Marshal(app)
	if err != nil {
		return fmt.Errorf("failed to marshal app data: %w", err)
	}

	if err := os.WriteFile(appPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write app file: %w", err)
	}

	// Register the app in memory
	r.Register(app)

	return nil
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
