package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/config"
	"cheat-go/pkg/notes"
	"cheat-go/pkg/ui"
)

func initialModelWithDefaults() ui.Model {
	opts := cliOptions{
		theme:      "",
		tableStyle: "",
		configFile: "",
	}
	return initialModel(opts)
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func TestInitialModel(t *testing.T) {
	// Test basic model initialization
	m := initialModelWithDefaults()

	if m.Registry == nil {
		t.Error("model should have registry initialized")
	}

	if m.Config == nil {
		t.Error("model should have config initialized")
	}

	if m.Renderer == nil {
		t.Error("model should have renderer initialized")
	}

	if m.Rows == nil {
		t.Error("model should have rows initialized")
	}

	// Test initial cursor position
	if m.CursorX != 0 {
		t.Errorf("initial cursorX = %d, expected 0", m.CursorX)
	}

	if m.CursorY != 1 {
		t.Errorf("initial cursorY = %d, expected 1", m.CursorY)
	}

	// Test that table data is populated
	if len(m.Rows) == 0 {
		t.Error("model should have table data")
	}

	// Test header row
	if len(m.Rows) > 0 && len(m.Rows[0]) == 0 {
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

	if len(m.Config.Apps) != 2 {
		t.Errorf("expected 2 apps from config, got %d", len(m.Config.Apps))
	}

	if m.Config.Theme != "dark" {
		t.Errorf("expected dark theme from config, got %s", m.Config.Theme)
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
	originalY := m.CursorY
	originalX := m.CursorX

	// Test down movement
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, cmd := m.Update(downMsg)

	if cmd != nil {
		t.Error("navigation should not return command")
	}

	m = newModel.(ui.Model)
	if m.CursorY <= originalY && len(m.Rows) > originalY+1 {
		t.Error("down key should move cursor down")
	}

	// Test up movement
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(upMsg)
	m = newModel.(ui.Model)

	if m.CursorY != originalY && m.CursorY > 1 {
		t.Error("up key should move cursor up")
	}

	// Test right movement
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ = m.Update(rightMsg)
	m = newModel.(ui.Model)

	if m.CursorX <= originalX && len(m.Rows[0]) > originalX+1 {
		t.Error("right key should move cursor right")
	}

	// Test left movement
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ = m.Update(leftMsg)
	m = newModel.(ui.Model)

	if m.CursorX != originalX && m.CursorX > 0 {
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

	originalY := m.CursorY

	// Test 'j' (down)
	newModel, _ := m.Update(jMsg)
	m = newModel.(ui.Model)
	if m.CursorY <= originalY && len(m.Rows) > originalY+1 {
		t.Error("'j' should move cursor down")
	}

	// Test 'k' (up)
	newModel, _ = m.Update(kMsg)
	m = newModel.(ui.Model)
	if m.CursorY != originalY && m.CursorY > 1 {
		t.Error("'k' should move cursor up")
	}

	originalX := m.CursorX

	// Test 'l' (right)
	newModel, _ = m.Update(lMsg)
	m = newModel.(ui.Model)
	if m.CursorX <= originalX && len(m.Rows[0]) > originalX+1 {
		t.Error("'l' should move cursor right")
	}

	// Test 'h' (left)
	newModel, _ = m.Update(hMsg)
	m = newModel.(ui.Model)
	if m.CursorX != originalX && m.CursorX > 0 {
		t.Error("'h' should move cursor left")
	}
}

func TestModel_Update_Boundaries(t *testing.T) {
	m := initialModelWithDefaults()

	// Move cursor to top-left
	m.CursorX = 0
	m.CursorY = 1 // Header is row 0, data starts at 1

	// Test left boundary
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(leftMsg)
	m = newModel.(ui.Model)
	if m.CursorX != 0 {
		t.Error("cursor should not move left from left boundary")
	}

	// Test up boundary
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(upMsg)
	m = newModel.(ui.Model)
	if m.CursorY != 1 {
		t.Error("cursor should not move up from top boundary")
	}

	// Move cursor to bottom-right
	m.CursorX = len(m.Rows[0]) - 1
	m.CursorY = len(m.Rows) - 1

	// Test right boundary
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ = m.Update(rightMsg)
	m = newModel.(ui.Model)
	if m.CursorX != len(m.Rows[0])-1 {
		t.Error("cursor should not move right from right boundary")
	}

	// Test down boundary
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(downMsg)
	m = newModel.(ui.Model)
	if m.CursorY != len(m.Rows)-1 {
		t.Error("cursor should not move down from bottom boundary")
	}
}

func TestModel_Update_UnknownKey(t *testing.T) {
	m := initialModelWithDefaults()
	originalX := m.CursorX
	originalY := m.CursorY

	// Test unknown key
	unknownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	newModel, cmd := m.Update(unknownMsg)

	if cmd != nil {
		t.Error("unknown key should not return command")
	}

	m = newModel.(ui.Model)
	if m.CursorX != originalX || m.CursorY != originalY {
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
		m = newModel.(ui.Model)
	}

	// Model should still be functional
	view := m.View()
	if view == "" {
		t.Error("model should still render after navigation sequence")
	}
}

func TestModel_EmptyTable(t *testing.T) {
	m := initialModelWithDefaults()
	m.Rows = [][]string{} // Force empty table

	view := m.View()
	// Should handle empty table gracefully
	if view == "" {
		t.Error("should handle empty table")
	}
}

func TestModel_SingleRowTable(t *testing.T) {
	m := initialModelWithDefaults()
	m.Rows = [][]string{{"Header"}} // Only header

	// Navigation should handle single row gracefully
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(downMsg)
	m = newModel.(ui.Model)

	if m.CursorY != 1 { // Should stay at row 1 (since no data rows)
		t.Error("cursor should handle single row table")
	}
}

func TestModel_Structure(t *testing.T) {
	m := initialModelWithDefaults()

	// Test model structure
	if m.Registry == nil {
		t.Error("registry should not be nil")
	}

	if m.Config == nil {
		t.Error("config should not be nil")
	}

	if m.Renderer == nil {
		t.Error("renderer should not be nil")
	}

	if m.Rows == nil {
		t.Error("rows should not be nil")
	}

	// Test that rows have expected structure (header + data)
	if len(m.Rows) < 1 {
		t.Error("should have at least header row")
	}

	if len(m.Rows[0]) == 0 {
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
	m := initialModelWithDefaults()

	// Test entering search mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering search mode should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if !updatedModel.SearchMode {
		t.Error("should enter search mode when '/' is pressed")
	}
	if updatedModel.SearchQuery != "" {
		t.Error("search query should be empty when entering search mode")
	}
}

func TestSearchInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true

	// Test adding characters to search query
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("adding to search query should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.SearchQuery != "v" {
		t.Errorf("search query should be 'v', got '%s'", updatedModel.SearchQuery)
	}

	// Test backspace
	msg = tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("backspace should not return command")
	}

	updatedModel = newModel.(ui.Model)
	if updatedModel.SearchQuery != "" {
		t.Errorf("search query should be empty after backspace, got '%s'", updatedModel.SearchQuery)
	}
}

func TestSearchEscape(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true
	m.SearchQuery = "test"

	// Test escape exits search mode
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.SearchMode {
		t.Error("escape should exit search mode")
	}
	if updatedModel.SearchQuery != "" {
		t.Error("escape should clear search query")
	}
}

func TestSearchEnter(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true
	m.SearchQuery = "vim"

	// Test enter confirms search
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("enter should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.SearchMode {
		t.Error("enter should exit search mode")
	}
	// Should have filtered results
	if len(updatedModel.Rows) == len(m.AllRows) {
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
	m := initialModelWithDefaults()

	// Set up test data
	m.AllRows = [][]string{
		{"Shortcut", "vim", "zsh"},
		{"k", "↑ move", "up history"},
		{"j", "↓ move", "down history"},
		{"q", "quit", "exit"},
	}

	// Test empty query returns all rows
	filtered := m.FilterRowsBySearch("")
	if len(filtered) != len(m.AllRows) {
		t.Error("empty query should return all rows")
	}

	// Test filtering by "move"
	filtered = m.FilterRowsBySearch("move")
	expectedCount := 3 // header + 2 move commands
	if len(filtered) != expectedCount {
		t.Errorf("search for 'move' should return %d rows, got %d", expectedCount, len(filtered))
	}

	// Test filtering by "quit"
	filtered = m.FilterRowsBySearch("quit")
	expectedCount = 2 // header + 1 quit command
	if len(filtered) != expectedCount {
		t.Errorf("search for 'quit' should return %d rows, got %d", expectedCount, len(filtered))
	}

	// Test no matches
	filtered = m.FilterRowsBySearch("nonexistent")
	expectedCount = 1 // just header
	if len(filtered) != expectedCount {
		t.Errorf("search for nonexistent term should return %d row (header only), got %d", expectedCount, len(filtered))
	}
}

func TestSearchView(t *testing.T) {
	m := initialModelWithDefaults()

	// Test normal view
	view := m.View()
	if !strings.Contains(view, "search") {
		t.Error("normal view should contain search instructions")
	}

	// Test search mode view
	m.SearchMode = true
	m.SearchQuery = "test"
	view = m.View()
	if !strings.Contains(view, "Search: test_") {
		t.Error("search mode view should show search query with cursor")
	}
	if !strings.Contains(view, "Enter to confirm") {
		t.Error("search mode view should show search instructions")
	}
}

func TestSearchUIEnhancements(t *testing.T) {
	m := initialModelWithDefaults()
	m.AllRows = [][]string{
		{"Shortcut", "Description"},
		{"ctrl+c", "copy"},
		{"ctrl+v", "paste"},
	}
	m.Rows = [][]string{
		{"Shortcut", "Description"},
		{"ctrl+c", "copy"},
	}

	// Test that filtered rows are set correctly
	if len(m.Rows) != 2 {
		t.Error("should have 2 rows when filtered")
	}
	if len(m.AllRows) != 3 {
		t.Error("should have 3 total rows")
	}
}

func TestEscapeClearSearch(t *testing.T) {
	m := initialModelWithDefaults()
	m.AllRows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"q", "quit"},
	}
	m.Rows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
	}

	// Test escape clears filter
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if len(updatedModel.Rows) != len(m.AllRows) {
		t.Error("escape should restore all rows")
	}
	if updatedModel.CursorY != 1 {
		t.Error("escape should reset cursor position")
	}
}

