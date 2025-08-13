package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the visual styling for the application
type Theme struct {
	HeaderStyle lipgloss.Style
	CellStyle   lipgloss.Style
	BorderColor lipgloss.Color
}

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return &Theme{
		HeaderStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		CellStyle:   lipgloss.NewStyle().Padding(0, 0),
		BorderColor: lipgloss.Color("240"),
	}
}

// DarkTheme returns a dark theme variant
func DarkTheme() *Theme {
	return &Theme{
		HeaderStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")),
		CellStyle:   lipgloss.NewStyle().Padding(0, 0),
		BorderColor: lipgloss.Color("238"),
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "dark":
		return DarkTheme()
	default:
		return DefaultTheme()
	}
}
