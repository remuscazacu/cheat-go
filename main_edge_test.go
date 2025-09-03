package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/config"
	"cheat-go/pkg/ui"
)

func TestInitialModel_ConfigLoadErrors(t *testing.T) {
	// Test model initialization with config errors
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to directory with invalid config
	os.Chdir(tmpDir)

	// Create invalid config file
	invalidConfig := "invalid: yaml: content: ["
	os.WriteFile("config.yaml", []byte(invalidConfig), 0644)

	// Should still initialize successfully with warnings
	m := initialModelWithDefaults()

	if m.Registry == nil {
		t.Error("should initialize registry even with config errors")
	}

	if m.Config == nil {
		t.Error("should initialize config even with invalid file")
	}

	// Should fallback to default config
	defaultConfig := config.DefaultConfig()
	if len(m.Config.Apps) != len(defaultConfig.Apps) {
		t.Error("should fallback to default config on invalid YAML")
	}
}

func TestInitialModel_AppLoadErrors(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	// Create config with non-existent apps
	testConfig := &config.Config{
		Apps:    []string{"nonexistent1", "nonexistent2"},
		Theme:   "default",
		DataDir: tmpDir,
	}

	configData, _ := yaml.Marshal(testConfig)
	os.WriteFile("config.yaml", configData, 0644)

	// Should still initialize successfully
	m := initialModelWithDefaults()

	if m.Registry == nil {
		t.Error("should initialize registry even with app load errors")
	}

	// Should have some table data (at least header)
	if len(m.Rows) == 0 {
		t.Error("should have at least header row")
	}
}

func TestInitialModel_ErrorLogging(t *testing.T) {
	// This test verifies the error handling paths in initialModel
	// Since we can't easily capture stderr in tests, we verify the model
	// still initializes correctly when errors occur

	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	// Create config that will cause warnings
	testConfig := &config.Config{
		Apps:    []string{"vim", "nonexistent"},
		DataDir: "/invalid/path/that/does/not/exist",
	}

	configData, _ := yaml.Marshal(testConfig)
	os.WriteFile("config.yaml", configData, 0644)

	m := initialModelWithDefaults()

	// Model should still be functional
	if m.Registry == nil || m.Config == nil || m.Renderer == nil {
		t.Error("model should be fully initialized despite warnings")
	}

	// Should have vim app (which exists in hardcoded data)
	if _, exists := m.Registry.Get("vim"); !exists {
		t.Error("should have vim app from hardcoded data")
	}
}

func TestInitialModel_DefaultPaths(t *testing.T) {
	// Test initialization without any config files (pure defaults)
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to empty directory
	os.Chdir(tmpDir)

	m := initialModelWithDefaults()

	// Should use all defaults
	defaultConfig := config.DefaultConfig()
	if m.Config.Theme != defaultConfig.Theme {
		t.Error("should use default theme")
	}

	if len(m.Config.Apps) != len(defaultConfig.Apps) {
		t.Error("should use default apps")
	}

	// Should have all hardcoded apps
	expectedApps := []string{"vim", "zsh", "dwm", "st", "lf", "zathura"}
	for _, appName := range expectedApps {
		if _, exists := m.Registry.Get(appName); !exists {
			t.Errorf("should have hardcoded app: %s", appName)
		}
	}
}

func TestModel_Update_ExtensiveKeyHandling(t *testing.T) {
	m := initialModelWithDefaults()

	// Test all key types that should be handled
	keyTests := []struct {
		key         string
		expectCmd   bool
		description string
	}{
		{"q", true, "quit key"},
		{"ctrl+c", true, "ctrl+c quit"},
		{"up", false, "up arrow"},
		{"down", false, "down arrow"},
		{"left", false, "left arrow"},
		{"right", false, "right arrow"},
		{"k", false, "vim up"},
		{"j", false, "vim down"},
		{"h", false, "vim left"},
		{"l", false, "vim right"},
		{"x", false, "unknown key"},
		{"enter", false, "enter key"},
		{"space", false, "space key"},
	}

	for _, test := range keyTests {
		var msg tea.KeyMsg

		switch test.key {
		case "q", "k", "j", "h", "l", "x":
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(test.key[0])}}
		case "ctrl+c":
			msg = tea.KeyMsg{Type: tea.KeyCtrlC}
		case "up":
			msg = tea.KeyMsg{Type: tea.KeyUp}
		case "down":
			msg = tea.KeyMsg{Type: tea.KeyDown}
		case "left":
			msg = tea.KeyMsg{Type: tea.KeyLeft}
		case "right":
			msg = tea.KeyMsg{Type: tea.KeyRight}
		case "enter":
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		case "space":
			msg = tea.KeyMsg{Type: tea.KeySpace}
		}

		newModel, cmd := m.Update(msg)

		if test.expectCmd && cmd == nil {
			t.Errorf("%s: expected command but got nil", test.description)
		}

		if !test.expectCmd && cmd != nil {
			t.Errorf("%s: expected no command but got one", test.description)
		}

		if newModel == nil {
			t.Errorf("%s: should always return model", test.description)
		}

		// Update model for next test (except for quit commands)
		if cmd == nil {
			m = newModel.(ui.Model)
		}
	}
}