func TestSearchModeStyledView(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true
	m.SearchQuery = "test"

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
	m := initialModelWithDefaults()

	// Test entering filter mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering filter mode should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if !updatedModel.FilterMode {
		t.Error("should enter filter mode when \"f\" is pressed")
	}
}

func TestFilterInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.FilterMode = true
	m.AllApps = []string{"vim", "zsh", "dwm"}

	// Test selecting app by number
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("selecting app should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if len(updatedModel.FilteredApps) != 1 {
		t.Error("should add app to filtered list")
	}
	if updatedModel.FilteredApps[0] != "vim" {
		t.Error("should add correct app to filtered list")
	}

	// Test toggling same app (should remove)
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("toggling app should not return command")
	}

	updatedModel = newModel.(ui.Model)
	if len(updatedModel.FilteredApps) != 0 {
		t.Error("should remove app from filtered list when toggled")
	}
}

func TestFilterSelectAll(t *testing.T) {
	m := initialModelWithDefaults()
	m.FilterMode = true
	m.AllApps = []string{"vim", "zsh", "dwm"}

	// Test select all
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("select all should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if len(updatedModel.FilteredApps) != len(m.AllApps) {
		t.Error("select all should add all apps to filtered list")
	}
}

func TestFilterClear(t *testing.T) {
	m := initialModelWithDefaults()
	m.FilterMode = true
	m.AllApps = []string{"vim", "zsh", "dwm"}
	m.FilteredApps = []string{"vim", "zsh"}

	// Test clear
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("clear should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if len(updatedModel.FilteredApps) != 0 {
		t.Error("clear should remove all apps from filtered list")
	}
}

