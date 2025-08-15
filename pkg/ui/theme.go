package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the visual styling for the application
type Theme struct {
	Name             string
	HeaderStyle      lipgloss.Style
	CellStyle        lipgloss.Style
	HighlightStyle   lipgloss.Style
	BorderColor      lipgloss.Color
	SelectedRowStyle lipgloss.Style
	CategoryStyle    lipgloss.Style
	SearchStyle      lipgloss.Style
	SearchInputStyle lipgloss.Style
	TableStyle       string
}

// DefaultTheme returns the default theme
func DefaultTheme() *Theme {
	return &Theme{
		Name:             "default",
		HeaderStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		CellStyle:        lipgloss.NewStyle(),
		HighlightStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")),
		BorderColor:      lipgloss.Color("240"),
		SelectedRowStyle: lipgloss.NewStyle().Background(lipgloss.Color("238")),
		CategoryStyle:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")),
		SearchStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")),
		SearchInputStyle: lipgloss.NewStyle().Background(lipgloss.Color("235")),
		TableStyle:       "simple",
	}
}

// DarkTheme returns a dark theme variant
func DarkTheme() *Theme {
	return &Theme{
		Name:             "dark",
		HeaderStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")),
		CellStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		HighlightStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")),
		BorderColor:      lipgloss.Color("238"),
		SelectedRowStyle: lipgloss.NewStyle().Background(lipgloss.Color("236")),
		CategoryStyle:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82")),
		SearchStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")),
		SearchInputStyle: lipgloss.NewStyle().Background(lipgloss.Color("234")),
		TableStyle:       "rounded",
	}
}

// LightTheme returns a light theme variant
func LightTheme() *Theme {
	return &Theme{
		Name:             "light",
		HeaderStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("25")),
		CellStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("235")),
		HighlightStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")),
		BorderColor:      lipgloss.Color("244"),
		SelectedRowStyle: lipgloss.NewStyle().Background(lipgloss.Color("254")),
		CategoryStyle:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("28")),
		SearchStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")),
		SearchInputStyle: lipgloss.NewStyle().Background(lipgloss.Color("255")),
		TableStyle:       "simple",
	}
}

// MinimalTheme returns a minimal theme variant
func MinimalTheme() *Theme {
	return &Theme{
		Name:             "minimal",
		HeaderStyle:      lipgloss.NewStyle().Bold(true),
		CellStyle:        lipgloss.NewStyle(),
		HighlightStyle:   lipgloss.NewStyle().Bold(true),
		BorderColor:      lipgloss.Color("250"),
		SelectedRowStyle: lipgloss.NewStyle().Underline(true),
		CategoryStyle:    lipgloss.NewStyle().Bold(true),
		SearchStyle:      lipgloss.NewStyle().Bold(true),
		SearchInputStyle: lipgloss.NewStyle().Underline(true),
		TableStyle:       "minimal",
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "dark":
		return DarkTheme()
	case "light":
		return LightTheme()
	case "minimal":
		return MinimalTheme()
	default:
		return DefaultTheme()
	}
}

// GetAvailableThemes returns a list of all available theme names
func GetAvailableThemes() []string {
	return []string{"default", "dark", "light", "minimal"}
}
