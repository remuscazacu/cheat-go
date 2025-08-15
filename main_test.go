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
	if !strings.Contains(view, "hjkl") {
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

// Test search functionality
func TestSearchFunctionality(t *testing.T) {
	m := initialModel()

	// Test entering search mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering search mode should not return command")
	}

	updatedModel := newModel.(model)
	if !updatedModel.searchMode {
		t.Error("should enter search mode when '/' is pressed")
	}
	if updatedModel.searchQuery != "" {
		t.Error("search query should be empty when entering search mode")
	}
}

func TestSearchInput(t *testing.T) {
	m := initialModel()
	m.searchMode = true

	// Test adding characters to search query
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("adding to search query should not return command")
	}

	updatedModel := newModel.(model)
	if updatedModel.searchQuery != "v" {
		t.Errorf("search query should be 'v', got '%s'", updatedModel.searchQuery)
	}

	// Test backspace
	msg = tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("backspace should not return command")
	}

	updatedModel = newModel.(model)
	if updatedModel.searchQuery != "" {
		t.Errorf("search query should be empty after backspace, got '%s'", updatedModel.searchQuery)
	}
}

func TestSearchEscape(t *testing.T) {
	m := initialModel()
	m.searchMode = true
	m.searchQuery = "test"

	// Test escape exits search mode
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}

	updatedModel := newModel.(model)
	if updatedModel.searchMode {
		t.Error("escape should exit search mode")
	}
	if updatedModel.searchQuery != "" {
		t.Error("escape should clear search query")
	}
}

func TestSearchEnter(t *testing.T) {
	m := initialModel()
	m.searchMode = true
	m.searchQuery = "vim"

	// Test enter confirms search
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("enter should not return command")
	}

	updatedModel := newModel.(model)
	if updatedModel.searchMode {
		t.Error("enter should exit search mode")
	}
	// Should have filtered results
	if len(updatedModel.rows) == len(m.allRows) {
		t.Error("search should filter results")
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	testCases := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "foo", false},
		{"vim", "V", true},
		{"VIM", "vim", true},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tc := range testCases {
		result := containsIgnoreCase(tc.s, tc.substr)
		if result != tc.expected {
			t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v",
				tc.s, tc.substr, result, tc.expected)
		}
	}
}

func TestFilterRowsBySearch(t *testing.T) {
	m := initialModel()

	// Set up test data
	m.allRows = [][]string{
		{"Shortcut", "vim", "zsh"},
		{"k", "↑ move", "up history"},
		{"j", "↓ move", "down history"},
		{"q", "quit", "exit"},
	}

	// Test empty query returns all rows
	filtered := m.filterRowsBySearch("")
	if len(filtered) != len(m.allRows) {
		t.Error("empty query should return all rows")
	}

	// Test filtering by "move"
	filtered = m.filterRowsBySearch("move")
	expectedCount := 3 // header + 2 move commands
	if len(filtered) != expectedCount {
		t.Errorf("search for 'move' should return %d rows, got %d", expectedCount, len(filtered))
	}

	// Test filtering by "quit"
	filtered = m.filterRowsBySearch("quit")
	expectedCount = 2 // header + 1 quit command
	if len(filtered) != expectedCount {
		t.Errorf("search for 'quit' should return %d rows, got %d", expectedCount, len(filtered))
	}

	// Test no matches
	filtered = m.filterRowsBySearch("nonexistent")
	expectedCount = 1 // just header
	if len(filtered) != expectedCount {
		t.Errorf("search for nonexistent term should return %d row (header only), got %d", expectedCount, len(filtered))
	}
}

func TestSearchView(t *testing.T) {
	m := initialModel()

	// Test normal view
	view := m.View()
	if !strings.Contains(view, "search") {
		t.Error("normal view should contain search instructions")
	}

	// Test search mode view
	m.searchMode = true
	m.searchQuery = "test"
	view = m.View()
	if !strings.Contains(view, "Search: test_") {
		t.Error("search mode view should show search query with cursor")
	}
	if !strings.Contains(view, "Enter to confirm") {
		t.Error("search mode view should show search instructions")
	}
}