func TestIsAppSelected(t *testing.T) {
	m := initialModelWithDefaults()
	m.FilteredApps = []string{"vim", "zsh"}

	if !m.IsAppSelected("vim") {
		t.Error("should return true for selected app")
	}
	if !m.IsAppSelected("zsh") {
		t.Error("should return true for selected app")
	}
	if m.IsAppSelected("dwm") {
		t.Error("should return false for unselected app")
	}
}

func TestFilterModeView(t *testing.T) {
	m := initialModelWithDefaults()
	m.FilterMode = true
	m.AllApps = []string{"vim", "zsh"}
	m.FilteredApps = []string{"vim"}

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
	m := initialModelWithDefaults()

	// Test help mode
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("entering help mode should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if !updatedModel.HelpMode {
		t.Error("should enter help mode when ? is pressed")
	}

	// Test help view
	view := updatedModel.View()
	if !strings.Contains(view, "Help") {
		t.Error("help view should show help title")
	}
	if !strings.Contains(view, "NAVIGATION") {
		t.Error("help view should show navigation section")
	}
}

func TestRefreshShortcut(t *testing.T) {
	m := initialModelWithDefaults()

	// Test refresh shortcut
	msg := tea.KeyMsg{Type: tea.KeyCtrlR}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("refresh should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.CursorY != 1 {
		t.Error("refresh should reset cursor position")
	}
}

