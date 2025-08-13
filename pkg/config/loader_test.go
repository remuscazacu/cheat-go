package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader("/test/path")

	if loader == nil {
		t.Fatal("NewLoader() returned nil")
	}

	if loader.configPath != "/test/path" {
		t.Errorf("configPath = %s, expected '/test/path'", loader.configPath)
	}
}

func TestLoader_Load_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	testConfig := &Config{
		Apps:  []string{"test-app"},
		Theme: "dark",
		Layout: LayoutConfig{
			Columns:        []string{"shortcut", "description"},
			ShowCategories: false,
			TableStyle:     "simple",
			MaxWidth:       80,
		},
		Keybinds: map[string]string{
			"quit":  "q",
			"up":    "k",
			"down":  "j",
			"left":  "h",
			"right": "l",
		},
		DataDir: "/test/data",
	}

	configData, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Test loading
	loader := NewLoader(configPath)
	config, err := loader.Load()

	if err != nil {
		t.Errorf("Load() failed: %v", err)
	}

	if !reflect.DeepEqual(config.Apps, testConfig.Apps) {
		t.Errorf("loaded Apps = %v, expected %v", config.Apps, testConfig.Apps)
	}

	if config.Theme != testConfig.Theme {
		t.Errorf("loaded Theme = %s, expected %s", config.Theme, testConfig.Theme)
	}
}

func TestLoader_Load_DefaultFallback(t *testing.T) {
	// Test with non-existent path
	loader := NewLoader("/non/existent/config.yaml")
	config, err := loader.Load()

	if err != nil {
		t.Errorf("Load() should not return error on fallback: %v", err)
	}

	// Should return default config
	defaultConfig := DefaultConfig()
	if !reflect.DeepEqual(config.Apps, defaultConfig.Apps) {
		t.Error("should fallback to default config")
	}
}

func TestLoader_Load_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// Write invalid YAML
	invalidYAML := "apps: [\ninvalid yaml"
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	loader := NewLoader(configPath)
	config, err := loader.Load()

	// Should fallback to default config
	if err != nil {
		t.Errorf("Load() should not return error on invalid YAML: %v", err)
	}

	defaultConfig := DefaultConfig()
	if !reflect.DeepEqual(config.Apps, defaultConfig.Apps) {
		t.Error("should fallback to default config on invalid YAML")
	}
}

func TestLoader_LoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewLoader("")

	// Test valid file
	validConfig := &Config{Apps: []string{"test"}}
	validData, _ := yaml.Marshal(validConfig)
	validPath := filepath.Join(tmpDir, "valid.yaml")
	os.WriteFile(validPath, validData, 0644)

	config, err := loader.loadFromFile(validPath)
	if err != nil {
		t.Errorf("should load valid file: %v", err)
	}
	if len(config.Apps) == 0 || config.Apps[0] != "test" {
		t.Error("loaded config should have correct apps")
	}

	// Test non-existent file
	_, err = loader.loadFromFile("/non/existent/file.yaml")
	if err == nil {
		t.Error("should return error for non-existent file")
	}

	// Test invalid YAML
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	os.WriteFile(invalidPath, []byte("invalid: yaml: ["), 0644)

	_, err = loader.loadFromFile(invalidPath)
	if err != ErrInvalidConfig {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestLoader_ValidateAndSetDefaults(t *testing.T) {
	loader := NewLoader("")

	// Test partial config
	partialConfig := &Config{
		Theme: "dark",
		// Missing Apps, Layout, Keybinds, DataDir
	}

	config, err := loader.validateAndSetDefaults(partialConfig)
	if err != nil {
		t.Errorf("validateAndSetDefaults failed: %v", err)
	}

	// Should fill in defaults
	defaultConfig := DefaultConfig()
	if !reflect.DeepEqual(config.Apps, defaultConfig.Apps) {
		t.Error("should set default Apps")
	}

	if config.Theme != "dark" {
		t.Error("should preserve existing Theme")
	}

	if !reflect.DeepEqual(config.Layout.Columns, defaultConfig.Layout.Columns) {
		t.Error("should set default Layout")
	}

	if !reflect.DeepEqual(config.Keybinds, defaultConfig.Keybinds) {
		t.Error("should set default Keybinds")
	}

	if config.DataDir != defaultConfig.DataDir {
		t.Error("should set default DataDir")
	}
}

func TestLoader_ValidateAndSetDefaults_CompleteConfig(t *testing.T) {
	loader := NewLoader("")

	// Test complete config
	completeConfig := &Config{
		Apps:  []string{"custom-app"},
		Theme: "dark",
		Layout: LayoutConfig{
			Columns:        []string{"shortcut", "description"},
			ShowCategories: true,
			TableStyle:     "rounded",
			MaxWidth:       100,
		},
		Keybinds: map[string]string{
			"quit":   "q",
			"up":     "k",
			"down":   "j",
			"left":   "h",
			"right":  "l",
			"custom": "c",
		},
		DataDir: "/custom/path",
	}

	config, err := loader.validateAndSetDefaults(completeConfig)
	if err != nil {
		t.Errorf("validateAndSetDefaults failed: %v", err)
	}

	// Should preserve all values
	if !reflect.DeepEqual(config.Apps, completeConfig.Apps) {
		t.Error("should preserve Apps")
	}

	if config.Theme != completeConfig.Theme {
		t.Error("should preserve Theme")
	}

	if !reflect.DeepEqual(config.Layout.Columns, completeConfig.Layout.Columns) {
		t.Error("should preserve Layout")
	}

	if !reflect.DeepEqual(config.Keybinds, completeConfig.Keybinds) {
		t.Error("should preserve Keybinds")
	}

	if config.DataDir != completeConfig.DataDir {
		t.Error("should preserve DataDir")
	}
}

func TestLoader_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	loader := NewLoader("")
	config := &Config{
		Apps:    []string{"test-app"},
		Theme:   "test-theme",
		DataDir: "/test/dir",
	}

	err := loader.Save(config, configPath)
	if err != nil {
		t.Errorf("Save() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file should be created")
	}

	// Verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	var savedConfig Config
	if err := yaml.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("saved config should be valid YAML: %v", err)
	}

	if !reflect.DeepEqual(savedConfig.Apps, config.Apps) {
		t.Error("saved config should match original")
	}
}

func TestLoader_Save_CreateDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nestedPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")

	loader := NewLoader("")
	config := DefaultConfig()

	err := loader.Save(config, nestedPath)
	if err != nil {
		t.Errorf("Save() should create nested directories: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Error("config file should be created in nested directory")
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

	// Test ~ expansion
	homeResult := expandPath("~")
	if homeResult == "~" {
		t.Error("~ should be expanded to home directory")
	}

	// Test ~/subpath expansion
	subpathResult := expandPath("~/test")
	if subpathResult == "~/test" {
		t.Error("~/test should be expanded")
	}

	if !filepath.IsAbs(subpathResult) {
		t.Error("expanded path should be absolute")
	}
}

func TestLoader_Load_DefaultPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config in one of the default locations (simulate current directory)
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	testConfig := &Config{Apps: []string{"default-path-test"}}
	configData, _ := yaml.Marshal(testConfig)
	os.WriteFile("config.yaml", configData, 0644)

	// Test loading without explicit path
	loader := NewLoader("")
	config, err := loader.Load()

	if err != nil {
		t.Errorf("Load() should find default config: %v", err)
	}

	if len(config.Apps) == 0 || config.Apps[0] != "default-path-test" {
		t.Error("should load config from default path")
	}
}
