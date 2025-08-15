# cheat-go

> A fast, interactive terminal application for displaying keyboard shortcuts and command cheat sheets

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue)](https://golang.org/dl/)
[![Test Coverage](https://img.shields.io/badge/Coverage-70.2%25-green)](./TEST_SUMMARY.md)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## Overview

**cheat-go** is a modern Terminal User Interface (TUI) application built with Go that provides an interactive way to browse and reference keyboard shortcuts for popular applications. It features a clean, navigable table interface that displays shortcuts across multiple applications simultaneously, making it easy to compare and learn commands.

### ✨ Features

- 🚀 **Fast & Lightweight** - Built with Go for minimal resource usage
- 🎨 **Beautiful TUI** - Modern terminal interface with multiple themes and table styles
- ⌨️ **Vim-style Navigation** - Support for both arrow keys and hjkl movement
- 🔍 **Interactive Search** - Real-time search through shortcuts and descriptions with highlighting
- 🎯 **App Filtering** - Select specific applications to focus on with visual interface
- ⌨️ **Rich Keyboard Shortcuts** - Comprehensive hotkey system with built-in help
- 📦 **Configurable** - YAML-based configuration with sensible defaults
- 🔧 **Extensible** - Plugin system for custom extensions and integrations
- 🌐 **Unicode Support** - Full support for international characters and emojis
- 📊 **Multiple Views** - Tabular display across multiple applications
- 🎭 **Multiple Themes** - Default, dark, light, and minimal theme options
- 📋 **Table Styles** - Simple, rounded, bold, and minimal table styles

#### 🆕 Phase 4 Features (Now with TUI Integration!)

- 🔌 **Plugin Architecture** - Extensible plugin system for custom functionality (press `p`)
- 🌍 **Community Integration** - Download and share cheat sheets from online repositories (press `o`)
- 📝 **Personal Notes** - Create and manage personal notes with shortcuts (press `n`)
- ☁️ **Cloud Sync** - Synchronize your data across multiple devices (press `s`)
- ⚡ **Performance Cache** - Multi-level caching for optimal performance (automatic)
- 🔄 **Auto-sync** - Automatic synchronization with conflict resolution (press `Ctrl+S` to force)

### 🎯 Supported Applications

Out of the box, cheat-go includes shortcuts for:

- **vim** - Vi IMproved text editor
- **zsh** - Z Shell
- **dwm** - Dynamic window manager  
- **st** - Simple terminal
- **lf** - Terminal file manager
- **zathura** - Document viewer

## 📸 Screenshots

### Main Interface
```
 Shortcut │ vim     │ zsh            │ dwm         │ st       │ lf      │ zathura  
──────────┼─────────┼────────────────┼─────────────┼──────────┼─────────┼──────────
 k        │ ↑ move  │ up history     │ focus up    │ ↑ scroll │ up      │ scroll ↑ 
 j        │ ↓ move  │ down history   │ focus down  │ ↓ scroll │ down    │ scroll ↓ 
 h        │ ← move  │ back char      │ focus left  │ ← move   │ left    │ scroll ← 
 l        │ → move  │ forward char   │ focus right │ → move   │ right   │ scroll → 
 /        │ search  │ search history │ -           │ search   │ search  │ search   
 q        │ quit    │ exit           │ close win   │ exit     │ quit    │ quit     

Arrow keys/hjkl: move • /: search • f: filter • Ctrl+R: refresh • ?: help • q: quit
```

### Search Mode
```
 Shortcut │ vim     │ zsh            │ dwm         │ st       │ lf      │ zathura  
──────────┼─────────┼────────────────┼─────────────┼──────────┼─────────┼──────────
 k        │ ↑ MOVE  │ up history     │ focus up    │ ↑ scroll │ up      │ scroll ↑ 
 j        │ ↓ MOVE  │ down history   │ focus down  │ ↓ scroll │ down    │ scroll ↓ 

Search: move_
Type to search, Enter to confirm, Esc to cancel
```

### App Filter Mode
```
Filter Apps:  [1] ✓vim [2] ✓zsh [3] dwm [4] st [5] lf [6] zathura
1-9: toggle apps, a: all, c: clear, Enter: apply, Esc: cancel
```

## 🚀 Installation

### Prerequisites

- **Go 1.22+** - [Download and install Go](https://golang.org/dl/)
- **Terminal** with Unicode support (most modern terminals)
- **Text Editor** (optional) - For notes editing functionality
  - Set `$EDITOR` environment variable to your preferred editor (vim, nano, emacs, code, etc.)
  - If not set, defaults to `nano`

### Method 1: Install from Source

```bash
# Clone the repository
git clone https://github.com/remuscazacu/cheat-go.git
cd cheat-go

# Build and install
go build -o cheat-go .
sudo mv cheat-go /usr/local/bin/

# Or install directly with go install
go install .
```

### Method 2: Download Binary

```bash
# Download the latest release (replace with actual release URL)
curl -L -o cheat-go https://github.com/remuscazacu/cheat-go/releases/latest/download/cheat-go
chmod +x cheat-go
sudo mv cheat-go /usr/local/bin/
```

### Method 3: Build from Source

```bash
git clone https://github.com/remuscazacu/cheat-go.git
cd cheat-go
go build .
./cheat-go
```

## 🎮 Usage

### Basic Usage

Simply run the application:

```bash
cheat-go
```

### Using Phase 4 Features

#### Interactive TUI Features

Simply use the keyboard shortcuts in the main application:

- Press `n` to open the **Notes Manager** - Create and manage personal notes
- Press `p` to open the **Plugin Manager** - View and manage installed plugins  
- Press `o` to browse **Online Repositories** - Discover and download community cheat sheets
- Press `s` to view **Sync Status** - Monitor cloud synchronization
- Press `Ctrl+S` to **Force Sync** - Manually trigger synchronization

Each view has its own set of keyboard shortcuts displayed at the bottom of the screen.

#### 📝 Notes Manager Features

The Notes Manager provides a full-featured personal notes system with the following capabilities:

**Note Creation**:
- Press `n` to create a new note with auto-generated timestamp title
- Notes include title, content, category, and tags
- Default category is "general" with "new" tag

**Note Editing**:
- Press `e` to edit notes in your default editor (✅ **Fixed in Phase 4**)
- Uses `$EDITOR` environment variable (falls back to `nano`)
- Opens with structured content format for easy editing:
  ```
  # Title: Your Note Title
  # Category: general
  # Tags: tag1, tag2
  
  Your note content here...
  ```
- Edit any field including title, category, tags, and content
- Changes are saved automatically when editor exits
- Works with vim, nano, emacs, VS Code, Sublime Text, and any terminal editor
- Supports both terminal and GUI editors that accept file arguments

**Note Management**:
- Press `d` to delete selected note
- Press `f` to toggle favorite status
- Navigate with arrow keys or `j/k` (vim-style)
- Visual indicators show favorites and categories

#### Notes Manager View (n)
- `n` - Create new note with timestamp-based title
- `e` - **Edit selected note in default editor** (✅ Fixed: opens $EDITOR or nano)
- `d` - Delete selected note 
- `f` - Toggle favorite status
- `up/down, j/k` - Navigate notes list
- `esc/q` - Return to main view

**Recent Fix**: The edit functionality now properly opens your default editor instead of just appending text. This provides a full editing experience with syntax highlighting, vim/emacs bindings, and your preferred editor features.

#### Plugin Manager View (p)
- `l` - Load selected plugin
- `u` - Unload selected plugin  
- `r` - Reload all plugins
- `up/down, j/k` - Navigate plugins list
- `esc/q` - Return to main view

#### Online Browser View (o)
- `enter` - Browse repository or download sheet
- `d` - Download selected cheat sheet
- `/` - Search online repositories
- `up/down, j/k` - Navigate repositories list
- `esc/q` - Return to main view

#### Sync Status View (s)
- `s` - Trigger sync now
- `r` - Resolve pending conflicts
- `a` - Toggle auto-sync enabled/disabled
- `up/down, j/k` - Navigate sync items
- `esc/q` - Return to main view

### Search Functionality

cheat-go includes powerful search capabilities to help you find shortcuts quickly:

- **Press `/`** to enter search mode
- **Type your query** to search through shortcut keys, descriptions, and categories
- **Matched terms are highlighted** in the results for easy identification
- **Press Enter** to confirm search and exit search mode
- **Press Esc** to cancel search and return to full table

### App Filtering

Focus on specific applications by filtering the displayed columns:

- **Press `f`** to enter filter mode
- **Use number keys (1-9)** to toggle individual applications
- **Press `a`** to select all applications
- **Press `c`** to clear all selections
- **Press Enter** to apply the filter
- **Press Esc** to cancel and return to previous state

### Interactive Help

- **Press `?`** at any time to see the comprehensive help screen
- The help screen shows all available keyboard shortcuts organized by category
- **Press `?` or `Esc`** to close the help screen

### Navigation

- **Arrow Keys** or **hjkl** - Navigate through the table
- **q** or **Ctrl+C** - Quit the application
- **/** - Enter search mode
- **f** - Enter filter mode
- **?** - Show help screen

### Keyboard Shortcuts

| Category | Key | Action |
|----------|-----|--------|
| **Navigation** | `↑` / `k` | Move cursor up |
| | `↓` / `j` | Move cursor down |
| | `←` / `h` | Move cursor left |
| | `→` / `l` | Move cursor right |
| | `Home` / `Ctrl+A` | Go to first row |
| | `End` / `Ctrl+E` | Go to last row |
| **Search** | `/` | Enter search mode |
| | `Enter` | Confirm search |
| | `Esc` | Exit search / clear filters |
| | `Backspace` | Delete character |
| | `Ctrl+U` | Clear search query |
| **Filtering** | `f` / `Ctrl+F` | Enter filter mode |
| | `1-9` | Toggle app selection |
| | `a` | Select all apps |
| | `c` | Clear all selections |
| | `Enter` | Apply filter |
| | `Esc` | Cancel filter |
| **Phase 4 Features** | `n` | Open notes manager |
| | `p` | Plugin manager |
| | `s` | Sync status |
| | `o` | Browse online repos |
| | `Ctrl+S` | Force sync |
| **General** | `Ctrl+R` | Refresh data |
| | `?` | Show/hide help |
| | `q` / `Ctrl+C` | Quit application |

## ⚙️ Configuration

cheat-go supports configuration through YAML files. The application looks for configuration files in the following order:

1. `~/.config/cheat-go/config.yaml`
2. `~/.cheat-go.yaml`
3. `./config.yaml` (current directory)

### Configuration File Example

```yaml
# ~/.config/cheat-go/config.yaml
apps:
  - vim
  - zsh
  - dwm
  - st
  - lf
  - zathura

theme: default  # or "dark"

layout:
  columns:
    - shortcut
    - description
  show_categories: false
  table_style: simple
  max_width: 120

keybinds:
  quit: q
  up: k
  down: j
  left: h
  right: l
  search: /
  next_app: tab
  prev_app: shift+tab

data_dir: ~/.config/cheat-go/apps

# New Phase 4 configuration options
plugins:
  enabled: true
  dirs:
    - ~/.config/cheat-go/plugins
    - /usr/local/share/cheat-go/plugins

sync:
  enabled: true
  service: cloud  # or "local"
  endpoint: https://sync.cheatsheets.com
  auto_sync: true
  interval: 15m

cache:
  enabled: true
  memory_size: 10485760  # 10MB
  disk_cache: ~/.cache/cheat-go

community:
  repositories:
    - https://github.com/cheat-go/community
    - https://github.com/awesome/cheatsheets
  auto_update: true

# Editor configuration for notes
editor:
  default: nano  # fallback if $EDITOR not set
  temp_prefix: cheat-go-note-
  backup: true
```

### Editor Setup for Notes

To use the notes editing feature effectively, configure your preferred editor:

```bash
# Set your preferred editor (choose one)
export EDITOR=vim      # Vim
export EDITOR=nano     # Nano (default)
export EDITOR=emacs    # Emacs
export EDITOR=code     # VS Code
export EDITOR=subl     # Sublime Text

# Add to your shell profile for persistence
echo 'export EDITOR=vim' >> ~/.bashrc    # Bash
echo 'export EDITOR=vim' >> ~/.zshrc     # Zsh
```

### Custom Applications

You can add custom applications by creating YAML files in your data directory:

```yaml
# ~/.config/cheat-go/apps/tmux.yaml
name: tmux
description: Terminal multiplexer
version: "1.0"
categories:
  - terminal
  - multiplexer
shortcuts:
  - keys: "Ctrl+b"
    description: "prefix key"
    category: "general"
  - keys: "Ctrl+b c"
    description: "new window"
    category: "window"
  - keys: "Ctrl+b d"
    description: "detach session"
    category: "session"
```

## 🏗️ Architecture

cheat-go is built with a clean, modular architecture:

```
cheat-go/
├── main.go              # Application entry point
├── pkg/
│   ├── apps/           # Application registry and data management
│   │   ├── types.go    # Data structures
│   │   └── registry.go # App loading and management
│   ├── config/         # Configuration management
│   │   ├── types.go    # Config structures
│   │   └── loader.go   # Config loading and validation
│   ├── ui/             # User interface components
│   │   ├── table.go    # Table rendering
│   │   └── theme.go    # Themes and styling
│   ├── plugins/        # Plugin system (NEW)
│   │   ├── types.go    # Plugin interfaces
│   │   ├── loader.go   # Plugin loading
│   │   └── builtin.go  # Built-in plugins
│   ├── online/         # Community integration (NEW)
│   │   ├── types.go    # Online data structures
│   │   └── client.go   # HTTP client for repositories
│   ├── notes/          # Personal notes (NEW)
│   │   ├── types.go    # Note structures
│   │   └── manager.go  # Note management
│   ├── sync/           # Cloud sync (NEW)
│   │   └── sync.go     # Synchronization logic
│   └── cache/          # Performance cache (NEW)
│       └── cache.go    # Multi-level caching
└── *_test.go           # Comprehensive test suite
```

### Key Components

- **Apps Package** - Manages application definitions and shortcut data
- **Config Package** - Handles configuration loading and validation
- **UI Package** - Provides table rendering and theming capabilities
- **Plugins Package** - Extensible plugin system for custom functionality
- **Online Package** - Community cheat sheet repository integration
- **Notes Package** - Personal notes and custom shortcuts management
- **Sync Package** - Cross-platform cloud synchronization
- **Cache Package** - Multi-level caching for performance
- **Main** - Coordinates the TUI application using Bubble Tea

## 🧪 Development

### Prerequisites for Development

- Go 1.22+
- Git

### Building

```bash
# Clone and enter directory
git clone https://github.com/remuscazacu/cheat-go.git
cd cheat-go

# Install dependencies
go mod download

# Build
go build .

# Run
./cheat-go
```

### Testing

The project includes a comprehensive test suite with **70.2% overall coverage**:

#### Package Coverage Details:
- **pkg/notes**: 90.7% - Full note management functionality
- **pkg/config**: 91.3% - Configuration loading and validation  
- **pkg/apps**: 93.6% - Application registry and data management
- **pkg/online**: 85.9% - Community integration features
- **pkg/ui**: 90.2% - Table rendering and themes
- **main package**: 60.4% - TUI application logic
- **pkg/cache**: 46.1% - Performance caching
- **pkg/sync**: 43.7% - Cloud synchronization

#### Test Commands:
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./pkg/notes -v
go test ./pkg/config -v

# Run tests with race detection
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Test notes manager editor functionality (new)
go test -v -run TestOpenEditorForNote

# Test all notes manager features
go test ./pkg/notes -v
```

#### Recent Test Improvements:
- **✅ Added editor functionality tests** for notes manager (`TestOpenEditorForNote`)
- **✅ Enhanced error handling tests** across all packages
- **✅ Improved validation tests** for config and apps packages  
- **✅ Added edge case coverage** for file operations
- **✅ Comprehensive integration tests** for main application
- **✅ Fixed notes manager editor bug** with proper external editor integration
- **✅ Achieved 90%+ coverage** in core packages (notes, config, apps, online, ui)

### Code Quality

```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## 📁 Project Structure

```
cheat-go/
├── README.md                    # This file
├── main.go                      # Application entry point and TUI
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums  
├── TEST_SUMMARY.md              # Test coverage documentation
├── AGENTS.md                    # Development guidelines
├── docs/                        # Documentation
│   ├── ARCHITECTURE_ANALYSIS.md # Architecture overview
│   ├── PHASE_4_IMPLEMENTATION.md # Phase 4 feature docs
│   └── PHASE_4_TUI_INTEGRATION.md # TUI integration guide
├── examples/                    # Example configurations
│   ├── config.yaml             # Sample config file
│   └── apps/                   # Sample app definitions
│       ├── vim.yaml
│       ├── zsh.yaml
│       └── ...
├── pkg/
│   ├── apps/                   # Application registry (93.6% coverage)
│   │   ├── types.go           # App and shortcut structures
│   │   ├── registry.go        # App loading and management
│   │   ├── types_test.go      # Type validation tests
│   │   ├── registry_test.go   # Registry functionality tests
│   │   └── registry_edge_test.go # Edge case coverage
│   ├── config/                 # Configuration system (91.3% coverage)
│   │   ├── types.go           # Config structures
│   │   ├── loader.go          # Config loading and validation
│   │   ├── types_test.go      # Config structure tests
│   │   ├── loader_test.go     # Loader functionality tests
│   │   └── loader_edge_test.go # Error handling tests
│   ├── notes/                  # Personal notes system (90.7% coverage)
│   │   ├── types.go           # Note data structures
│   │   ├── manager.go         # Note CRUD and management
│   │   └── manager_test.go    # Comprehensive note tests
│   ├── online/                 # Community integration (85.9% coverage)
│   │   ├── types.go           # Online repository structures
│   │   ├── client.go          # HTTP client for repositories
│   │   └── client_test.go     # Client and mock tests
│   ├── plugins/                # Plugin system
│   │   ├── types.go           # Plugin interfaces
│   │   ├── loader.go          # Plugin loading logic
│   │   ├── builtin.go         # Built-in plugins
│   │   └── loader_test.go     # Plugin system tests
│   ├── sync/                   # Cloud synchronization (43.7% coverage)
│   │   ├── sync.go            # Sync logic and conflict resolution
│   │   └── sync_test.go       # Sync functionality tests
│   ├── cache/                  # Performance caching (46.1% coverage)
│   │   ├── cache.go           # LRU cache implementation
│   │   └── cache_test.go      # Cache functionality tests
│   └── ui/                     # User interface (90.2% coverage)
│       ├── table.go           # Table rendering engine
│       ├── theme.go           # Theme and styling
│       ├── table_test.go      # Table rendering tests
│       └── theme_test.go      # Theme system tests
├── main_test.go                # Integration tests (60.4% coverage)
├── main_edge_test.go          # Additional edge case tests
└── coverage.html              # Test coverage report
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Guidelines

- Follow standard Go conventions
- Write tests for new functionality
- Maintain or improve test coverage
- Use meaningful commit messages
- Update documentation as needed

## 📋 Roadmap

### Completed Features ✅

- [x] Search functionality within shortcuts
- [x] Interactive search with highlighting
- [x] App filtering and selection
- [x] Comprehensive keyboard shortcuts
- [x] Multiple theme support (default, dark, light, minimal)
- [x] Multiple table styles (simple, rounded, bold, minimal)
- [x] Built-in help system
- [x] Plugin system for external applications (Phase 4)
- [x] Community repository integration (Phase 4)
- [x] Personal notes system with editor integration (Phase 4)
- [x] Cloud synchronization with conflict resolution (Phase 4) 
- [x] Performance caching with LRU eviction (Phase 4)
- [x] Export functionality (JSON, Markdown via plugins)
- [x] Notes editor with $EDITOR support and structured editing (✅ **Fixed editor bug**)
- [x] Comprehensive test suite with 70.2% coverage and 90%+ in core packages
- [x] Editor integration tests for notes management functionality

### Planned Features

- [ ] Custom key bindings configuration
- [ ] Interactive tutorial mode
- [ ] Fuzzy search across applications
- [ ] Shortcut categories and advanced filtering
- [ ] Application profiles (work, gaming, etc.)
- [ ] Search history and saved searches
- [ ] Web interface for remote access
- [ ] Mobile companion apps
- [ ] AI-powered shortcut suggestions
- [ ] Team collaboration features

### Version History

- **v0.1.0** - Initial release with basic TUI functionality
- **v0.2.0** - Added configuration system and custom apps
- **v0.3.0** - Comprehensive test suite and improved error handling
- **v0.4.0** - Enhanced UI with multiple table styles and improved themes
- **v0.5.0** - Interactive search, app filtering, and comprehensive keyboard shortcuts
- **v1.0.0-phase4** - Phase 4: Full TUI integration with plugins, notes (editor support), online repos, sync, and caching
  - ✅ **Fixed editor bug**: Notes editor now properly opens external editor ($EDITOR or nano)
  - ✅ **Added comprehensive editor tests**: Full test coverage for editor functionality 
  - ✅ **Enhanced keyboard shortcuts**: Complete Phase 4 TUI integration with intuitive navigation
  - ✅ **Achieved 70.2% overall test coverage** with 90%+ coverage in core packages:
    - pkg/notes: 90.7% (includes new editor tests)
    - pkg/config: 91.3% 
    - pkg/apps: 93.6%
    - pkg/online: 85.9%
    - pkg/ui: 90.2%
  - ✅ **Structured note editing**: Title, category, tags, and content in user's preferred editor
  - ✅ **Improved error handling**: Better user feedback and graceful failure handling

## 🐛 Troubleshooting

### Common Issues

**Q: Application doesn't start**
- Ensure Go 1.22+ is installed
- Check terminal Unicode support
- Verify binary permissions

**Q: Configuration not loading**
- Check YAML syntax in config files
- Verify file permissions
- Use absolute paths for custom data directories

**Q: Display issues in terminal**
- Ensure terminal supports Unicode
- Try different themes: `default`, `dark`, `light`, or `minimal`
- Try different table styles: `simple`, `rounded`, `bold`, or `minimal`
- Check terminal size (minimum 80x24 recommended)
- Use `cheat-go --help` to see all available options

**Q: Search not finding results**
- Search is case-insensitive and searches keys, descriptions, and categories
- Try shorter search terms or partial matches
- Use Esc to clear filters and return to full table
- Press ? for help with search keyboard shortcuts

**Q: Notes editor not working**
- ✅ **Recently Fixed**: Editor now properly opens external editor (Phase 4 update)
- Ensure `$EDITOR` environment variable is set: `echo $EDITOR`
- If not set, it defaults to `nano` - ensure nano is installed: `which nano`
- Test your editor manually: `$EDITOR test.txt` or `nano test.txt`
- **Supported editors**: vim, nano, emacs, code, subl, micro, helix, and most terminal editors
- **For GUI editors**: Ensure they support command-line file opening (e.g., `code --wait`)
- **Example setup**: `export EDITOR="vim"` or `export EDITOR="code --wait"`
- **Troubleshooting**: 
  - Run `go test -v -run TestOpenEditorForNote` to test editor integration
  - Check file permissions in `~/.config/cheat-go/` directory
  - Verify editor accepts file arguments: `your-editor --help`

**Q: Phase 4 features not working**
- Ensure you have proper file permissions in `~/.config/cheat-go/`
- Check internet connection for online repository features
- Verify plugin directories exist and are readable
- Use `Ctrl+S` to force sync if cloud sync appears stuck

### Reporting Issues

Please report issues on [GitHub Issues](https://github.com/remuscazacu/cheat-go/issues) with:

1. Operating system and version
2. Go version (`go version`)
3. Terminal emulator
4. Steps to reproduce
5. Expected vs actual behavior

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the excellent Bubble Tea TUI framework and Lipgloss styling
- [mattn/go-runewidth](https://github.com/mattn/go-runewidth) for Unicode width calculations
- Go community for the robust standard library and tooling
- Contributors and users who provide feedback and improvements

## 📞 Support

- 📖 **Documentation**: Check this README and inline code comments
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/remuscazacu/cheat-go/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/remuscazacu/cheat-go/discussions)
- 📧 **Contact**: [Project maintainer](mailto:your-email@example.com)

---

<div align="center">

**⭐ Star this repository if you find it helpful!**

Made with ❤️ and Go

</div>