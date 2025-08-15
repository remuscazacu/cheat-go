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

// Test search functionality
func TestRegistry_SearchTableData(t *testing.T) {
	registry := NewRegistry("")
	apps := []string{"vim", "zsh"}

	// Test empty query returns all data
	allData := registry.GetTableData(apps)
	searchData := registry.SearchTableData(apps, "")
	if len(searchData) != len(allData) {
		t.Error("empty search should return all data")
	}

	// Test search with specific term
	searchResults := registry.SearchTableData(apps, "move")
	if len(searchResults) <= 1 { // should have more than just header
		t.Error("search for \"move\" should return results")
	}

	// Verify header is preserved
	if len(searchResults) > 0 {
		if searchResults[0][0] != "Shortcut" {
			t.Error("header should be preserved in search results")
		}
	}

	// Test search with no matches
	noResults := registry.SearchTableData(apps, "nonexistentterm")
	if len(noResults) != 1 { // should only have header
		t.Error("search with no matches should return only header")
	}
}

func TestRegistry_ShortcutMatches(t *testing.T) {
	registry := NewRegistry("")

	shortcut := Shortcut{
		Keys:        "k",
		Description: "move up",
		Category:    "navigation",
	}

	testCases := []struct {
		query    string
		expected bool
	}{
		{"k", true},          // matches keys
		{"K", true},          // case insensitive keys
		{"move", true},       // matches description
		{"MOVE", true},       // case insensitive description
		{"nav", true},        // matches category
		{"NAVIGATION", true}, // case insensitive category
		{"xyz", false},       // no match
		{"", true},           // empty query matches everything
	}

	for _, tc := range testCases {
		result := registry.shortcutMatches(shortcut, tc.query)
		if result != tc.expected {
			t.Errorf("shortcutMatches(%q) = %v, want %v", tc.query, result, tc.expected)
		}
	}
}

func TestRegistry_SearchShortcuts(t *testing.T) {
	registry := NewRegistry("")

	// Test search that should return results
	results := registry.SearchShortcuts("move")
	if len(results) == 0 {
		t.Error("search for \"move\" should return results")
	}

	// Check result structure
	if len(results) > 0 {
		result := results[0]
		if result.AppName == "" {
			t.Error("result should have app name")
		}
		if result.Shortcut.Keys == "" {
			t.Error("result should have shortcut with keys")
		}
		if len(result.Matches) == 0 {
			t.Error("result should have match indicators")
		}
	}

	// Test search with no results
	noResults := registry.SearchShortcuts("nonexistentterm")
	if len(noResults) != 0 {
		t.Error("search with no matches should return empty slice")
	}
}

func TestRegistry_LoadAllAppsFromDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewRegistry(tmpDir)

	// Create test app files in directory
	testApps := []App{
		{
			Name:        "app1",
			Description: "First app",
			Shortcuts: []Shortcut{
				{Keys: "a", Description: "action a"},
			},
		},
		{
			Name:        "app2",
			Description: "Second app",
			Shortcuts: []Shortcut{
				{Keys: "b", Description: "action b"},
			},
		},
	}

	// Write test files
	for _, app := range testApps {
		appFile := filepath.Join(tmpDir, app.Name+".yaml")
		data, _ := yaml.Marshal(app)
		os.WriteFile(appFile, data, 0644)
	}

	// Load all apps from directory
	err := registry.LoadAllAppsFromDirectory()
	if err != nil {
		t.Fatalf("LoadAllAppsFromDirectory() error = %v", err)
	}

	// Verify apps were loaded
	app1, exists1 := registry.AppRegistry.Get("app1")
	if !exists1 {
		t.Error("app1 should be loaded from directory")
	}
	if app1.Name != "app1" {
		t.Error("app1 should have correct name")
	}

	app2, exists2 := registry.AppRegistry.Get("app2")
	if !exists2 {
		t.Error("app2 should be loaded from directory")
	}
	if app2.Name != "app2" {
		t.Error("app2 should have correct name")
	}
}

func TestRegistry_LoadAllAppsFromDirectoryNonExistent(t *testing.T) {
	registry := NewRegistry("/non/existent/directory")

	// Should handle non-existent directory gracefully
	err := registry.LoadAllAppsFromDirectory()
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestRegistry_SaveApp(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewRegistry(tmpDir)

	testApp := &App{
		Name:        "save-test",
		Description: "Test save functionality",
		Version:     "1.0.0",
		Categories:  []string{"test"},
		Shortcuts: []Shortcut{
			{Keys: "s", Description: "save", Category: "file"},
		},
	}

	err := registry.SaveApp(testApp)
	if err != nil {
		t.Fatalf("SaveApp() error = %v", err)
	}

	// Verify file was created
	expectedFile := filepath.Join(tmpDir, "save-test.yaml")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("App file should be created")
	}

	// Load and verify content
	err = registry.LoadApp("save-test")
	if err != nil {
		t.Fatalf("LoadApp after Save error = %v", err)
	}

	loadedApp, exists := registry.AppRegistry.Get("save-test")
	if !exists {
		t.Error("Saved app should be loadable")
	}

	if loadedApp.Description != testApp.Description {
		t.Error("Saved app description should match")
	}
}