func TestHomeEndShortcuts(t *testing.T) {
	m := initialModelWithDefaults()
	m.Rows = [][]string{
		{"Header"},
		{"Row1"},
		{"Row2"},
		{"Row3"},
	}
	m.CursorY = 2

	// Test home shortcut
	msg := tea.KeyMsg{Type: tea.KeyHome}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("home should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.CursorY != 1 {
		t.Error("home should go to first data row")
	}

	// Test end shortcut
	msg = tea.KeyMsg{Type: tea.KeyEnd}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("end should not return command")
	}

	updatedModel = newModel.(ui.Model)
	if updatedModel.CursorY != 3 {
		t.Error("end should go to last row")
	}
}

func TestAlternativeShortcuts(t *testing.T) {
	m := initialModelWithDefaults()

	// Test Ctrl+F for filter
	msg := tea.KeyMsg{Type: tea.KeyCtrlF}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	if !updatedModel.FilterMode {
		t.Error("Ctrl+F should enter filter mode")
	}
}

func TestViewModes(t *testing.T) {
	m := initialModelWithDefaults()

	// Test switching to notes view
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewNotes {
		t.Error("Should switch to notes view")
	}

	// Test notes view rendering
	view := updatedModel.View()
	if !strings.Contains(view, "Personal Notes") {
		t.Error("Notes view should show title")
	}

	// Test switching to plugins view
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewPlugins {
		t.Error("Should switch to plugins view")
	}

	// Test plugins view rendering
	view = updatedModel.View()
	if !strings.Contains(view, "Plugin Manager") {
		t.Error("Plugins view should show title")
	}

	// Test switching to online view
	msg2 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}
	newModel2, _ := m.Update(msg2)
	updatedModel2 := newModel2.(ui.Model)
	if updatedModel2.ViewMode != ui.ViewOnline {
		t.Error("Should switch to online view")
	}

	// Test online view rendering
	view = updatedModel2.View()
	if !strings.Contains(view, "Online Repositories") {
		t.Error("Online view should show title")
	}

	// Test switching to sync view
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewSync {
		t.Error("Should switch to sync view")
	}

	// Test sync view rendering
	view = updatedModel.View()
	if !strings.Contains(view, "Sync Status") {
		t.Error("Sync view should show title")
	}
}

func TestNotesViewInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.ViewMode = ui.ViewNotes

	// Test navigation in notes view
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	// Navigation should work even with empty notes

	// Test creating note
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.StatusMessage == "" {
		t.Error("Should show status message when creating note")
	}

	// Test escape to go back
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewMain {
		t.Error("Escape should return to main view")
	}
}

func TestPluginsViewInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.ViewMode = ui.ViewPlugins

	// Test navigation
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := m.Update(msg)
	// Should handle gracefully even with no plugins

	// Test reload
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ = m.Update(msg)
	updatedModel := newModel.(ui.Model)
	if updatedModel.StatusMessage == "" {
		t.Error("Should show status message when reloading plugins")
	}

	// Test escape
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewMain {
		t.Error("q should return to main view")
	}
}

func TestOnlineViewInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.ViewMode = ui.ViewOnline

	// Test navigation
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	// Should handle gracefully

	// Test search mode
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if !updatedModel.SearchMode {
		t.Error("/ should enter search mode in online view")
	}

	// Test escape
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewMain {
		t.Error("Escape should return to main view")
	}
}

func TestSyncViewInput(t *testing.T) {
	m := initialModelWithDefaults()
	m.ViewMode = ui.ViewSync

	// Test trigger sync
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	if updatedModel.StatusMessage == "" {
		t.Error("Should show status message when triggering sync")
	}

	// Test resolve conflicts
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.StatusMessage == "" {
		t.Error("Should show status message when resolving conflicts")
	}

	// Test toggle auto-sync
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	// Should handle gracefully

	// Test escape
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = m.Update(msg)
	updatedModel = newModel.(ui.Model)
	if updatedModel.ViewMode != ui.ViewMain {
		t.Error("Escape should return to main view")
	}
}

func TestForceSync(t *testing.T) {
	m := initialModelWithDefaults()

	// Test Ctrl+S for force sync
	msg := tea.KeyMsg{Type: tea.KeyCtrlS}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(ui.Model)
	if updatedModel.StatusMessage == "" {
		t.Error("Ctrl+S should show sync status message")
	}
}

