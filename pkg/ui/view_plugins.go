package ui

import (
	"fmt"
	"strings"
)

func (m Model) ViewPlugins() string {
	var output strings.Builder

	output.WriteString("╭─ Plugin Manager ─────────────────────────────────────────╮\n")

	if len(m.PluginsList) == 0 {
		output.WriteString("│  No plugins loaded. Place plugins in plugins directory. │\n")
	} else {
		for i, plugin := range m.PluginsList {
			if i > 10 {
				output.WriteString(fmt.Sprintf("│  ... and %d more plugins                                 │\n", len(m.PluginsList)-10))
				break
			}

			cursor := "  "
			if i == m.PluginCursor {
				cursor = "▶ "
			}

			line := fmt.Sprintf("%s%-20s v%-8s %s", cursor, plugin.Metadata.Name, plugin.Metadata.Version, plugin.Metadata.Author)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: l: load • u: unload • r: reload all • esc: back\n")

	if m.StatusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.StatusMessage))
	}

	return output.String()
}
