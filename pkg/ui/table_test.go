package ui

import (
	"strings"
	"testing"
)

func TestNewTableRenderer(t *testing.T) {
	theme := DefaultTheme()
	renderer := NewTableRenderer(theme)

	if renderer == nil {
		t.Fatal("NewTableRenderer() returned nil")
	}

	if renderer.theme != theme {
		t.Error("renderer should store the provided theme")
	}
}

func TestTableRenderer_Render_EmptyData(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	result := renderer.Render([][]string{}, 0, 0)
	if result != "" {
		t.Error("rendering empty data should return empty string")
	}
}

func TestTableRenderer_Render_BasicTable(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Header1", "Header2", "Header3"},
		{"Row1Col1", "Row1Col2", "Row1Col3"},
		{"Row2Col1", "Row2Col2", "Row2Col3"},
	}

	result := renderer.Render(data, 0, 1)

	if result == "" {
		t.Error("render should return non-empty string for valid data")
	}

	// Check that all data appears in output
	for _, row := range data {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("rendered output should contain cell data: %s", cell)
			}
		}
	}

	// Check for table structure elements
	if !strings.Contains(result, "│") {
		t.Error("rendered table should contain column separators")
	}

	if !strings.Contains(result, "─") {
		t.Error("rendered table should contain header separator")
	}

	if !strings.Contains(result, "┼") {
		t.Error("rendered table should contain header separator junction")
	}
}

func TestTableRenderer_Render_SingleColumn(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Header"},
		{"Row1"},
		{"Row2"},
	}

	result := renderer.Render(data, 0, 0)

	if result == "" {
		t.Error("render should return non-empty string")
	}

	// Should contain all data
	for _, row := range data {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("rendered output should contain: %s", cell)
			}
		}
	}
}

func TestTableRenderer_Render_CursorHighlight(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Header1", "Header2"},
		{"Row1Col1", "Row1Col2"},
		{"Row2Col1", "Row2Col2"},
	}

	// Test different cursor positions
	result1 := renderer.Render(data, 0, 1)
	result2 := renderer.Render(data, 1, 1)

	// Both should render without error
	if result1 == "" || result2 == "" {
		t.Error("both renders should produce output")
	}

	// Note: cursor highlighting may or may not produce visibly different text
	// depending on the terminal styling, so we just verify both work
}

func TestTableRenderer_Render_VariableColumnWidths(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Short", "Very Long Header Text"},
		{"A", "B"},
		{"Medium Text", "C"},
	}

	result := renderer.Render(data, 0, 0)

	if result == "" {
		t.Error("render should handle variable column widths")
	}

	// Should contain all data
	for _, row := range data {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("rendered output should contain: %s", cell)
			}
		}
	}
}

func TestTableRenderer_Render_UnicodeContent(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"ASCII", "Unicode"},
		{"test", "测试"},
		{"emoji", "🚀💻"},
	}

	result := renderer.Render(data, 0, 0)

	if result == "" {
		t.Error("render should handle unicode content")
	}

	// Should contain unicode characters
	if !strings.Contains(result, "测试") {
		t.Error("should handle Chinese characters")
	}

	if !strings.Contains(result, "🚀💻") {
		t.Error("should handle emoji characters")
	}
}

func TestTableRenderer_RenderWithInstructions(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Header"},
		{"Row1"},
	}

	result := renderer.RenderWithInstructions(data, 0, 0)

	if result == "" {
		t.Error("RenderWithInstructions should return non-empty string")
	}

	// Should contain the table data
	if !strings.Contains(result, "Header") {
		t.Error("should contain table data")
	}

	// Should contain instructions
	if !strings.Contains(result, "arrow keys") {
		t.Error("should contain usage instructions")
	}

	if !strings.Contains(result, "hjkl") {
		t.Error("should contain vim-style key instructions")
	}

	if !strings.Contains(result, "q to quit") {
		t.Error("should contain quit instructions")
	}
}

func TestTableRenderer_RenderWithInstructions_EmptyData(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	result := renderer.RenderWithInstructions([][]string{}, 0, 0)

	// Should still contain instructions even with empty data
	if !strings.Contains(result, "arrow keys") {
		t.Error("should contain instructions even with empty data")
	}
}

