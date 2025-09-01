package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) HandleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "ctrl+[":
		m.SearchMode = false
		m.SearchQuery = ""
		m.Rows = m.AllRows
		m.LastSearch = ""
		m.CursorY = 1
		return m, nil
	case "ctrl+u":
		m.SearchQuery = ""
		return m, nil
	case "enter":
		m.SearchMode = false
		if m.SearchQuery == "" {
			m.Rows = m.AllRows
			m.LastSearch = ""
		} else {
			m.Rows = m.Registry.SearchTableData(m.Config.Apps, m.SearchQuery)
			m.LastSearch = m.SearchQuery
		}
		m.CursorY = 1
		return m, nil
	case "backspace":
		if len(m.SearchQuery) > 0 {
			m.SearchQuery = m.SearchQuery[:len(m.SearchQuery)-1]
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.SearchQuery += msg.String()
		}
		return m, nil
	}
}

func (m Model) HandleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "ctrl+[":
		m.FilterMode = false
		return m, nil
	case "ctrl+u":
		m.FilteredApps = make([]string, 0)
		return m, nil
	case "enter":
		m.FilterMode = false
		if len(m.FilteredApps) == 0 {
			m.Rows = m.Registry.GetTableData(m.AllApps)
		} else {
			m.Rows = m.Registry.GetTableData(m.FilteredApps)
		}
		m.CursorY = 1
		return m, nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		appIndex := int(msg.String()[0] - '1')
		if appIndex < len(m.AllApps) {
			appName := m.AllApps[appIndex]
			found := false
			for i, name := range m.FilteredApps {
				if name == appName {
					m.FilteredApps = append(m.FilteredApps[:i], m.FilteredApps[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				m.FilteredApps = append(m.FilteredApps, appName)
			}
		}
		return m, nil
	case "c":
		m.FilteredApps = make([]string, 0)
		return m, nil
	case "a":
		m.FilteredApps = make([]string, len(m.AllApps))
		copy(m.FilteredApps, m.AllApps)
		return m, nil
	}
	return m, nil
}

func (m Model) HandleHelpInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?", "esc", "q":
		m.HelpMode = false
		m.ViewMode = ViewMain
		return m, nil
	}
	return m, nil
}

func (m Model) HandlePluginsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ViewMode = ViewMain
		return m, nil
	case "up", "k":
		if m.PluginCursor > 0 {
			m.PluginCursor--
		}
		return m, nil
	case "down", "j":
		if m.PluginCursor < len(m.PluginsList)-1 {
			m.PluginCursor++
		}
		return m, nil
	case "l":
		if m.PluginCursor < len(m.PluginsList) {
			plugin := m.PluginsList[m.PluginCursor]
			m.StatusMessage = fmt.Sprintf("Loading plugin: %s", plugin.Metadata.Name)
		}
		return m, nil
	case "u":
		if m.PluginCursor < len(m.PluginsList) {
			plugin := m.PluginsList[m.PluginCursor]
			m.PluginLoader.UnloadPlugin(plugin.Metadata.Name)
			m.LoadPlugins()
			m.StatusMessage = fmt.Sprintf("Unloaded plugin: %s", plugin.Metadata.Name)
		}
		return m, nil
	case "r":
		m.PluginLoader.LoadAll()
		m.LoadPlugins()
		m.StatusMessage = "Plugins reloaded"
		return m, nil
	}
	return m, nil
}

func (m Model) HandleOnlineInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ViewMode = ViewMain
		return m, nil
	case "up", "k":
		if m.RepoCursor > 0 {
			m.RepoCursor--
		}
		return m, nil
	case "down", "j":
		if m.RepoCursor < len(m.ReposList)-1 {
			m.RepoCursor++
		}
		return m, nil
	case "enter":
		if m.RepoCursor < len(m.ReposList) {
			repo := m.ReposList[m.RepoCursor]
			m.LoadCheatSheets(repo.URL)
		}
		return m, nil
	case "d":
		if m.SheetCursor < len(m.CheatSheets) {
			sheet := m.CheatSheets[m.SheetCursor]
			m.StatusMessage = fmt.Sprintf("Downloading: %s", sheet.Name)
		}
		return m, nil
	case "/":
		m.SearchMode = true
		return m, nil
	}
	return m, nil
}

func (m Model) HandleSyncInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ViewMode = ViewMain
		return m, nil
	case "s":
		m.StatusMessage = "Syncing..."
		if m.SyncManager != nil {
			go m.SyncManager.Sync()
		}
		return m, nil
	case "r":
		m.StatusMessage = "Resolving conflicts..."
		return m, nil
	case "a":
		if m.SyncManager != nil {
			m.StatusMessage = "Auto-sync toggled"
		}
		return m, nil
	}
	return m, nil
}
