package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) ViewMain() string {
	var output strings.Builder

	tableStr := m.Renderer.RenderWithHighlighting(
		m.Rows,
		m.CursorX,
		m.CursorY,
		m.LastSearch,
	)
	output.WriteString(tableStr)
	output.WriteString("\n")

	if m.SearchMode {
		output.WriteString(fmt.Sprintf("\nSearch: %s_\nType to search, Enter to confirm, Esc to cancel\n", m.SearchQuery))
	} else if m.FilterMode {
		output.WriteString("\nFilter Apps: ")
		for i, app := range m.AllApps {
			isSelected := false
			for _, selected := range m.FilteredApps {
				if app == selected {
					isSelected = true
					break
				}
			}
			if isSelected {
				output.WriteString(fmt.Sprintf(" [%d] ✓%s", i+1, app))
			} else {
				output.WriteString(fmt.Sprintf(" [%d] %s", i+1, app))
			}
		}
		output.WriteString("\n1-9: toggle apps, a: all, c: clear, Enter: apply, Esc: cancel\n")
	} else {
		output.WriteString("\nArrow keys/hjkl: move • /: search • f: filter • n: notes • p: plugins • o: online • s: sync • ?: help • q: quit\n")
	}

	if m.StatusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.StatusMessage))
	}

	return output.String()
}

func (m Model) HandleMainInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "?":
		m.HelpMode = true
		m.ViewMode = ViewHelp
		return m, nil
	case "/":
		m.SearchMode = true
		return m, nil
	case "f", "ctrl+f":
		m.FilterMode = true
		return m, nil
	case "n":
		m.ViewMode = ViewNotes
		m.LoadNotes()
		return m, nil
	case "p":
		m.ViewMode = ViewPlugins
		m.LoadPlugins()
		return m, nil
	case "o":
		m.ViewMode = ViewOnline
		m.LoadRepositories()
		return m, nil
	case "s":
		m.ViewMode = ViewSync
		m.LoadSyncStatus()
		return m, nil
	case "ctrl+s":
		m.StatusMessage = "Syncing..."
		return m, nil
	case "up", "k":
		if m.CursorY > 1 {
			m.CursorY--
		}
		return m, nil
	case "down", "j":
		if m.CursorY < len(m.Rows)-1 {
			m.CursorY++
		}
		return m, nil
	case "left", "h":
		if m.CursorX > 0 {
			m.CursorX--
		}
		return m, nil
	case "right", "l":
		if m.CursorX < len(m.Rows[0])-1 {
			m.CursorX++
		}
		return m, nil
	case "ctrl+a", "home":
		m.CursorY = 1
		return m, nil
	case "ctrl+e", "end":
		m.CursorY = len(m.Rows) - 1
		return m, nil
	case "ctrl+r":
		m.Rows = m.Registry.GetTableData(m.Config.Apps)
		m.AllRows = m.Rows
		return m, nil
	case "esc", "ctrl+[":
		m.SearchMode = false
		m.SearchQuery = ""
		m.LastSearch = ""
		m.Rows = m.AllRows
		m.CursorY = 1
		return m, nil
	}
	return m, nil
}
