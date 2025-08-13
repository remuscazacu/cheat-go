package ui

import (
	"github.com/charmbracelet/lipgloss"
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()

	if theme == nil {
		t.Fatal("DefaultTheme() returned nil")
	}

	// Test that styles are functional
	headerText := theme.HeaderStyle.Render("test")
	if headerText == "" {
		t.Error("HeaderStyle should be functional")
	}

	cellText := theme.CellStyle.Render("test")
	if cellText == "" {
		t.Error("CellStyle should be functional")
	}

	if theme.BorderColor == "" {
		t.Error("BorderColor should be set")
	}

	// Test specific default values
	if theme.BorderColor != lipgloss.Color("240") {
		t.Errorf("BorderColor = %v, expected lipgloss.Color(\"240\")", theme.BorderColor)
	}
}

func TestDarkTheme(t *testing.T) {
	theme := DarkTheme()

	if theme == nil {
		t.Fatal("DarkTheme() returned nil")
	}

	// Test that styles are functional
	headerText := theme.HeaderStyle.Render("test")
	if headerText == "" {
		t.Error("HeaderStyle should be functional")
	}

	cellText := theme.CellStyle.Render("test")
	if cellText == "" {
		t.Error("CellStyle should be functional")
	}

	if theme.BorderColor == "" {
		t.Error("BorderColor should be set")
	}

	// Test specific dark theme values
	if theme.BorderColor != lipgloss.Color("238") {
		t.Errorf("BorderColor = %v, expected lipgloss.Color(\"238\")", theme.BorderColor)
	}
}

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name     string
		expected *Theme
	}{
		{"dark", DarkTheme()},
		{"default", DefaultTheme()},
		{"unknown", DefaultTheme()},
		{"", DefaultTheme()},
	}

	for _, test := range tests {
		theme := GetTheme(test.name)

		if theme == nil {
			t.Errorf("GetTheme(%s) returned nil", test.name)
			continue
		}

		// Compare border colors as a way to verify correct theme
		if theme.BorderColor != test.expected.BorderColor {
			t.Errorf("GetTheme(%s) BorderColor = %v, expected %v",
				test.name, theme.BorderColor, test.expected.BorderColor)
		}
	}
}

func TestTheme_Structure(t *testing.T) {
	theme := &Theme{
		HeaderStyle: lipgloss.NewStyle().Bold(true),
		CellStyle:   lipgloss.NewStyle().Padding(1, 2),
		BorderColor: lipgloss.Color("123"),
	}

	// Test that the theme structure works as expected
	headerText := theme.HeaderStyle.Render("test")
	if headerText == "" {
		t.Error("HeaderStyle should be functional")
	}

	cellText := theme.CellStyle.Render("test")
	if cellText == "" {
		t.Error("CellStyle should be functional")
	}

	if theme.BorderColor != lipgloss.Color("123") {
		t.Error("BorderColor should be set correctly")
	}
}

func TestTheme_StyleApplication(t *testing.T) {
	theme := DefaultTheme()

	// Test that styles can be applied to text
	headerText := theme.HeaderStyle.Render("Header")
	if headerText == "" {
		t.Error("HeaderStyle should render text")
	}

	cellText := theme.CellStyle.Render("Cell")
	if cellText == "" {
		t.Error("CellStyle should render text")
	}

	// The actual styled text should contain the original text
	if len(headerText) < len("Header") {
		t.Error("styled header text should contain original text")
	}

	if len(cellText) < len("Cell") {
		t.Error("styled cell text should contain original text")
	}
}

func TestTheme_DifferentThemes(t *testing.T) {
	defaultTheme := DefaultTheme()
	darkTheme := DarkTheme()

	// Themes should be different
	if defaultTheme.BorderColor == darkTheme.BorderColor {
		t.Error("default and dark themes should have different border colors")
	}

	// Both should be valid themes
	if defaultTheme == nil || darkTheme == nil {
		t.Error("both themes should be valid")
	}
}

func TestGetTheme_CaseHandling(t *testing.T) {
	// Test case sensitivity - Go's switch is case-sensitive
	darkTheme := GetTheme("dark")
	upperDarkTheme := GetTheme("DARK") // Should return default

	expectedDark := DarkTheme()
	expectedDefault := DefaultTheme()

	if darkTheme.BorderColor != expectedDark.BorderColor {
		t.Error("lowercase 'dark' should return dark theme")
	}

	if upperDarkTheme.BorderColor != expectedDefault.BorderColor {
		t.Error("uppercase 'DARK' should return default theme (case sensitive)")
	}
}
