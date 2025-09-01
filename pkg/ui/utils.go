package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"cheat-go/pkg/notes"
	"cheat-go/pkg/online"
	"cheat-go/pkg/sync"
)

func (m *Model) LoadNotes() {
	notes, _ := m.NotesManager.ListNotes()
	m.NotesList = notes
	m.NoteCursor = 0
}

func (m *Model) LoadPlugins() {
	m.PluginsList = m.PluginLoader.ListPlugins()
	m.PluginCursor = 0
}

func (m *Model) LoadRepositories() {
	repos, _ := m.OnlineClient.GetRepositories()
	m.ReposList = repos
	m.RepoCursor = 0
}

func (m *Model) LoadCheatSheets(repoURL string) {
	sheets, _ := m.OnlineClient.SearchCheatSheets(online.SearchOptions{
		Repository: repoURL,
		Limit:      50,
	})
	m.CheatSheets = sheets
	m.SheetCursor = 0
}

func (m *Model) LoadSyncStatus() {
	if m.SyncManager != nil {
		m.SyncStatus = m.SyncManager.GetSyncStatus()
	} else {
		m.SyncStatus = sync.SyncStatus{
			LastSync:  time.Time{},
			IsSyncing: false,
			DeviceID:  "not-configured",
		}
	}
}

func (m Model) FilterRowsBySearch(query string) [][]string {
	if query == "" {
		return m.AllRows
	}

	var filtered [][]string
	if len(m.AllRows) > 0 {
		filtered = append(filtered, m.AllRows[0])
	}

	for i := 1; i < len(m.AllRows); i++ {
		row := m.AllRows[i]
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), strings.ToLower(query)) {
				filtered = append(filtered, row)
				break
			}
		}
	}

	return filtered
}

func (m Model) IsAppSelected(appName string) bool {
	for _, name := range m.FilteredApps {
		if name == appName {
			return true
		}
	}
	return false
}

func (m Model) OpenEditorForNote(note *notes.Note) (*notes.Note, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	tmpFile, err := ioutil.TempFile("", "cheat-go-note-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := fmt.Sprintf("# Title: %s\n# Category: %s\n# Tags: %s\n\n%s",
		note.Title, note.Category, strings.Join(note.Tags, ", "), note.Content)

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor exited with error: %v", err)
	}

	editedContent, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read edited content: %v", err)
	}

	lines := strings.Split(string(editedContent), "\n")
	updatedNote := &notes.Note{
		ID:       note.ID,
		Title:    note.Title,
		Category: note.Category,
		Tags:     note.Tags,
		Content:  note.Content,
	}

	var contentStart int
	for i, line := range lines {
		if strings.HasPrefix(line, "# Title: ") {
			updatedNote.Title = strings.TrimSpace(strings.TrimPrefix(line, "# Title: "))
		} else if strings.HasPrefix(line, "# Category: ") {
			updatedNote.Category = strings.TrimSpace(strings.TrimPrefix(line, "# Category: "))
		} else if strings.HasPrefix(line, "# Tags: ") {
			tagStr := strings.TrimSpace(strings.TrimPrefix(line, "# Tags: "))
			if tagStr != "" {
				updatedNote.Tags = strings.Split(tagStr, ", ")
				for j := range updatedNote.Tags {
					updatedNote.Tags[j] = strings.TrimSpace(updatedNote.Tags[j])
				}
			} else {
				updatedNote.Tags = []string{}
			}
		} else if line == "" && i > 0 && strings.HasPrefix(lines[i-1], "# Tags: ") {
			contentStart = i + 1
			break
		}
	}

	if contentStart < len(lines) {
		updatedNote.Content = strings.Join(lines[contentStart:], "\n")
	}

	return updatedNote, nil
}
