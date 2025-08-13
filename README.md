# cheat-go

> A fast, interactive terminal application for displaying keyboard shortcuts and command cheat sheets

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue)](https://golang.org/dl/)
[![Test Coverage](https://img.shields.io/badge/Coverage-94.3%25-brightgreen)](./TEST_SUMMARY.md)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## Overview

**cheat-go** is a modern Terminal User Interface (TUI) application built with Go that provides an interactive way to browse and reference keyboard shortcuts for popular applications. It features a clean, navigable table interface that displays shortcuts across multiple applications simultaneously, making it easy to compare and learn commands.

### ✨ Features

- 🚀 **Fast & Lightweight** - Built with Go for minimal resource usage
- 🎨 **Beautiful TUI** - Modern terminal interface with themes and styling
- ⌨️ **Vim-style Navigation** - Support for both arrow keys and hjkl movement
- 📦 **Configurable** - YAML-based configuration with sensible defaults
- 🔧 **Extensible** - Easy to add custom applications and shortcuts
- 🌐 **Unicode Support** - Full support for international characters and emojis
- 📊 **Multiple Views** - Tabular display across multiple applications

### 🎯 Supported Applications

Out of the box, cheat-go includes shortcuts for:

- **vim** - Vi IMproved text editor
- **zsh** - Z Shell
- **dwm** - Dynamic window manager  
- **st** - Simple terminal
- **lf** - Terminal file manager
- **zathura** - Document viewer

## 📸 Screenshots

```
 Shortcut │ vim     │ zsh            │ dwm         │ st       │ lf      │ zathura  
──────────┼─────────┼────────────────┼─────────────┼──────────┼─────────┼──────────
 k        │ ↑ move  │ up history     │ focus up    │ ↑ scroll │ up      │ scroll ↑ 
 j        │ ↓ move  │ down history   │ focus down  │ ↓ scroll │ down    │ scroll ↓ 
 h        │ ← move  │ back char      │ focus left  │ ← move   │ left    │ scroll ← 
 l        │ → move  │ forward char   │ focus right │ → move   │ right   │ scroll → 
 /        │ search  │ search history │ -           │ search   │ search  │ search   
 q        │ quit    │ exit           │ close win   │ exit     │ quit    │ quit     

Use arrow keys or hjkl to move. Press q to quit.
```

## 🚀 Installation

### Prerequisites

- **Go 1.22+** - [Download and install Go](https://golang.org/dl/)
- **Terminal** with Unicode support (most modern terminals)

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

### Navigation

- **Arrow Keys** or **hjkl** - Navigate through the table
- **q** or **Ctrl+C** - Quit the application

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `←` / `h` | Move cursor left |
| `→` / `l` | Move cursor right |
| `q` | Quit application |
| `Ctrl+C` | Force quit |

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
│   └── ui/             # User interface components
│       ├── table.go    # Table rendering
│       └── theme.go    # Themes and styling
└── *_test.go           # Comprehensive test suite
```

### Key Components

- **Apps Package** - Manages application definitions and shortcut data
- **Config Package** - Handles configuration loading and validation
- **UI Package** - Provides table rendering and theming capabilities
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

The project includes a comprehensive test suite with 94.3% coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with race detection
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

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
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── TEST_SUMMARY.md              # Test documentation
├── pkg/
│   ├── apps/
│   │   ├── types.go            # App and shortcut data structures
│   │   ├── registry.go         # App registry and loading
│   │   ├── types_test.go       # Unit tests for types
│   │   ├── registry_test.go    # Unit tests for registry
│   │   └── registry_edge_test.go # Edge case tests
│   ├── config/
│   │   ├── types.go            # Configuration structures
│   │   ├── loader.go           # Config loading logic
│   │   ├── types_test.go       # Config type tests
│   │   ├── loader_test.go      # Loader tests
│   │   └── loader_edge_test.go # Edge case tests
│   └── ui/
│       ├── table.go            # Table rendering
│       ├── theme.go            # Theme management
│       ├── table_test.go       # Table tests
│       └── theme_test.go       # Theme tests
├── main_test.go                 # Integration tests
└── main_edge_test.go           # Additional edge case tests
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

### Planned Features

- [ ] Search functionality within shortcuts
- [ ] Custom key bindings configuration
- [ ] Plugin system for external applications
- [ ] Export functionality (JSON, CSV, Markdown)
- [ ] Interactive tutorial mode
- [ ] Fuzzy search across applications
- [ ] Shortcut categories and filtering
- [ ] Multiple theme support
- [ ] Application profiles (work, gaming, etc.)

### Version History

- **v0.1.0** - Initial release with basic TUI functionality
- **v0.2.0** - Added configuration system and custom apps
- **v0.3.0** - Comprehensive test suite and improved error handling

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
- Try different themes (default vs dark)
- Check terminal size (minimum 80x24 recommended)

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