package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrInvalidConfig  = errors.New("invalid configuration format")
)

// Loader handles configuration loading and validation
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load reads and parses the configuration file
func (l *Loader) Load() (*Config, error) {
	// Try to load from file first
	if l.configPath != "" {
		if config, err := l.loadFromFile(l.configPath); err == nil {
			return config, nil
		}
	}

	// Try default config locations
	defaultPaths := []string{
		"~/.config/cheat-go/config.yaml",
		"~/.cheat-go.yaml",
		"./config.yaml",
	}

	for _, path := range defaultPaths {
		expandedPath := expandPath(path)
		if config, err := l.loadFromFile(expandedPath); err == nil {
			return config, nil
		}
	}

	// Fall back to default configuration
	return DefaultConfig(), nil
}

// loadFromFile loads configuration from a specific file
func (l *Loader) loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, ErrInvalidConfig
	}

	return l.validateAndSetDefaults(&config)
}

// validateAndSetDefaults validates config and sets defaults for missing values
func (l *Loader) validateAndSetDefaults(config *Config) (*Config, error) {
	defaults := DefaultConfig()

	// Set defaults for missing values
	if len(config.Apps) == 0 {
		config.Apps = defaults.Apps
	}

	if config.Theme == "" {
		config.Theme = defaults.Theme
	}

	if len(config.Layout.Columns) == 0 {
		config.Layout = defaults.Layout
	} else {
		// Merge layout defaults for missing fields
		if config.Layout.TableStyle == "" {
			config.Layout.TableStyle = defaults.Layout.TableStyle
		}
		if config.Layout.MaxWidth == 0 {
			config.Layout.MaxWidth = defaults.Layout.MaxWidth
		}
	}

	if len(config.Keybinds) == 0 {
		config.Keybinds = defaults.Keybinds
	} else {
		// Merge keybind defaults for missing required keys
		for key, value := range defaults.Keybinds {
			if _, exists := config.Keybinds[key]; !exists {
				config.Keybinds[key] = value
			}
		}
	}

	if config.DataDir == "" {
		config.DataDir = defaults.DataDir
	}

	// Validate the configuration
	validation := config.Validate()
	if !validation.Valid {
		return nil, fmt.Errorf("configuration validation failed: %v", validation.Errors)
	}

	return config, nil
}

// Save writes the configuration to file
func (l *Loader) Save(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
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
