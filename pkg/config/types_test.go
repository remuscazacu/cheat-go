package config

import (
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test default apps
	expectedApps := []string{"vim", "zsh", "dwm", "st", "lf", "zathura"}
	if !reflect.DeepEqual(config.Apps, expectedApps) {
		t.Errorf("default apps = %v, expected %v", config.Apps, expectedApps)
	}

	// Test default theme
	if config.Theme != "default" {
		t.Errorf("default theme = %s, expected 'default'", config.Theme)
	}

	// Test default layout
	expectedColumns := []string{"shortcut", "description"}
	if !reflect.DeepEqual(config.Layout.Columns, expectedColumns) {
		t.Errorf("default columns = %v, expected %v", config.Layout.Columns, expectedColumns)
	}

	if config.Layout.ShowCategories != false {
		t.Error("default ShowCategories should be false")
	}

	if config.Layout.TableStyle != "simple" {
		t.Errorf("default TableStyle = %s, expected 'simple'", config.Layout.TableStyle)
	}

	if config.Layout.MaxWidth != 120 {
		t.Errorf("default MaxWidth = %d, expected 120", config.Layout.MaxWidth)
	}

	// Test default keybinds
	expectedKeybinds := map[string]string{
		"quit":     "q",
		"up":       "k",
		"down":     "j",
		"left":     "h",
		"right":    "l",
		"search":   "/",
		"next_app": "tab",
		"prev_app": "shift+tab",
	}

	if !reflect.DeepEqual(config.Keybinds, expectedKeybinds) {
		t.Errorf("default keybinds = %v, expected %v", config.Keybinds, expectedKeybinds)
	}

	// Test default data directory
	if config.DataDir != "~/.config/cheat-go/apps" {
		t.Errorf("default DataDir = %s, expected '~/.config/cheat-go/apps'", config.DataDir)
	}
}

func TestLayoutConfig_Structure(t *testing.T) {
	layout := LayoutConfig{
		Columns:        []string{"key", "vim", "zsh"},
		ShowCategories: true,
		TableStyle:     "bordered",
		MaxWidth:       100,
	}

	if len(layout.Columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(layout.Columns))
	}

	if !layout.ShowCategories {
		t.Error("ShowCategories should be true")
	}

	if layout.TableStyle != "bordered" {
		t.Errorf("TableStyle = %s, expected 'bordered'", layout.TableStyle)
	}

	if layout.MaxWidth != 100 {
		t.Errorf("MaxWidth = %d, expected 100", layout.MaxWidth)
	}
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		Apps:     []string{"custom-app"},
		Theme:    "dark",
		Layout:   LayoutConfig{Columns: []string{"test"}},
		Keybinds: map[string]string{"test": "t"},
		DataDir:  "/custom/path",
	}

	if len(config.Apps) != 1 || config.Apps[0] != "custom-app" {
		t.Error("Apps should be set correctly")
	}

	if config.Theme != "dark" {
		t.Error("Theme should be set correctly")
	}

	if len(config.Layout.Columns) != 1 || config.Layout.Columns[0] != "test" {
		t.Error("Layout.Columns should be set correctly")
	}

	if config.Keybinds["test"] != "t" {
		t.Error("Keybinds should be set correctly")
	}

	if config.DataDir != "/custom/path" {
		t.Error("DataDir should be set correctly")
	}
}

func TestConfig_EmptyValues(t *testing.T) {
	config := &Config{}

	// Test zero values
	if config.Apps != nil && len(config.Apps) > 0 {
		t.Error("empty config Apps should be nil or empty")
	}

	if config.Theme != "" {
		t.Error("empty config Theme should be empty string")
	}

	if config.DataDir != "" {
		t.Error("empty config DataDir should be empty string")
	}

	if config.Keybinds != nil && len(config.Keybinds) > 0 {
		t.Error("empty config Keybinds should be nil or empty")
	}
}

func TestLayoutConfig_EmptyValues(t *testing.T) {
	layout := LayoutConfig{}

	if layout.Columns != nil && len(layout.Columns) > 0 {
		t.Error("empty layout Columns should be nil or empty")
	}

	if layout.ShowCategories != false {
		t.Error("empty layout ShowCategories should be false")
	}

	if layout.TableStyle != "" {
		t.Error("empty layout TableStyle should be empty string")
	}

	if layout.MaxWidth != 0 {
		t.Error("empty layout MaxWidth should be 0")
	}
}
