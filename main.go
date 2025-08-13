package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/config"
	"cheat-go/pkg/ui"
)

const (
	version = "v0.4.0"
	appName = "cheat-go"
)

type model struct {
	registry *apps.Registry
	config   *config.Config
	renderer *ui.TableRenderer
	rows     [][]string
	cursorX  int
	cursorY  int
}

type cliOptions struct {
	showHelp    bool
	showVersion bool
	theme       string
	tableStyle  string
	configFile  string
}

func printHelp() {
	fmt.Printf(`%s %s - Interactive terminal cheat sheet viewer

USAGE:
    %s [OPTIONS]

DESCRIPTION:
    A fast, interactive terminal application for displaying keyboard shortcuts 
    and command cheat sheets. Navigate through shortcuts for popular applications
    like vim, zsh, dwm, st, lf, and zathura in a beautiful tabular interface.

OPTIONS:
    -h, --help              Show this help message and exit
    -v, --version           Show version information and exit
    -t, --theme THEME       Set the display theme
                            Options: default, dark, light, minimal
                            Default: default
    -s, --style STYLE       Set the table style  
                            Options: simple, rounded, bold, minimal
                            Default: simple
    -c, --config FILE       Use custom configuration file
                            Default: ~/.config/cheat-go/config.yaml

NAVIGATION:
    Arrow Keys / hjkl       Navigate through the table
    q / Ctrl+C             Quit the application

THEMES:
    default                Balanced colors for general use
    dark                   High-contrast for dark terminals
    light                  Clean appearance for light terminals  
    minimal                Reduced visual elements

TABLE STYLES:
    simple                 Clean borders with lines
    rounded                Elegant rounded corners
    bold                   Thick borders for visibility
    minimal                Spacing-based separation

EXAMPLES:
    %s                     # Start with default settings
    %s --theme dark        # Use dark theme
    %s --style rounded     # Use rounded table borders
    %s -t dark -s bold     # Dark theme with bold borders
    %s --config my.yaml    # Use custom config file

CONFIGURATION:
    Configuration files are loaded in this order:
    1. Command line --config option
    2. ~/.config/cheat-go/config.yaml
    3. ~/.cheat-go.yaml
    4. ./config.yaml

    See documentation for configuration file format.

SUPPORTED APPLICATIONS:
    vim, zsh, dwm, st, lf, zathura

For more information, visit: https://github.com/remuscazacu/cheat-go
`, appName, version, appName, appName, appName, appName, appName, appName)
}

func printVersion() {
	fmt.Printf("%s %s\n", appName, version)
}

func parseFlags() cliOptions {
	var opts cliOptions

	flag.BoolVar(&opts.showHelp, "h", false, "Show help message")
	flag.BoolVar(&opts.showHelp, "help", false, "Show help message")
	flag.BoolVar(&opts.showVersion, "v", false, "Show version")
	flag.BoolVar(&opts.showVersion, "version", false, "Show version")
	flag.StringVar(&opts.theme, "t", "", "Theme (default, dark, light, minimal)")
	flag.StringVar(&opts.theme, "theme", "", "Theme (default, dark, light, minimal)")
	flag.StringVar(&opts.tableStyle, "s", "", "Table style (simple, rounded, bold, minimal)")
	flag.StringVar(&opts.tableStyle, "style", "", "Table style (simple, rounded, bold, minimal)")
	flag.StringVar(&opts.configFile, "c", "", "Configuration file path")
	flag.StringVar(&opts.configFile, "config", "", "Configuration file path")

	// Custom usage function
	flag.Usage = func() {
		printHelp()
	}

	flag.Parse()

	// Validate theme option
	if opts.theme != "" {
		validThemes := []string{"default", "dark", "light", "minimal"}
		valid := false
		for _, t := range validThemes {
			if opts.theme == t {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Error: Invalid theme '%s'. Valid options: %s\n",
				opts.theme, strings.Join(validThemes, ", "))
			os.Exit(1)
		}
	}

	// Validate table style option
	if opts.tableStyle != "" {
		validStyles := []string{"simple", "rounded", "bold", "minimal"}
		valid := false
		for _, s := range validStyles {
			if opts.tableStyle == s {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("Error: Invalid table style '%s'. Valid options: %s\n",
				opts.tableStyle, strings.Join(validStyles, ", "))
			os.Exit(1)
		}
	}

	return opts
}

func initialModelWithDefaults() model {
	return initialModel(cliOptions{})
}

func initialModel(opts cliOptions) model {
	// Load configuration
	loader := config.NewLoader(opts.configFile)
	cfg, err := loader.Load()
	if err != nil {
		// Log error but continue with defaults
		fmt.Printf("Warning: Could not load config (%v), using defaults\n", err)
		cfg = config.DefaultConfig()
	}

	// Override config with CLI options
	if opts.theme != "" {
		cfg.Theme = opts.theme
	}
	if opts.tableStyle != "" {
		cfg.Layout.TableStyle = opts.tableStyle
	}

	// Initialize app registry
	registry := apps.NewRegistry(cfg.DataDir)
	if err := registry.LoadApps(cfg.Apps); err != nil {
		// Log warning but continue with hardcoded data
		fmt.Printf("Warning: Could not load some apps (%v), using defaults\n", err)
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
	// Parse command-line arguments
	opts := parseFlags()

	// Handle help and version flags
	if opts.showHelp {
		printHelp()
		return
	}

	if opts.showVersion {
		printVersion()
		return
	}

	// Start the TUI application
	p := tea.NewProgram(initialModel(opts))
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
