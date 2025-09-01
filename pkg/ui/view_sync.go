package ui

import (
	"fmt"
	"strings"
)

func (m Model) ViewSync() string {
	var output strings.Builder

	output.WriteString("╭─ Sync Status ────────────────────────────────────────────╮\n")

	if m.SyncManager == nil {
		output.WriteString("│  Sync is not configured.                                 │\n")
		output.WriteString("│  Configure sync in ~/.config/cheat-go/config.yaml        │\n")
	} else {
		status := "Idle"
		if m.SyncStatus.IsSyncing {
			status = "Syncing..."
		}

		lastSync := "Never"
		if !m.SyncStatus.LastSync.IsZero() {
			lastSync = m.SyncStatus.LastSync.Format("2006-01-02 15:04:05")
		}

		output.WriteString(fmt.Sprintf("│  Status:     %-43s │\n", status))
		output.WriteString(fmt.Sprintf("│  Last Sync:  %-43s │\n", lastSync))
		output.WriteString(fmt.Sprintf("│  Device ID:  %-43s │\n", m.SyncStatus.DeviceID[:16]+"..."))

		if m.SyncStatus.HasConflicts {
			output.WriteString(fmt.Sprintf("│  ⚠ Conflicts: %-42d │\n", len(m.SyncStatus.Conflicts)))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: s: sync now • r: resolve conflicts • a: auto-sync • esc: back\n")

	if m.StatusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.StatusMessage))
	}

	return output.String()
}
