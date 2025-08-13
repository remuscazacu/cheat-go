package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/config"
	"cheat-go/pkg/ui"
)

type model struct {
	registry *apps.Registry
	config   *config.Config
	renderer *ui.TableRenderer
	rows     [][]string
	cursorX  int
	cursorY  int
}

func initialModel() model {
	// Load configuration
	loader := config.NewLoader("")
	cfg, err := loader.Load()
	if err != nil {
		// Log error but continue with defaults
		fmt.Fprintf(os.Stderr, "Warning: Could not load config (%v), using defaults\n", err)
		cfg = config.DefaultConfig()
	}

	// Initialize app registry
	registry := apps.NewRegistry(cfg.DataDir)
	if err := registry.LoadApps(cfg.Apps); err != nil {
		// Log warning but continue with hardcoded data
		fmt.Fprintf(os.Stderr, "Warning: Could not load some apps (%v), using defaults\n", err)
	}

	// Create theme and renderer
	theme := ui.GetTheme(cfg.Theme)
	renderer := ui.NewTableRenderer(theme)
	renderer.SetTableStyle(cfg.Layout.TableStyle)
	renderer.SetMaxWidth(cfg.Layout.MaxWidth)

	// Generate table data
	rows := registry.GetTableData(cfg.Apps)

	return model{
		registry: registry,
		config:   cfg,
		renderer: renderer,
		rows:     rows,
		cursorX:  0,
		cursorY:  1,
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
	return m.renderer.RenderWithInstructions(m.rows, m.cursorX, m.cursorY)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
