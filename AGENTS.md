# AGENTS.md - Development Guidelines

## Build & Test Commands
- **Build**: `go build .` or `go build -o cheat-go .`
- **Run**: `go run .` or `go run main.go`
- **Test**: `go test ./...` (no tests exist yet)
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...`
- **Clean**: `go clean`

## Code Style Guidelines
- **Language**: Go 1.24.2+
- **Formatting**: Use `go fmt` - tabs for indentation, standard Go formatting
- **Imports**: Group standard, third-party, local imports with blank lines between groups
- **Naming**: camelCase for unexported, PascalCase for exported, descriptive names
- **Variables**: Declare close to use, use short names for short scopes (`i`, `err`)
- **Types**: Define structs with consistent field ordering, embed interfaces when appropriate
- **Error Handling**: Always check errors, use early returns, descriptive error messages
- **Comments**: Package-level comments required, exported functions/types documented

## Dependencies
- Uses Charm libraries: bubbletea (TUI), lipgloss (styling)
- Minimal external dependencies, prefer standard library when possible
- No test framework configured yet - use standard `testing` package for new tests