func TestSearchUIEnhancements(t *testing.T) {
	m := initialModel()
	
	// Test filter status display
	m.allRows = [][]string{
		{"Shortcut", "vim", "zsh"},
		{"k", "↑ move", "up history"},
		{"j", "↓ move", "down history"},
		{"q", "quit", "exit"},
	}
	
	// Test when results are filtered
	m.rows = [][]string{
		{"Shortcut", "vim", "zsh"},
		{"k", "↑ move", "up history"},
		{"j", "↓ move", "down history"},
	}
	
	view := m.View()
	if !strings.Contains(view, "2/3 results") {
		t.Error("should show filter status when results are filtered")
	}
	
	// Test normal state (no filter)
	m.rows = m.allRows
	view = m.View()
	if strings.Contains(view, "results") {
		t.Error("should not show filter status when not filtered")
	}
}

func TestEscapeClearSearch(t *testing.T) {
	m := initialModel()
	m.allRows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"q", "quit"},
	}
	m.rows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
	}
	
	// Test escape clears filter
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}
	
	updatedModel := newModel.(model)
	if len(updatedModel.rows) != len(m.allRows) {
		t.Error("escape should restore all rows")
	}
	if updatedModel.cursorY != 1 {
		t.Error("escape should reset cursor position")
	}
}

func TestSearchModeStyledView(t *testing.T) {
	m := initialModel()
	m.searchMode = true
	m.searchQuery = "test"
	
	view := m.View()
	if !strings.Contains(view, "Search:") {
		t.Error("search mode view should contain search prompt")
	}
	if !strings.Contains(view, "test_") {
		t.Error("search mode view should show query with cursor")
	}
	if !strings.Contains(view, "Type to search") {
		t.Error("search mode view should show search instructions")
	}
}


func TestAppFiltering(t *testing.T) {
	m := initialModel()
	
	// Test entering filter mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering filter mode should not return command")
	}
	
	updatedModel := newModel.(model)
	if !updatedModel.filterMode {
		t.Error("should enter filter mode when \"f\" is pressed")
	}
}

func TestFilterInput(t *testing.T) {
	m := initialModel()
	m.filterMode = true
	m.allApps = []string{"vim", "zsh", "dwm"}
	
	// Test selecting app by number
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("selecting app should not return command")
	}
	
	updatedModel := newModel.(model)
	if len(updatedModel.filteredApps) != 1 {
		t.Error("should add app to filtered list")
	}
	if updatedModel.filteredApps[0] != "vim" {
		t.Error("should add correct app to filtered list")
	}
	
	// Test toggling same app (should remove)
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("toggling app should not return command")
	}
	
	updatedModel = newModel.(model)
	if len(updatedModel.filteredApps) != 0 {
		t.Error("should remove app from filtered list when toggled")
	}
}

func TestFilterSelectAll(t *testing.T) {
	m := initialModel()
	m.filterMode = true
	m.allApps = []string{"vim", "zsh", "dwm"}
	
	// Test select all
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("select all should not return command")
	}
	
	updatedModel := newModel.(model)
	if len(updatedModel.filteredApps) != len(m.allApps) {
		t.Error("select all should add all apps to filtered list")
	}
}

func TestFilterClear(t *testing.T) {
	m := initialModel()
	m.filterMode = true
	m.allApps = []string{"vim", "zsh", "dwm"}
	m.filteredApps = []string{"vim", "zsh"}
	
	// Test clear
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("clear should not return command")
	}
	
	updatedModel := newModel.(model)
	if len(updatedModel.filteredApps) != 0 {
		t.Error("clear should remove all apps from filtered list")
	}
}

func TestIsAppSelected(t *testing.T) {
	m := initialModel()
	m.filteredApps = []string{"vim", "zsh"}
	
	if !m.isAppSelected("vim") {
		t.Error("should return true for selected app")
	}
	if !m.isAppSelected("zsh") {
		t.Error("should return true for selected app")
	}
	if m.isAppSelected("dwm") {
		t.Error("should return false for unselected app")
	}
}

