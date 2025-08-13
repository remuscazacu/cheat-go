package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_Save_EdgeCases(t *testing.T) {
	loader := NewLoader("")
	config := DefaultConfig()

	// Test save to invalid path (should fail)
	err := loader.Save(config, "/root/nonexistent/config.yaml")
	if err == nil {
		t.Error("Save to invalid path should return error")
	}

	// Test save with invalid config that can't be marshaled
	// (This is hard to test as Go's yaml.Marshal is quite robust)
}

func TestLoader_Save_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	loader := NewLoader("")
	config := &Config{
		Apps:    []string{"test"},
		Theme:   "dark",
		DataDir: "/test",
	}

	// Test successful save with directory creation
	err := loader.Save(config, configPath)
	if err != nil {
		t.Errorf("Save should succeed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file should exist after save")
	}

	// Verify directory was created
	if _, err := os.Stat(filepath.Dir(configPath)); os.IsNotExist(err) {
		t.Error("directory should be created")
	}
}

func TestExpandPath_ErrorCase(t *testing.T) {
	// Test path that doesn't start with ~
	result := expandPath("/absolute/path")
	if result != "/absolute/path" {
		t.Errorf("absolute path should not be modified: got %s", result)
	}

	// Test empty path
	result = expandPath("")
	if result != "" {
		t.Errorf("empty path should remain empty: got %s", result)
	}

	// Test relative path
	result = expandPath("relative/path")
	if result != "relative/path" {
		t.Errorf("relative path should not be modified: got %s", result)
	}

	// Test ~ only
	result = expandPath("~")
	if result == "~" {
		t.Error("~ should be expanded")
	}
	if !filepath.IsAbs(result) {
		t.Error("expanded ~ should be absolute path")
	}
}
