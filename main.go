package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	cellStyle   = lipgloss.NewStyle().Padding(0, 0)
)

type model struct {
	rows    [][]string
	cursorX int
	cursorY int
}

func initialModel() model {
	return model{
		rows: [][]string{
			{"Shortcut", "vim", "zsh", "dwm", "st", "lf", "zathura"},
			{"h", "← move", "back char", "focus left", "← move", "left", "scroll ←"},
			{"l", "→ move", "forward char", "focus right", "→ move", "right", "scroll →"},
			{"j", "↓ move", "down history", "focus down", "↓ scroll", "down", "scroll ↓"},
			{"k", "↑ move", "up history", "focus up", "↑ scroll", "up", "scroll ↑"},
			{"gg", "top", "-", "-", "-", "top", "-"},
			{"G", "bottom", "-", "-", "-", "bottom", "-"},
			{"/", "search", "search history", "-", "search", "search", "search"},
			{":", "command", "prompt", "command", "-", "command", "-"},
			{"q", "quit", "exit", "close win", "exit", "quit", "quit"},
		},
		cursorX: 0,
		cursorY: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursorY > 1 {
				m.cursorY--
			}
		case "down", "j":
			if m.cursorY < len(m.rows)-1 {
				m.cursorY++
			}
		case "left", "h":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "right", "l":
			if m.cursorX < len(m.rows[0])-1 {
				m.cursorX++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Determine column widths using runewidth
	colWidths := make([]int, len(m.rows[0]))
	for _, row := range m.rows {
		for i, cell := range row {
			if w := runewidth.StringWidth(cell); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	// Render rows
	for y, row := range m.rows {
		for x, cell := range row {
			cellWidth := runewidth.StringWidth(cell)
			pad := colWidths[x] - cellWidth
			content := " " + cell + strings.Repeat(" ", pad) + " "

			style := cellStyle
			if y == 0 {
				style = headerStyle
			}
			if x == m.cursorX && y == m.cursorY {
				style = style.Reverse(true)
			}

			b.WriteString(style.Render(content))
			if x < len(row)-1 {
				b.WriteString("│")
			}
		}
		b.WriteString("\n")

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

	b.WriteString("\nUse arrow keys or hjkl to move. Press q to quit.")
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