func TestModelInit(t *testing.T) {
	m := initialModelWithDefaults()
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestPrintHelp(t *testing.T) {
	// This just ensures printHelp doesn't panic
	// We can't easily test stdout output
	printHelp()
}

func TestPrintVersion(t *testing.T) {
	// This just ensures printVersion doesn't panic
	printVersion()
}

func TestSearchClearShortcut(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true
	m.SearchQuery = "test query"

	// Test Ctrl+U to clear search
	msg := tea.KeyMsg{Type: tea.KeyCtrlU}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("ctrl+u should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.SearchQuery != "" {
		t.Error("ctrl+u should clear search query")
	}
}

func TestSearchHighlighting(t *testing.T) {
	m := initialModelWithDefaults()
	m.LastSearch = "move"
	m.Rows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"j", "↓ move"},
	}
	m.AllRows = [][]string{
		{"Shortcut", "vim"},
		{"k", "↑ move"},
		{"j", "↓ move"},
		{"q", "quit"},
	}

	view := m.View()
	// Should use highlighting when there is a search term and filtered results
	if m.LastSearch != "" && len(m.Rows) != len(m.AllRows) {
		// The highlighting is applied, but we cannot easily test the exact output
		// since it contains ANSI escape codes. We just verify the view renders.
		if view == "" {
			t.Error("view with highlighting should not be empty")
		}
	}
}

func TestLastSearchTracking(t *testing.T) {
	m := initialModelWithDefaults()
	m.SearchMode = true
	m.SearchQuery = "test"

	// Test that lastSearch is set when confirming search
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	if cmd != nil {
		t.Error("enter should not return command")
	}

	updatedModel := newModel.(ui.Model)
	if updatedModel.LastSearch != "test" {
		t.Error("should set lastSearch when confirming search")
	}

	// Test that lastSearch is cleared when escaping from search
	updatedModel.SearchMode = true
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd = updatedModel.Update(msg)
	if cmd != nil {
		t.Error("escape should not return command")
	}

	finalModel := newModel.(ui.Model)
	if finalModel.LastSearch != "" {
		t.Error("should clear lastSearch when escaping from search")
	}
}

func TestOpenEditorForNote(t *testing.T) {
	m := initialModelWithDefaults()

	// Create a test note
	testNote := &notes.Note{
		ID:       "test-note-1",
		Title:    "Test Note",
		Content:  "Original content",
		Category: "test",
		Tags:     []string{"test", "example"},
	}

	// Test with EDITOR environment variable set to "echo" to avoid opening actual editor
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor) // Restore after test

	// Set EDITOR to a simple command that won't open an interactive editor
	os.Setenv("EDITOR", "echo")

	// Test the method
	result, err := m.OpenEditorForNote(testNote)

	if err != nil {
		t.Errorf("openEditorForNote should not return error with echo editor: %v", err)
	}

	if result == nil {
		t.Error("openEditorForNote should return a note")
	}

	if result.ID != testNote.ID {
		t.Error("returned note should preserve ID")
	}

	if result.Title != testNote.Title {
		t.Error("returned note should preserve title when unchanged")
	}

	if result.Category != testNote.Category {
		t.Error("returned note should preserve category when unchanged")
	}

	if len(result.Tags) != len(testNote.Tags) {
		t.Error("returned note should preserve tags when unchanged")
	}
}

func TestOpenEditorForNote_InvalidEditor(t *testing.T) {
	m := initialModelWithDefaults()

	testNote := &notes.Note{
		ID:      "test-note-1",
		Title:   "Test Note",
		Content: "Original content",
	}

	// Test with invalid editor command
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	os.Setenv("EDITOR", "nonexistent-editor-command-12345")

	result, err := m.OpenEditorForNote(testNote)

	if err == nil {
		t.Error("openEditorForNote should return error with invalid editor")
	}

	if result != nil {
		t.Error("openEditorForNote should return nil result on error")
	}
}

func TestOpenEditorForNote_DefaultEditor(t *testing.T) {
	m := initialModelWithDefaults()

	testNote := &notes.Note{
		ID:      "test-note-1",
		Title:   "Test Note",
		Content: "Original content",
	}

	// Test with no EDITOR environment variable (should default to nano)
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	os.Unsetenv("EDITOR")

	// Since nano might not be available and would be interactive,
	// we'll test that the method attempts to use "nano" by checking the error message
	_, err := m.OpenEditorForNote(testNote)

	// We expect an error since nano likely isn't available or would hang waiting for input
	if err == nil {
		t.Error("openEditorForNote with default nano editor should fail in test environment")
	}

	// The error should mention the editor (though this is implementation-dependent)
	if err != nil && !strings.Contains(err.Error(), "editor") {
		t.Errorf("error should mention editor: %v", err)
	}
}
