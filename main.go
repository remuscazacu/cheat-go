package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"cheat-go/pkg/apps"
	"cheat-go/pkg/cache"
	"cheat-go/pkg/config"
	"cheat-go/pkg/notes"
	"cheat-go/pkg/online"
	"cheat-go/pkg/plugins"
	"cheat-go/pkg/ui"
)

const (
	version = "v1.0.0-phase4"
	appName = "cheat-go"
)

type cliOptions struct {
	showHelp    bool
	showVersion bool
	theme       string
	tableStyle  string
	configFile  string
}

func printHelp() {
	fmt.Printf(`%s %s - Interactive terminal cheat sheet viewer with cloud features

USAGE:
    %s [OPTIONS]

DESCRIPTION:
    A fast, interactive terminal application for displaying keyboard shortcuts 
    and command cheat sheets. Navigate through shortcuts for popular applications
    with plugin support, personal notes, online repositories, and cloud sync.

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
    /                       Search mode
    f                       Filter apps
    n                       Open notes manager
    p                       Open plugin manager
    o                       Browse online repositories
    s                       Show sync status
    Ctrl+S                  Force sync
    ?                       Show help
    q / Ctrl+C              Quit the application

PHASE 4 FEATURES:
    Notes Manager (n)       Create and manage personal notes
    Plugin Manager (p)      Load and manage plugins
    Online Browser (o)      Browse community cheat sheets
    Sync Status (s)         View and manage cloud sync
    
THEMES:
    default                 Balanced colors for general use
    dark                    High-contrast for dark terminals
    light                   Clean appearance for light terminals  
    minimal                 Reduced visual elements

TABLE STYLES:
    simple                  Clean borders with lines
    rounded                 Elegant rounded corners
    bold                    Thick borders for visibility
    minimal                 Spacing-based separation

EXAMPLES:
    %s                      # Start with default settings
    %s --theme dark         # Use dark theme
    %s --style rounded      # Use rounded table borders
    %s -t dark -s bold      # Dark theme with bold borders
    %s --config my.yaml     # Use custom config file

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
	flag.StringVar(&opts.theme, "t", "", "Theme")
	flag.StringVar(&opts.theme, "theme", "", "Theme")
	flag.StringVar(&opts.tableStyle, "s", "", "Table style")
	flag.StringVar(&opts.tableStyle, "style", "", "Table style")
	flag.StringVar(&opts.configFile, "c", "", "Configuration file path")
	flag.StringVar(&opts.configFile, "config", "", "Configuration file path")

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

func initialModel(opts cliOptions) ui.Model {
	// Load configuration
	loader := config.NewLoader(opts.configFile)
	cfg, err := loader.Load()
	if err != nil {
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
		fmt.Printf("Warning: Could not load some apps (%v), using defaults\n", err)
	}

	// Create theme and renderer
	theme := ui.GetTheme(cfg.Theme)
	renderer := ui.NewTableRenderer(theme)
	renderer.SetTableStyle(cfg.Layout.TableStyle)
	renderer.SetMaxWidth(cfg.Layout.MaxWidth)

	// Generate table data
	rows := registry.GetTableData(cfg.Apps)

	// Initialize Phase 4 components
	m := ui.Model{
		Registry:     registry,
		Config:       cfg,
		Renderer:     renderer,
		Rows:         rows,
		FilteredRows: rows,
		CursorX:      0,
		CursorY:      1,
		SearchMode:   false,
		SearchQuery:  "",
		LastSearch:   "",
		AllRows:      rows,
		FilterMode:   false,
		FilteredApps: make([]string, 0),
		AllApps:      cfg.Apps,
		HelpMode:     false,
		ViewMode:     ui.ViewMain,
	}

	// Initialize cache
	m.Cache = cache.NewLRUCache(10*1024*1024, 1000) // 10MB, 1000 items

	// Initialize notes manager
	notesDir := os.ExpandEnv("$HOME/.config/cheat-go/notes")
	if cfg.DataDir != "" {
		notesDir = cfg.DataDir + "/notes"
	}
	m.NotesManager, _ = notes.NewFileManager(notesDir)

	// Initialize plugin loader
	pluginDirs := []string{
		os.ExpandEnv("$HOME/.config/cheat-go/plugins"),
		"/usr/local/share/cheat-go/plugins",
	}
	if cfg.DataDir != "" {
		pluginDirs = append([]string{cfg.DataDir + "/plugins"}, pluginDirs...)
	}
	m.PluginLoader = plugins.NewLoader(pluginDirs...)
	m.PluginLoader.LoadAll()

	// Initialize online client (mock for now)
	m.OnlineClient = online.NewMockClient()

	// Initialize sync manager (disabled by default)
	// m.syncManager would be initialized if sync is enabled in config

	return m
}

func main() {
	opts := parseFlags()

	if opts.showHelp {
		printHelp()
		os.Exit(0)
	}

	if opts.showVersion {
		printVersion()
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel(opts))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
