package config

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidTheme      = errors.New("invalid theme")
	ErrInvalidTableStyle = errors.New("invalid table style")
	ErrInvalidColumn     = errors.New("invalid column")
	ErrInvalidKeybind    = errors.New("invalid keybind")
	ErrInvalidMaxWidth   = errors.New("invalid max width")
)

// Config represents the main application configuration
type Config struct {
	Apps     []string          `yaml:"apps" json:"apps"`
	Theme    string            `yaml:"theme" json:"theme"`
	Layout   LayoutConfig      `yaml:"layout" json:"layout"`
	Keybinds map[string]string `yaml:"keybinds" json:"keybinds"`
	DataDir  string            `yaml:"data_dir" json:"data_dir"`
}

// LayoutConfig controls the display layout
type LayoutConfig struct {
	Columns        []string `yaml:"columns" json:"columns"`
	ShowCategories bool     `yaml:"show_categories" json:"show_categories"`
	TableStyle     string   `yaml:"table_style" json:"table_style"`
	MaxWidth       int      `yaml:"max_width" json:"max_width"`
}

// ValidationResult contains validation information
type ValidationResult struct {
	Valid  bool
	Errors []error
}

// ValidThemes contains all supported themes
var ValidThemes = []string{"default", "dark", "light", "minimal"}

// ValidTableStyles contains all supported table styles
var ValidTableStyles = []string{"simple", "rounded", "bold", "minimal"}

// ValidColumns contains all supported columns
var ValidColumns = []string{"shortcut", "description", "category", "tags", "platform"}

// RequiredKeybinds contains all required keybind actions
var RequiredKeybinds = []string{"quit", "up", "down", "left", "right"}

// Validate validates the configuration and returns validation results
func (c *Config) Validate() ValidationResult {
	var errors []error

	// Validate theme
	if !isValidTheme(c.Theme) {
		errors = append(errors, fmt.Errorf("%w: %s (valid: %v)", ErrInvalidTheme, c.Theme, ValidThemes))
	}

	// Validate layout
	if validationErrors := c.Layout.validate(); len(validationErrors) > 0 {
		errors = append(errors, validationErrors...)
	}

	// Validate keybinds
	if validationErrors := c.validateKeybinds(); len(validationErrors) > 0 {
		errors = append(errors, validationErrors...)
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

// validate validates the layout configuration
func (l *LayoutConfig) validate() []error {
	var errors []error

	// Validate columns
	for _, column := range l.Columns {
		if !isValidColumn(column) {
			errors = append(errors, fmt.Errorf("%w: %s (valid: %v)", ErrInvalidColumn, column, ValidColumns))
		}
	}

	// Validate table style
	if !isValidTableStyle(l.TableStyle) {
		errors = append(errors, fmt.Errorf("%w: %s (valid: %v)", ErrInvalidTableStyle, l.TableStyle, ValidTableStyles))
	}

	// Validate max width
	if l.MaxWidth < 40 || l.MaxWidth > 200 {
		errors = append(errors, fmt.Errorf("%w: %d (must be between 40 and 200)", ErrInvalidMaxWidth, l.MaxWidth))
	}

	return errors
}

// validateKeybinds validates the keybind configuration
func (c *Config) validateKeybinds() []error {
	var errors []error

	// Check for required keybinds
	for _, required := range RequiredKeybinds {
		if _, exists := c.Keybinds[required]; !exists {
			errors = append(errors, fmt.Errorf("%w: missing required keybind '%s'", ErrInvalidKeybind, required))
		}
	}

	// Check for duplicate keybind values
	usedKeys := make(map[string]string)
	for action, key := range c.Keybinds {
		if existingAction, exists := usedKeys[key]; exists {
			errors = append(errors, fmt.Errorf("%w: key '%s' used by both '%s' and '%s'", ErrInvalidKeybind, key, existingAction, action))
		}
		usedKeys[key] = action
	}

	return errors
}

// isValidTheme checks if the theme is valid
func isValidTheme(theme string) bool {
	for _, valid := range ValidThemes {
		if theme == valid {
			return true
		}
	}
	return false
}

// isValidTableStyle checks if the table style is valid
func isValidTableStyle(style string) bool {
	for _, valid := range ValidTableStyles {
		if style == valid {
			return true
		}
	}
	return false
}

// isValidColumn checks if the column is valid
func isValidColumn(column string) bool {
	for _, valid := range ValidColumns {
		if column == valid {
			return true
		}
	}
	return false
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Apps:  []string{"vim", "zsh", "dwm", "st", "lf", "zathura"},
		Theme: "default",
		Layout: LayoutConfig{
			Columns:        []string{"shortcut", "description"},
			ShowCategories: false,
			TableStyle:     "simple",
			MaxWidth:       120,
		},
		Keybinds: map[string]string{
			"quit":     "q",
			"up":       "k",
			"down":     "j",
			"left":     "h",
			"right":    "l",
			"search":   "/",
			"next_app": "tab",
			"prev_app": "shift+tab",
		},
		DataDir: "~/.config/cheat-go/apps",
	}
}