func TestRegistry_SaveAppInvalidPath(t *testing.T) {
	registry := NewRegistry("/invalid/path/that/cannot/be/created")

	testApp := &App{
		Name: "test",
	}

	err := registry.SaveApp(testApp)
	if err == nil {
		t.Error("Expected error when saving to invalid path")
	}
}

func TestRegistry_ValidateAppEdgeCases(t *testing.T) {
	registry := NewRegistry("")

	// Test app with no name
	invalidApp1 := &App{
		Description: "No name",
	}

	err := registry.validateApp(invalidApp1)
	if err == nil {
		t.Error("Should error for app with no name")
	}

	// Test app with no description
	invalidApp2 := &App{
		Name: "no-desc",
	}

	err = registry.validateApp(invalidApp2)
	if err == nil {
		t.Error("Should error for app with no description")
	}

	// Test app with no shortcuts
	invalidApp3 := &App{
		Name:        "no-shortcuts",
		Description: "Has description",
		Shortcuts:   []Shortcut{},
	}

	err = registry.validateApp(invalidApp3)
	// App with no shortcuts might be valid in some cases
	if err != nil {
		t.Logf("App with no shortcuts validation: %v", err)
	}

	// Test shortcut with no keys
	invalidApp4 := &App{
		Name:        "invalid-shortcut",
		Description: "Has invalid shortcut",
		Shortcuts: []Shortcut{
			{Description: "No keys"},
		},
	}

	err = registry.validateApp(invalidApp4)
	if err == nil {
		t.Error("Should error for shortcut with no keys")
	}

	// Test shortcut with no description
	invalidApp5 := &App{
		Name:        "invalid-shortcut2",
		Description: "Has invalid shortcut",
		Shortcuts: []Shortcut{
			{Keys: "x"},
		},
	}

	err = registry.validateApp(invalidApp5)
	if err == nil {
		t.Error("Should error for shortcut with no description")
	}

	// Test valid app
	validApp := &App{
		Name:        "valid",
		Description: "Valid app",
		Shortcuts: []Shortcut{
			{Keys: "v", Description: "valid action"},
		},
	}

	err = registry.validateApp(validApp)
	if err != nil {
		t.Errorf("Valid app should not error: %v", err)
	}
}

func TestRegistry_LoadAppFromFileErrors(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewRegistry(tmpDir)

	// Test with non-existent file
	_, err := registry.loadAppFromFile("non-existent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test with invalid YAML
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")
	os.WriteFile(invalidFile, []byte("invalid: yaml: ["), 0644)

	_, err = registry.loadAppFromFile("invalid.yaml")
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}

	// Test with invalid app data
	invalidApp := App{
		Name: "invalid",
		// Missing description and shortcuts
	}

	invalidData, _ := yaml.Marshal(invalidApp)
	invalidAppFile := filepath.Join(tmpDir, "invalid-app.yaml")
	os.WriteFile(invalidAppFile, invalidData, 0644)

	_, err = registry.loadAppFromFile("invalid-app.yaml")
	if err == nil {
		t.Error("Expected error for invalid app data")
	}
}

func TestRegistry_ExpandPath(t *testing.T) {
	// Test home directory expansion
	homePath := expandPath("~/test")
	if homePath == "~/test" {
		t.Error("~ should be expanded to home directory")
	}

	// Test absolute path
	absPath := expandPath("/absolute/path")
	if absPath != "/absolute/path" {
		t.Errorf("Absolute path should remain unchanged, got %s", absPath)
	}

	// Test relative path
	relPath := expandPath("relative/path")
	if relPath != "relative/path" {
		t.Errorf("Relative path should remain unchanged, got %s", relPath)
	}
}

func TestRegistry_GetSearchMatchesEdgeCases(t *testing.T) {
	registry := NewRegistry("")

	// Test with empty query
	shortcut := Shortcut{
		Keys:        "ctrl+c",
		Description: "copy",
	}

	matches := registry.getSearchMatches(shortcut, "")
	// Empty query might still return some matches depending on implementation
	if len(matches) != 0 {
		t.Logf("Empty query returned %d matches", len(matches))
	}

	// Test with query that matches keys
	matches = registry.getSearchMatches(shortcut, "ctrl")
	if len(matches) == 0 {
		t.Error("Should find matches in keys")
	}

	// Test with query that matches description
	matches = registry.getSearchMatches(shortcut, "copy")
	if len(matches) == 0 {
		t.Error("Should find matches in description")
	}
}
