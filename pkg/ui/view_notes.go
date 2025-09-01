package ui

import (
	"fmt"
	"strings"
	"time"

	"cheat-go/pkg/notes"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) ViewNotes() string {
	var output strings.Builder

	output.WriteString("╭─ Personal Notes ─────────────────────────────────────────╮\n")

	if len(m.NotesList) == 0 {
		output.WriteString("│  No notes found. Press 'n' to create a new note.        │\n")
	} else {
		for i, note := range m.NotesList {
			if i > 10 {
				output.WriteString(fmt.Sprintf("│  ... and %d more notes                                   │\n", len(m.NotesList)-10))
				break
			}

			cursor := "  "
			if i == m.NoteCursor {
				cursor = "▶ "
			}

			favorite := " "
			if note.IsFavorite {
				favorite = "⭐"
			}

			line := fmt.Sprintf("%s%s %-30s %s", cursor, favorite, note.Title, note.AppName)
			if len(line) > 58 {
				line = line[:58]
			}
			output.WriteString(fmt.Sprintf("│%-58s│\n", line))
		}
	}

	output.WriteString("╰──────────────────────────────────────────────────────────╯\n")
	output.WriteString("\nKeys: n: new • e: edit • d: delete • f: favorite • esc: back\n")

	if m.StatusMessage != "" {
		output.WriteString(fmt.Sprintf("\nStatus: %s\n", m.StatusMessage))
	}

	return output.String()
}

func (m Model) HandleNotesInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ViewMode = ViewMain
		return m, nil
	case "up", "k":
		if m.NoteCursor > 0 {
			m.NoteCursor--
		}
		return m, nil
	case "down", "j":
		if m.NoteCursor < len(m.NotesList)-1 {
			m.NoteCursor++
		}
		return m, nil
	case "n":
		newNote := &notes.Note{
			Title:    fmt.Sprintf("New Note %d", time.Now().Unix()),
			Content:  "Enter your note content here",
			Category: "general",
			Tags:     []string{"new"},
		}
		err := m.NotesManager.CreateNote(newNote)
		if err != nil {
			m.StatusMessage = fmt.Sprintf("Error creating note: %v", err)
		} else {
			m.LoadNotes()
			m.StatusMessage = "Note created successfully"
		}
		return m, nil
	case "e":
		if m.NoteCursor < len(m.NotesList) {
			note := m.NotesList[m.NoteCursor]
			updatedNote, err := m.OpenEditorForNote(note)
			if err != nil {
				m.StatusMessage = fmt.Sprintf("Error opening editor: %v", err)
			} else if updatedNote != nil {
				err := m.NotesManager.UpdateNote(note.ID, updatedNote)
				if err != nil {
					m.StatusMessage = fmt.Sprintf("Error updating note: %v", err)
				} else {
					m.LoadNotes()
					m.StatusMessage = fmt.Sprintf("Note '%s' updated", updatedNote.Title)
				}
			}
		}
		return m, nil
	case "d":
		if m.NoteCursor < len(m.NotesList) {
			noteID := m.NotesList[m.NoteCursor].ID
			m.NotesManager.DeleteNote(noteID)
			m.LoadNotes()
			m.StatusMessage = "Note deleted"
		}
		return m, nil
	case "f":
		if m.NoteCursor < len(m.NotesList) {
			noteID := m.NotesList[m.NoteCursor].ID
			m.NotesManager.ToggleFavorite(noteID)
			m.LoadNotes()
		}
		return m, nil
	}
	return m, nil
}