func TestFilterModeView(t *testing.T) {
	m := initialModel()
	m.filterMode = true
	m.allApps = []string{"vim", "zsh"}
	m.filteredApps = []string{"vim"}
	
	view := m.View()
	if !strings.Contains(view, "Filter Apps:") {
		t.Error("filter mode view should contain filter prompt")
	}
	if !strings.Contains(view, "[1]") {
		t.Error("filter mode view should show app numbers")
	}
	if !strings.Contains(view, "✓vim") {
		t.Error("filter mode view should show selected apps with checkmark")
	}
	if !strings.Contains(view, " zsh") {
		t.Error("filter mode view should show unselected apps without checkmark")
	}
}


func TestKeyboardShortcuts(t *testing.T) {
	m := initialModel()
	
	// Test help mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering help mode should not return command")
	}
	
	updatedModel := newModel.(model)
	if !updatedModel.helpMode {
		t.Error("should enter help mode when ? is pressed")
	}
	
	// Test help view
	view := updatedModel.View()
	if !strings.Contains(view, "KEYBOARD SHORTCUTS") {
		t.Error("help view should show shortcuts title")
	}
	if !strings.Contains(view, "NAVIGATION:") {
		t.Error("help view should show navigation section")
	}
}

func TestRefreshShortcut(t *testing.T) {
	m := initialModel()
	
	// Test refresh shortcut
	msg := tea.KeyMsg{Type: tea.KeyCtrlR}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("refresh should not return command")
	}
	
	updatedModel := newModel.(model)
	if updatedModel.cursorY != 1 {
		t.Error("refresh should reset cursor position")
	}
}

func TestHomeEndShortcuts(t *testing.T) {
	m := initialModel()
	m.rows = [][]string{
		{"Header"},
		{"Row1"},
		{"Row2"},
		{"Row3"},
	}
	m.cursorY = 2
	
	// Test home shortcut
	msg := tea.KeyMsg{Type: tea.KeyHome}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("home should not return command")
	}
	
	updatedModel := newModel.(model)
	if updatedModel.cursorY != 1 {
		t.Error("home should go to first data row")
	}
	
	// Test end shortcut
	msg = tea.KeyMsg{Type: tea.KeyEnd}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("end should not return command")
	}
	
	updatedModel = newModel.(model)
	if updatedModel.cursorY != 3 {
		t.Error("end should go to last row")
	}
}

func TestAlternativeShortcuts(t *testing.T) {
	m := initialModel()
	
	// Test Ctrl+F for filter mode
	msg := tea.KeyMsg{Type: tea.KeyCtrlF}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("ctrl+f should not return command")
	}
	
	updatedModel := newModel.(model)
	if !updatedModel.filterMode {
		t.Error("ctrl+f should enter filter mode")
	}
}

func TestSearchClearShortcut(t *testing.T) {
	m := initialModel()
	m.searchMode = true
	m.searchQuery = "test query"
	
	// Test Ctrl+U to clear search
	msg := tea.KeyMsg{Type: tea.KeyCtrlU}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("ctrl+u should not return command")
	}
	
	updatedModel := newModel.(model)
	if updatedModel.searchQuery != "" {
		t.Error("ctrl+u should clear search query")
	}
}


func TestSearchHighlighting(t *testing.T) {
	m := initialModel()
	m.lastSearch = "move"
	m.rows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"j", "↓ move"},
	}
	m.allRows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"j", "↓ move"},
		{"q", "quit"},
	}
	
	view := m.View()
	// Should use highlighting when there is a search term and filtered results
	if m.lastSearch != "" && len(m.rows) != len(m.allRows) {
		// The highlighting is applied, but we cannot easily test the exact output
		// since it contains ANSI escape codes. We just verify the view renders.
		if view == "" {
			t.Error("view with highlighting should not be empty")
		}
	}
}

func TestLastSearchTracking(t *testing.T) {
	m := initialModel()
	m.searchMode = true
	m.searchQuery = "test"
	
	// Test that lastSearch is set when confirming search
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("enter should not return command")
	}
	
	updatedModel := newModel.(model)
	if updatedModel.lastSearch != "test" {
		t.Error("should set lastSearch when confirming search")
	}
	
	// Test that lastSearch is cleared when escaping from search
	updatedModel.searchMode = true
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}
	
	finalModel := newModel.(model)
	if finalModel.lastSearch != "" {
		t.Error("should clear lastSearch when escaping from search")
	}
}

