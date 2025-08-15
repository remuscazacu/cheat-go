package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// TableRenderer handles the rendering of tabular data
type TableRenderer struct {
	theme *Theme
}

// NewTableRenderer creates a new table renderer with the given theme
func NewTableRenderer(theme *Theme) *TableRenderer {
	return &TableRenderer{
		theme: theme,
	}
}

// Render renders a table from the given data with cursor position
func (r *TableRenderer) Render(rows [][]string, cursorX, cursorY int) string {
	if len(rows) == 0 {
		return ""
	}

	var b strings.Builder

	// Determine column widths using runewidth
	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if w := runewidth.StringWidth(cell); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Render rows
	for y, row := range rows {
		for x, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			pad := colWidths[x] - cellWidth
			content := " " + cell + strings.Repeat(" ", pad) + " "

			style := r.theme.CellStyle
			if y == 0 {
				style = r.theme.HeaderStyle
			}
			if x == cursorX && y == cursorY {
				style = style.Reverse(true)
			}

			b.WriteString(style.Render(content))
			if x < len(row)-1 {
				b.WriteString("│")
			}
		}
		b.WriteString("\n")

		// Add separator after header
		if y == 0 {
			for i, w := range colWidths {
				b.WriteString(strings.Repeat("─", w+2))
				if i < len(colWidths)-1 {
					b.WriteString("┼")
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// GetTheme returns the current theme
func (r *TableRenderer) GetTheme() *Theme {
	return r.theme
}

// highlightSearchTerm highlights search terms in the given text
func (r *TableRenderer) highlightSearchTerm(text, searchTerm string) string {
	if searchTerm == "" {
		return text
	}

	// Case-insensitive highlighting
	lowerText := strings.ToLower(text)
	lowerTerm := strings.ToLower(searchTerm)

	if !strings.Contains(lowerText, lowerTerm) {
		return text
	}

	// Find the position of the match
	index := strings.Index(lowerText, lowerTerm)
	if index == -1 {
		return text
	}

	// Extract the actual case-preserved match
	beforeMatch := text[:index]
	actualMatch := text[index : index+len(searchTerm)]
	afterMatch := text[index+len(searchTerm):]

	// Apply highlighting style
	highlightedMatch := r.theme.HighlightStyle.Render(actualMatch)

	// Recursively highlight remaining occurrences
	highlightedAfter := r.highlightSearchTerm(afterMatch, searchTerm)

	return beforeMatch + highlightedMatch + highlightedAfter
}

// RenderWithHighlighting renders the table with search term highlighting
func (r *TableRenderer) RenderWithHighlighting(rows [][]string, cursorX, cursorY int, searchTerm string) string {
	if len(rows) == 0 {
		return ""
	}

	var b strings.Builder

	// Determine column widths using runewidth (without highlight markup)
	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if w := runewidth.StringWidth(cell); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Render rows with highlighting
	for y, row := range rows {
		for x, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			pad := colWidths[x] - cellWidth

			// Apply highlighting if not header row and search term exists
			content := cell
			if y > 0 && searchTerm != "" {
				content = r.highlightSearchTerm(cell, searchTerm)
			}

			contentWithPadding := " " + content + strings.Repeat(" ", pad) + " "

			style := r.theme.CellStyle
			if y == 0 {
				style = r.theme.HeaderStyle
			}
			if x == cursorX && y == cursorY {
				style = style.Copy().Inherit(r.theme.SelectedRowStyle)
			}

			b.WriteString(style.Render(contentWithPadding))
			if x < len(row)-1 {
				b.WriteString("│")
			}
		}
		b.WriteString("\n")

		// Add separator after header
		if y == 0 {
			for i, w := range colWidths {
				b.WriteString(strings.Repeat("─", w+2))
				if i < len(colWidths)-1 {
					b.WriteString("┼")
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// RenderWithInstructions renders the table with usage instructions
func (r *TableRenderer) RenderWithInstructions(rows [][]string, cursorX, cursorY int) string {
	table := r.Render(rows, cursorX, cursorY)
	return table + "\nUse arrow keys or hjkl to move. Press q to quit."
}
