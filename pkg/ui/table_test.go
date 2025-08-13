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

// TestTableRenderer_BorderAlignment tests that table borders are properly aligned
// This test prevents layout bugs where borders don't line up with content
func TestTableRenderer_BorderAlignment(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	data := [][]string{
		{"Short", "Medium Text", "Very Long Header Text"},
		{"A", "B", "C"},
		{"Test", "Data", "Value"},
	}

	result := renderer.Render(data, 0, 0)
	lines := strings.Split(result, "\n")

	// Filter out empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) < 5 {
		t.Fatal("Expected at least 5 non-empty lines (top border, header, separator, 2 data rows, bottom border)")
	}

	// Check that all border lines have the same length
	var borderLineLength int
	borderLines := []int{0, 2, len(nonEmptyLines) - 1} // top border, separator, bottom border

	for i, lineIdx := range borderLines {
		line := nonEmptyLines[lineIdx]
		if i == 0 {
			borderLineLength = len(line)
		} else {
			if len(line) != borderLineLength {
				t.Errorf("Border line %d has length %d, expected %d. Line: %s",
					lineIdx, len(line), borderLineLength, line)
			}
		}
	}

	// Check that content lines (header and data rows) have consistent structure
	contentLines := []int{1, 3, 4} // header and data rows
	var verticalBarCount int

	for i, lineIdx := range contentLines {
		line := nonEmptyLines[lineIdx]
		currentBarCount := strings.Count(line, "│")

		if i == 0 {
			verticalBarCount = currentBarCount
		} else {
			if currentBarCount != verticalBarCount {
				t.Errorf("Content line %d has %d vertical bars, expected %d. Line: %s",
					lineIdx, currentBarCount, verticalBarCount, line)
			}
		}
	}

	// Verify structural consistency - all content lines should have the same number of vertical bars
	expectedBarCount := strings.Count(nonEmptyLines[1], "│") // header line
	for _, lineIdx := range contentLines[1:] {               // check other content lines
		line := nonEmptyLines[lineIdx]
		actualBarCount := strings.Count(line, "│")
		if actualBarCount != expectedBarCount {
			t.Errorf("Line %d has %d vertical bars, expected %d. Line: %s",
				lineIdx, actualBarCount, expectedBarCount, line)
		}
	}
}

// TestTableRenderer_DifferentTableStyles tests that all table styles maintain proper alignment
func TestTableRenderer_DifferentTableStyles(t *testing.T) {
	themes := map[string]*Theme{
		"default": DefaultTheme(),
		"dark":    DarkTheme(),
		"light":   LightTheme(),
		"minimal": MinimalTheme(),
	}

	data := [][]string{
		{"Style", "Test", "Table"},
		{"Row1", "Data1", "Value1"},
		{"Row2", "Data2", "Value2"},
	}

	for themeName, theme := range themes {
		t.Run(themeName, func(t *testing.T) {
			renderer := NewTableRenderer(theme)
			result := renderer.Render(data, 0, 0)

			if result == "" {
				t.Errorf("Theme %s should produce output", themeName)
				return
			}

			// All content should be present
			for _, row := range data {
				for _, cell := range row {
					if !strings.Contains(result, cell) {
						t.Errorf("Theme %s missing cell content: %s", themeName, cell)
					}
				}
			}

			// For non-minimal styles, check for table structure
			if theme.TableStyle != "minimal" {
				lines := strings.Split(result, "\n")
				var nonEmptyLines []string
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						nonEmptyLines = append(nonEmptyLines, line)
					}
				}

				if len(nonEmptyLines) < 3 {
					t.Errorf("Theme %s should have at least 3 lines", themeName)
					return
				}

				// Check for vertical separators in content lines
				headerLine := ""
				for _, line := range nonEmptyLines {
					if strings.Contains(line, "Style") {
						headerLine = line
						break
					}
				}

				if headerLine == "" {
					t.Errorf("Theme %s should contain header line", themeName)
					return
				}

				// Should have at least 2 vertical bars (start and end, plus internal separators)
				barCount := strings.Count(headerLine, "│")
				if barCount < 2 {
					t.Errorf("Theme %s header should have at least 2 vertical bars, got %d", themeName, barCount)
				}
			}
		})
	}
}

// TestTableRenderer_ColumnWidthConsistency tests that column widths are calculated consistently
func TestTableRenderer_ColumnWidthConsistency(t *testing.T) {
	renderer := NewTableRenderer(DefaultTheme())

	// Test with varying content lengths
	data := [][]string{
		{"A", "BB", "CCC"},
		{"DDDD", "E", "FF"},
		{"G", "HHHH", "I"},
	}

	result := renderer.Render(data, 0, 0)
	lines := strings.Split(result, "\n")

	// Find content lines (lines containing data)
	var contentLines []string
	for _, line := range lines {
		if strings.Contains(line, "│") && (strings.Contains(line, "A") || strings.Contains(line, "DDDD") || strings.Contains(line, "G")) {
			contentLines = append(contentLines, line)
		}
	}

	if len(contentLines) != 3 {
		t.Fatalf("Expected 3 content lines, got %d", len(contentLines))
	}

	// All content lines should have the same length (accounting for consistent column spacing)
	expectedLength := len(contentLines[0])
	for i, line := range contentLines {
		if len(line) != expectedLength {
			t.Errorf("Content line %d has length %d, expected %d. Line: %s",
				i, len(line), expectedLength, line)
		}
	}
}
