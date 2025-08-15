# AGENTS.md - Development Guidelines

## Build & Test Commands
- **Build**: `go build .` or `go build -o cheat-go .`
- **Run**: `go run .` or `go run main.go`
- **Test**: `go test ./...` (✅ 70.2% coverage, 150+ tests)
- **Test Coverage**: `go test ./... -cover`
- **Race Detection**: `go test ./... -race`
- **Editor Tests**: `go test -v -run TestOpenEditorForNote`
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
- **Testing**: Comprehensive test suite using standard `testing` package
- **Phase 4 Features**: Notes manager with external editor integration ($EDITOR support)

## Testing Guidelines
- **Coverage Target**: Maintain 90%+ coverage in core packages (apps, config, notes, ui, online)
- **Test Types**: Unit tests, integration tests, edge case testing, race condition testing
- **Editor Testing**: Test external editor integration with mock editors and error scenarios
- **Naming**: Use descriptive test names like `TestOpenEditorForNote_InvalidEditor`
- **Structure**: Group related tests and use table-driven tests for multiple scenarios