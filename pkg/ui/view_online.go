package ui

import (
	"fmt"
	"strings"
)

func (m Model) ViewOnline() string {
	var output strings.Builder

	output.WriteString("╭─ Online Repositories ────────────────────────────────────╮\n")

	if len(m.ReposList) == 0 {
		output.WriteString("│  Loading repositories...                                 │\n")
	} else {
		for i, repo := range m.ReposList {
			cursor := "  "
			if i == m.RepoCursor {
				cursor = "▶ "
			}

			line := fmt.Sprintf("%s%-30s ⭐%d", cursor, repo.Name, repo.Stars)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	if len(m.CheatSheets) > 0 {
		output.WriteString("│──────────────────────────────────────────────────────────│\n")
		output.WriteString("│ Cheat Sheets:                                            │\n")
		for i, sheet := range m.CheatSheets {
			if i > 5 {
				break
			}
			line := fmt.Sprintf("  %-25s ⬇%d ★%.1f", sheet.Name, sheet.Downloads, sheet.Rating)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: enter: browse • d: download • /: search • esc: back\n")

	if m.StatusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.StatusMessage))
	}

	return output.String()
}
