package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/config"
)

func TestInitialModel(t *testing.T) {
	// Test basic model initialization
	m := initialModelWithDefaults()

	if m.registry == nil {
		t.Error("model should have registry initialized")
	}

	if m.config == nil {
		t.Error("model should have config initialized")
	}

	if m.renderer == nil {
		t.Error("model should have renderer initialized")
	}

	if m.rows == nil {
		t.Error("model should have rows initialized")
	}

	// Test initial cursor position
	if m.cursorX != 0 {
		t.Errorf("initial cursorX = %d, expected 0", m.cursorX)
	}

	if m.cursorY != 1 {
		t.Errorf("initial cursorY = %d, expected 1", m.cursorY)
	}

	// Test that table data is populated
	if len(m.rows) == 0 {
		t.Error("model should have table data")
	}

	// Test header row
	if len(m.rows) > 0 && len(m.rows[0]) == 0 {
		t.Error("header row should not be empty")
	}
}

func TestInitialModel_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	testConfig := &config.Config{
		Apps:  []string{"vim", "zsh"},
		Theme: "dark",
	}

	configData, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Change to temp directory to test config loading
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)
	os.Rename(configPath, "config.yaml")

	m := initialModelWithDefaults()

	if len(m.config.Apps) != 2 {
		t.Errorf("expected 2 apps from config, got %d", len(m.config.Apps))
	}

	if m.config.Theme != "dark" {
		t.Errorf("expected dark theme from config, got %s", m.config.Theme)
	}
}

func TestModel_Init(t *testing.T) {
	m := initialModelWithDefaults()

	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil command")
	}
}

func TestModel_Update_Quit(t *testing.T) {
	m := initialModelWithDefaults()

	// Test quit with 'q'
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := m.Update(quitMsg)

	// Check if quit command is returned
	if cmd == nil {
		t.Error("'q' key should return quit command")
	}

	// Model should be returned (bubbletea requirement)
	if newModel == nil {
		t.Error("Update should return model")
	}

	// Test quit with 'ctrl+c'
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = m.Update(ctrlCMsg)

	if cmd == nil {
		t.Error("'ctrl+c' should return quit command")
	}
}

func TestModel_Update_Navigation(t *testing.T) {
	m := initialModelWithDefaults()
	originalY := m.cursorY
	originalX := m.cursorX

	// Test down movement
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, cmd := m.Update(downMsg)

	if cmd != nil {
		t.Error("navigation should not return command")
	}

	m = newModel.(model)
	if m.cursorY <= originalY && len(m.rows) > originalY+1 {
		t.Error("down key should move cursor down")
	}

	// Test up movement
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(upMsg)
	m = newModel.(model)

	if m.cursorY != originalY && m.cursorY > 1 {
		t.Error("up key should move cursor up")
	}

	// Test right movement
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ = m.Update(rightMsg)
	m = newModel.(model)

	if m.cursorX <= originalX && len(m.rows[0]) > originalX+1 {
		t.Error("right key should move cursor right")
	}

	// Test left movement
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ = m.Update(leftMsg)
	m = newModel.(model)

	if m.cursorX != originalX && m.cursorX > 0 {
		t.Error("left key should move cursor left")
	}
}

func TestModel_Update_VimNavigation(t *testing.T) {
	m := initialModelWithDefaults()

	// Test vim-style navigation
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	hMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	lMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}

	originalY := m.cursorY

	// Test 'j' (down)
	newModel, _ := m.Update(jMsg)
	m = newModel.(model)
	if m.cursorY <= originalY && len(m.rows) > originalY+1 {
		t.Error("'j' should move cursor down")
	}

	// Test 'k' (up)
	newModel, _ = m.Update(kMsg)
	m = newModel.(model)
	if m.cursorY != originalY && m.cursorY > 1 {
		t.Error("'k' should move cursor up")
	}

	originalX := m.cursorX

	// Test 'l' (right)
	newModel, _ = m.Update(lMsg)
	m = newModel.(model)
	if m.cursorX <= originalX && len(m.rows[0]) > originalX+1 {
		t.Error("'l' should move cursor right")
	}

	// Test 'h' (left)
	newModel, _ = m.Update(hMsg)
	m = newModel.(model)
	if m.cursorX != originalX && m.cursorX > 0 {
		t.Error("'h' should move cursor left")
	}
}

