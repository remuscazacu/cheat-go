package apps

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry("/tmp/test-data")

	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if registry.dataDir != "/tmp/test-data" {
		t.Error("dataDir should be set correctly")
	}
	if registry.AppRegistry == nil {
		t.Error("AppRegistry should be initialized")
	}

	// Should have hardcoded apps loaded
	apps := registry.GetAll()
	if len(apps) == 0 {
		t.Error("registry should have hardcoded apps loaded")
	}

	// Check for expected hardcoded apps
	expectedApps := []string{"vim", "zsh", "dwm", "st", "lf", "zathura"}
	for _, name := range expectedApps {
		if _, exists := registry.Get(name); !exists {
			t.Errorf("hardcoded app %s should exist", name)
		}
	}
}

func TestRegistry_LoadApp_FromFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create test app file
	testApp := &App{
		Name:        "test-app",
		Description: "Test application",
		Version:     "1.0",
		Categories:  []string{"test"},
		Shortcuts: []Shortcut{
			{Keys: "q", Description: "quit", Category: "general"},
		},
	}

	appData, err := yaml.Marshal(testApp)
	if err != nil {
		t.Fatalf("failed to marshal test app: %v", err)
	}

	appPath := filepath.Join(tmpDir, "test-app.yaml")
	if err := os.WriteFile(appPath, appData, 0644); err != nil {
		t.Fatalf("failed to write test app file: %v", err)
	}

	// Test loading
	registry := NewRegistry(tmpDir)
	err = registry.LoadApp("test-app")
	if err != nil {
		t.Errorf("LoadApp should succeed: %v", err)
	}

	// Verify app was loaded
	app, exists := registry.Get("test-app")
	if !exists {
		t.Error("loaded app should exist in registry")
	}
	if app.Name != "test-app" {
		t.Error("loaded app should have correct name")
	}
}

func TestRegistry_LoadApp_NotFound(t *testing.T) {
	registry := NewRegistry("/non/existent/path")

	err := registry.LoadApp("non-existent-app")
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}

func TestRegistry_LoadApp_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid YAML file
	invalidYAML := "invalid: yaml: content: ["
	appPath := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(appPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write invalid YAML file: %v", err)
	}

	registry := NewRegistry(tmpDir)
	err := registry.LoadApp("invalid")

	// Should fallback to hardcoded or return error
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound for invalid YAML, got %v", err)
	}
}

func TestRegistry_LoadApps(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test app files
	apps := []struct {
		name string
		app  *App
	}{
		{"app1", &App{Name: "app1", Description: "First app"}},
		{"app2", &App{Name: "app2", Description: "Second app"}},
	}

	for _, appInfo := range apps {
		appData, err := yaml.Marshal(appInfo.app)
		if err != nil {
			t.Fatalf("failed to marshal app %s: %v", appInfo.name, err)
		}

		appPath := filepath.Join(tmpDir, appInfo.name+".yaml")
		if err := os.WriteFile(appPath, appData, 0644); err != nil {
			t.Fatalf("failed to write app file %s: %v", appInfo.name, err)
		}
	}

	registry := NewRegistry(tmpDir)
	err := registry.LoadApps([]string{"app1", "app2", "non-existent"})

	// Should not return error even if some apps fail to load
	if err != nil {
		t.Errorf("LoadApps should not return error: %v", err)
	}

	// Verify loaded apps
	for _, appInfo := range apps {
		if _, exists := registry.Get(appInfo.name); !exists {
			t.Errorf("app %s should be loaded", appInfo.name)
		}
	}
}

func TestRegistry_GetTableData(t *testing.T) {
	registry := NewRegistry("")

	// Test with known hardcoded apps
	appNames := []string{"vim", "zsh"}
	tableData := registry.GetTableData(appNames)

	if len(tableData) == 0 {
		t.Error("table data should not be empty")
	}

	// Check header row
	header := tableData[0]
	if len(header) != len(appNames)+1 {
		t.Errorf("header should have %d columns, got %d", len(appNames)+1, len(header))
	}
	if header[0] != "Shortcut" {
		t.Error("first column should be 'Shortcut'")
	}
	for i, appName := range appNames {
		if header[i+1] != appName {
			t.Errorf("column %d should be %s, got %s", i+1, appName, header[i+1])
		}
	}

	// Check that we have data rows
	if len(tableData) < 2 {
		t.Error("should have at least header + data rows")
	}
}

func TestRegistry_GetTableData_EmptyApps(t *testing.T) {
	registry := NewRegistry("")
	tableData := registry.GetTableData([]string{})

	// Should return at least header
	if len(tableData) == 0 {
		t.Error("should return at least header row")
	}

	header := tableData[0]
	if len(header) != 1 || header[0] != "Shortcut" {
		t.Error("header should contain only 'Shortcut' column")
	}
}

func TestRegistry_loadAppFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewRegistry("")

	// Test valid file
	validApp := &App{Name: "valid", Description: "Valid app"}
	validData, _ := yaml.Marshal(validApp)
	validPath := filepath.Join(tmpDir, "valid.yaml")
	os.WriteFile(validPath, validData, 0644)

	app, err := registry.loadAppFromFile(validPath)
	if err != nil {
		t.Errorf("should load valid file: %v", err)
	}
	if app.Name != "valid" {
		t.Error("loaded app should have correct name")
	}

	// Test non-existent file
	_, err = registry.loadAppFromFile("/non/existent/file.yaml")
	if err == nil {
		t.Error("should return error for non-existent file")
	}

	// Test invalid YAML
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	os.WriteFile(invalidPath, []byte("invalid: yaml: ["), 0644)

	_, err = registry.loadAppFromFile(invalidPath)
	if !errors.Is(err, ErrInvalidAppFile) {
		t.Errorf("expected ErrInvalidAppFile, got %v", err)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, test := range tests {
		result := expandPath(test.input)
		if test.input != "~" && result != test.expected {
			t.Errorf("expandPath(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}

	// Test ~ expansion (result will vary by system)
	homeResult := expandPath("~")
	if homeResult == "~" {
		t.Error("~ should be expanded to home directory")
	}

	// Test ~/subpath expansion
	subpathResult := expandPath("~/test")
	if subpathResult == "~/test" {
		t.Error("~/test should be expanded")
	}
}

func TestRegistry_HardcodedAppsData(t *testing.T) {
	registry := NewRegistry("")

	// Test that hardcoded apps have expected structure
	vim, exists := registry.Get("vim")
	if !exists {
		t.Error("vim should exist in hardcoded apps")
	}
	if vim.Name != "vim" {
		t.Error("vim app should have correct name")
	}
	if len(vim.Shortcuts) == 0 {
		t.Error("vim should have shortcuts")
	}

	// Check that shortcuts have expected fields
	for _, shortcut := range vim.Shortcuts {
		if shortcut.Keys == "" {
			t.Error("shortcut should have keys")
		}
		if shortcut.Description == "" {
			t.Error("shortcut should have description")
		}
	}
}
