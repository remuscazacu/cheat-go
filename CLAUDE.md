# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**cheat-go** is a Terminal User Interface (TUI) application built with Go that provides an interactive way to browse and reference keyboard shortcuts for popular applications. It uses the Bubble Tea framework for the TUI and includes advanced features like notes management, plugin system, online repository browsing, and cloud sync.

## Build & Test Commands

```bash
# Build
go build -o cheat-go .

# Run
go run .

# Test all packages (Note: main package tests currently have issues)
go test ./pkg/...

# Test with coverage
go test ./pkg/... -cover

# Test with race detection
go test ./pkg/... -race

# Format code
go fmt ./...

# Static analysis
go vet ./pkg/...

# Generate coverage report
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Architecture

The codebase follows a modular package structure under `pkg/`:

- **pkg/apps/** - Application registry and shortcut data management (93.6% coverage)
  - Manages application definitions and shortcuts
  - Provides hardcoded defaults and YAML loading

- **pkg/config/** - Configuration loading and validation (91.3% coverage)
  - Handles YAML config files from multiple locations
  - Manages themes, layouts, and keybindings

- **pkg/notes/** - Personal notes system (90.7% coverage)
  - CRUD operations for notes with categories and tags
  - External editor integration via $EDITOR

- **pkg/ui/** - UI components with table rendering and themes (90.2% coverage)
  - Multiple view types: main, help, notes, plugins, online, sync
  - Theme support: default, dark, light, minimal

- **pkg/online/** - Community cheat sheet repository integration (85.9% coverage)
  - HTTP client for downloading sheets
  - Repository browsing functionality

- **pkg/plugins/** - Plugin system for extensions
  - Dynamic loading of external plugins
  - Built-in plugins for common features

- **pkg/sync/** - Cloud synchronization (43.7% coverage)
  - Conflict resolution and auto-sync
  - Support for multiple sync providers

- **pkg/cache/** - Performance caching with LRU (46.1% coverage)
  - Multi-level caching for performance
  - Memory and disk cache support

**main.go** - Application entry point and TUI orchestration using Bubble Tea

## Key Dependencies

- **github.com/charmbracelet/bubbletea** - TUI framework
- **github.com/charmbracelet/lipgloss** - Terminal styling
- **github.com/mattn/go-runewidth** - Unicode width calculations
- **gopkg.in/yaml.v3** - YAML parsing

## Code Style

- Go 1.22+ with standard formatting (`go fmt`)
- Comprehensive error handling with descriptive messages
- Table-driven tests for multiple scenarios
- Package-level documentation for exported types

## Known Issues

- Main package tests have undefined references and need fixing
- Some formatted files may need to be committed (pkg/cache/cache_test.go, pkg/sync/sync_test.go)

## Testing Strategy

The project maintains high test coverage (70.2% overall) with focus on:
- Core packages maintaining 85%+ coverage
- Edge case testing for error conditions
- External editor integration testing for notes
- Mock HTTP servers for online features