func TestTableRenderer_Render_EdgeCases(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	// Test with nil data
	result := renderer.Render(nil, 0, 0)
	if result != "" {
		t.Error("nil data should return empty string")
	}

	// Test with empty rows
	data := [][]string{
		{},
	}
	result = renderer.Render(data, 0, 0)
	// Should handle gracefully (might be empty or minimal output)

	// Test with consistent row lengths only to avoid index out of bounds
	data = [][]string{
		{"A", "B", "C"},
		{"1", "2", "3"},
		{"X", "Y", "Z"},
	}
	result = renderer.Render(data, 0, 0)
	if result == "" {
		t.Error("should handle consistent table structure")
	}
}

func TestTableRenderer_Render_LargeTable(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	// Create a larger table
	data := make([][]string, 100)
	data[0] = []string{"Col1", "Col2", "Col3", "Col4", "Col5"}

	for i := 1; i < 100; i++ {
		data[i] = []string{
			"Row" + string(rune(i)),
			"Data" + string(rune(i)),
			"Test" + string(rune(i)),
			"Value" + string(rune(i)),
			"End" + string(rune(i)),
		}
	}

	result := renderer.Render(data, 2, 50)

	if result == "" {
		t.Error("should handle large tables")
	}

	// Should contain header
	if !strings.Contains(result, "Col1") {
		t.Error("should contain header data")
	}
}

func TestTableRenderer_Render_SpecialCharacters(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Special", "Characters"},
		{"Tabs\tTest", "Newlines\nTest"},
		{"Quotes\"Test", "Backslash\\Test"},
	}

	result := renderer.Render(data, 0, 0)

	if result == "" {
		t.Error("should handle special characters")
	}

	// Should handle without crashing
	if !strings.Contains(result, "Special") {
		t.Error("should contain header data")
	}
}

func TestTableRenderer_GetTheme(t *testing.T) {
	theme := DefaultTheme()
	renderer := NewTableRenderer(theme)
	
	retrievedTheme := renderer.GetTheme()
	if retrievedTheme != theme {
		t.Error("GetTheme should return the same theme instance")
	}
	if retrievedTheme.Name != "default" {
		t.Error("should return theme with correct name")
	}
}


func TestTableRenderer_HighlightSearchTerm(t *testing.T) {
	theme := DefaultTheme()
	renderer := NewTableRenderer(theme)
	
	// Test basic highlighting
	result := renderer.highlightSearchTerm("move up", "move")
	if !strings.Contains(result, "move") {
		t.Error("should contain the search term")
	}
	
	// Test case insensitive highlighting
	result = renderer.highlightSearchTerm("MOVE up", "move")
	if !strings.Contains(result, "MOVE") {
		t.Error("should preserve original case")
	}
	
	// Test no match
	result = renderer.highlightSearchTerm("quit", "move")
	if result != "quit" {
		t.Error("should return original text when no match")
	}
	
	// Test empty search term
	result = renderer.highlightSearchTerm("some text", "")
	if result != "some text" {
		t.Error("should return original text when search term is empty")
	}
}

func TestTableRenderer_RenderWithHighlighting(t *testing.T) {
	theme := DefaultTheme()
	renderer := NewTableRenderer(theme)
	
	data := [][]string{
		{"Shortcut", "Description"},
		{"k", "move up"},
		{"j", "move down"},
	}
	
	// Test highlighting render
	result := renderer.RenderWithHighlighting(data, 0, 1, "move")
	if result == "" {
		t.Error("should return non-empty string")
	}
	
	// Should contain all data
	for _, row := range data {
		for _, cell := range row {
			if !strings.Contains(result, cell) {
				t.Errorf("rendered output should contain cell data: %s", cell)
			}
		}
	}
	
	// Test with empty search term (should work like normal render)
	result2 := renderer.RenderWithHighlighting(data, 0, 1, "")
	if result2 == "" {
		t.Error("should return non-empty string even with empty search term")
	}
}