func TestModel_Update_Boundaries(t *testing.T) {
	m := initialModelWithDefaults()

	// Move cursor to top-left
	m.cursorX = 0
	m.cursorY = 1 // Header is row 0, data starts at 1

	// Test left boundary
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(leftMsg)
	m = newModel.(model)
	if m.cursorX != 0 {
		t.Error("cursor should not move left from left boundary")
	}

	// Test up boundary
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(upMsg)
	m = newModel.(model)
	if m.cursorY != 1 {
		t.Error("cursor should not move up from top boundary")
	}

	// Move cursor to bottom-right
	m.cursorX = len(m.rows[0]) - 1
	m.cursorY = len(m.rows) - 1

	// Test right boundary
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ = m.Update(rightMsg)
	m = newModel.(model)
	if m.cursorX != len(m.rows[0])-1 {
		t.Error("cursor should not move right from right boundary")
	}

	// Test down boundary
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(downMsg)
	m = newModel.(model)
	if m.cursorY != len(m.rows)-1 {
		t.Error("cursor should not move down from bottom boundary")
	}
}

func TestModel_Update_UnknownKey(t *testing.T) {
	m := initialModelWithDefaults()
	originalX := m.cursorX
	originalY := m.cursorY

	// Test unknown key
	unknownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	newModel, cmd := m.Update(unknownMsg)

	if cmd != nil {
		t.Error("unknown key should not return command")
	}

	m = newModel.(model)
	if m.cursorX != originalX || m.cursorY != originalY {
		t.Error("unknown key should not change cursor position")
	}
}

func TestModel_View(t *testing.T) {
	m := initialModelWithDefaults()

	view := m.View()

	if view == "" {
		t.Error("View() should return non-empty string")
	}

	// Should contain table structure
	if !strings.Contains(view, "│") && !strings.Contains(view, "─") {
		t.Error("View should contain table formatting")
	}

	// Should contain instructions
	if !strings.Contains(view, "arrow keys") {
		t.Error("View should contain usage instructions")
	}
}

func TestModel_Integration(t *testing.T) {
	m := initialModelWithDefaults()

	// Test a sequence of operations
	keys := []tea.KeyMsg{
		{Type: tea.KeyDown},
		{Type: tea.KeyDown},
		{Type: tea.KeyRight},
		{Type: tea.KeyUp},
		{Type: tea.KeyLeft},
	}

	for _, key := range keys {
		newModel, cmd := m.Update(key)
		if cmd != nil {
			t.Error("navigation should not produce commands")
		}
		m = newModel.(model)
	}

	// Model should still be functional
	view := m.View()
	if view == "" {
		t.Error("model should still render after navigation sequence")
	}
}

func TestModel_EmptyTable(t *testing.T) {
	m := initialModelWithDefaults()
	m.rows = [][]string{} // Force empty table

	view := m.View()
	// Should handle empty table gracefully
	if view == "" {
		t.Error("should handle empty table")
	}
}

func TestModel_SingleRowTable(t *testing.T) {
	m := initialModelWithDefaults()
	m.rows = [][]string{{"Header"}} // Only header

	// Navigation should handle single row gracefully
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(downMsg)
	m = newModel.(model)

	if m.cursorY != 1 { // Should stay at row 1 (since no data rows)
		t.Error("cursor should handle single row table")
	}
}

func TestModel_Structure(t *testing.T) {
	m := initialModelWithDefaults()

	// Test model structure
	if m.registry == nil {
		t.Error("registry should not be nil")
	}

	if m.config == nil {
		t.Error("config should not be nil")
	}

	if m.renderer == nil {
		t.Error("renderer should not be nil")
	}

	if m.rows == nil {
		t.Error("rows should not be nil")
	}

	// Test that rows have expected structure (header + data)
	if len(m.rows) < 1 {
		t.Error("should have at least header row")
	}

	if len(m.rows[0]) == 0 {
		t.Error("header row should not be empty")
	}
}

func TestModel_Update_KeyTypes(t *testing.T) {
	m := initialModelWithDefaults()

	// Test different key message types
	testCases := []struct {
		msg  tea.KeyMsg
		desc string
	}{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, "rune key"},
		{tea.KeyMsg{Type: tea.KeyCtrlC}, "ctrl key"},
		{tea.KeyMsg{Type: tea.KeyUp}, "arrow key"},
		{tea.KeyMsg{Type: tea.KeyEnter}, "special key"},
	}

	for _, tc := range testCases {
		newModel, cmd := m.Update(tc.msg)

		// Should always return a model
		if newModel == nil {
			t.Errorf("%s: Update should return model", tc.desc)
		}

		// Some keys should return commands, others shouldn't
		switch tc.msg.Type {
		case tea.KeyCtrlC:
			if cmd == nil {
				t.Errorf("%s: should return quit command", tc.desc)
			}
		case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
			if cmd != nil {
				t.Errorf("%s: navigation should not return command", tc.desc)
			}
		}
	}
}
