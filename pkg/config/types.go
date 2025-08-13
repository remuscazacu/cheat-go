package config

// Config represents the main application configuration
type Config struct {
	Apps     []string          `yaml:"apps" json:"apps"`
	Theme    string            `yaml:"theme" json:"theme"`
	Layout   LayoutConfig      `yaml:"layout" json:"layout"`
	Keybinds map[string]string `yaml:"keybinds" json:"keybinds"`
	DataDir  string            `yaml:"data_dir" json:"data_dir"`
}

// LayoutConfig controls the display layout
type LayoutConfig struct {
	Columns        []string `yaml:"columns" json:"columns"`
	ShowCategories bool     `yaml:"show_categories" json:"show_categories"`
	TableStyle     string   `yaml:"table_style" json:"table_style"`
	MaxWidth       int      `yaml:"max_width" json:"max_width"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Apps:  []string{"vim", "zsh", "dwm", "st", "lf", "zathura"},
		Theme: "default",
		Layout: LayoutConfig{
			Columns:        []string{"shortcut", "description"},
			ShowCategories: false,
			TableStyle:     "simple",
			MaxWidth:       120,
		},
		Keybinds: map[string]string{
			"quit":     "q",
			"up":       "k",
			"down":     "j",
			"left":     "h",
			"right":    "l",
			"search":   "/",
			"next_app": "tab",
			"prev_app": "shift+tab",
		},
		DataDir: "~/.config/cheat-go/apps",
	}
}
