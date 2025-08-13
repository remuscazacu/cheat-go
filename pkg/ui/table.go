package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// TableRenderer handles the rendering of tabular data
type TableRenderer struct {
	theme      *Theme
	tableStyle string
	maxWidth   int
}

// TableBorders defines border characters for different table styles
type TableBorders struct {
	Horizontal  string
	Vertical    string
	Cross       string
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	TopCross    string
	BottomCross string
	LeftCross   string
	RightCross  string
}

// NewTableRenderer creates a new table renderer with the given theme
func NewTableRenderer(theme *Theme) *TableRenderer {
	return &TableRenderer{
		theme:      theme,
		tableStyle: theme.TableStyle,
		maxWidth:   120, // default max width
	}
}

// SetTableStyle sets the table style
func (r *TableRenderer) SetTableStyle(style string) {
	r.tableStyle = style
}

// SetMaxWidth sets the maximum table width
func (r *TableRenderer) SetMaxWidth(width int) {
	r.maxWidth = width
}

// getBorders returns the appropriate border characters for the current style
func (r *TableRenderer) getBorders() TableBorders {
	switch r.tableStyle {
	case "rounded":
		return TableBorders{
			Horizontal:  "─",
			Vertical:    "│",
			Cross:       "┼",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
			TopCross:    "┬",
			BottomCross: "┴",
			LeftCross:   "├",
			RightCross:  "┤",
		}
	case "bold":
		return TableBorders{
			Horizontal:  "━",
			Vertical:    "┃",
			Cross:       "╋",
			TopLeft:     "┏",
			TopRight:    "┓",
			BottomLeft:  "┗",
			BottomRight: "┛",
			TopCross:    "┳",
			BottomCross: "┻",
			LeftCross:   "┣",
			RightCross:  "┫",
		}
	case "minimal":
		return TableBorders{
			Horizontal:  " ",
			Vertical:    " ",
			Cross:       " ",
			TopLeft:     " ",
			TopRight:    " ",
			BottomLeft:  " ",
			BottomRight: " ",
			TopCross:    " ",
			BottomCross: " ",
			LeftCross:   " ",
			RightCross:  " ",
		}
	default: // simple
		return TableBorders{
			Horizontal:  "─",
			Vertical:    "│",
			Cross:       "┼",
			TopLeft:     "─",
			TopRight:    "─",
			BottomLeft:  "─",
			BottomRight: "─",
			TopCross:    "┬",
			BottomCross: "┴",
			LeftCross:   "├",
			RightCross:  "┤",
		}
	}
}

// Render renders a table from the given data with cursor position
func (r *TableRenderer) Render(rows [][]string, cursorX, cursorY int) string {
	if len(rows) == 0 {
		return ""
	}

	borders := r.getBorders()
	var b strings.Builder

	// Determine column widths using runewidth
	colWidths := make([]int, len(rows[0]))
	totalWidth := 0
	for _, row := range rows {
		for i, cell := range row {
			if w := runewidth.StringWidth(cell); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Calculate total width and adjust if necessary
	for i, w := range colWidths {
		totalWidth += w + 2 // padding
		if i < len(colWidths)-1 {
			totalWidth += 1 // separator
		}
	}

	// Adjust column widths if exceeding max width
	if totalWidth > r.maxWidth {
		excessWidth := totalWidth - r.maxWidth
		avgReduction := excessWidth / len(colWidths)
		for i := range colWidths {
			colWidths[i] = max(colWidths[i]-avgReduction, 10) // minimum column width
		}
	}

	// Render top border for non-minimal styles
	if r.tableStyle != "minimal" {
		r.renderTopBorder(&b, colWidths, borders)
	}

	// Render rows
	for y, row := range rows {
		r.renderRow(&b, row, colWidths, borders, cursorX, cursorY, y)

		// Add separator after header
		if y == 0 && r.tableStyle != "minimal" {
			r.renderSeparator(&b, colWidths, borders)
		}
	}

	// Render bottom border for non-minimal styles
	if r.tableStyle != "minimal" {
		r.renderBottomBorder(&b, colWidths, borders)
	}

	return b.String()
}

// renderTopBorder renders the top border of the table
func (r *TableRenderer) renderTopBorder(b *strings.Builder, colWidths []int, borders TableBorders) {
	b.WriteString(borders.TopLeft)
	for i, w := range colWidths {
		b.WriteString(strings.Repeat(borders.Horizontal, w+2))
		if i < len(colWidths)-1 {
			b.WriteString(borders.TopCross)
		}
	}
	b.WriteString(borders.TopRight + "\n")
}

// renderBottomBorder renders the bottom border of the table
func (r *TableRenderer) renderBottomBorder(b *strings.Builder, colWidths []int, borders TableBorders) {
	b.WriteString(borders.BottomLeft)
	for i, w := range colWidths {
		b.WriteString(strings.Repeat(borders.Horizontal, w+2))
		if i < len(colWidths)-1 {
			b.WriteString(borders.BottomCross)
		}
	}
	b.WriteString(borders.BottomRight + "\n")
}

// renderSeparator renders a separator line
func (r *TableRenderer) renderSeparator(b *strings.Builder, colWidths []int, borders TableBorders) {
	b.WriteString(borders.LeftCross)
	for i, w := range colWidths {
		b.WriteString(strings.Repeat(borders.Horizontal, w+2))
		if i < len(colWidths)-1 {
			b.WriteString(borders.Cross)
		}
	}
	b.WriteString(borders.RightCross + "\n")
}

// renderRow renders a single table row
func (r *TableRenderer) renderRow(b *strings.Builder, row []string, colWidths []int, borders TableBorders, cursorX, cursorY, y int) {
	if r.tableStyle != "minimal" {
		b.WriteString(borders.Vertical)
	}

	for x, cell := range row {
		cellWidth := runewidth.StringWidth(cell)
		pad := colWidths[x] - cellWidth
		content := " " + cell + strings.Repeat(" ", pad) + " "

		style := r.theme.CellStyle
		if y == 0 {
			style = r.theme.HeaderStyle
		}
		if x == cursorX && y == cursorY {
			style = r.theme.SelectedRowStyle.Copy().Inherit(style)
		}

		b.WriteString(style.Render(content))
		if x < len(row)-1 {
			if r.tableStyle != "minimal" {
				b.WriteString(borders.Vertical)
			} else {
				b.WriteString("  ") // spacing for minimal style
			}
		}
	}

	if r.tableStyle != "minimal" {
		b.WriteString(borders.Vertical)
	}
	b.WriteString("\n")
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RenderWithInstructions renders the table with usage instructions
func (r *TableRenderer) RenderWithInstructions(rows [][]string, cursorX, cursorY int) string {
	table := r.Render(rows, cursorX, cursorY)
	return table + "\nUse arrow keys or hjkl to move. Press q to quit."
}
