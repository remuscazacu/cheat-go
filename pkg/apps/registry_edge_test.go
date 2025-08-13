package apps

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath_EdgeCases(t *testing.T) {
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

	// Test ~/ path
	result = expandPath("~/subdir")
	if result == "~/subdir" {
		t.Error("~/subdir should be expanded")
	}
	if !filepath.IsAbs(result) {
		t.Error("expanded ~/subdir should be absolute path")
	}
}

func TestRegistry_LoadApp_FallbackScenarios(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewRegistry(tmpDir)

	// Test loading app that exists in hardcoded data but not in file
	err := registry.LoadApp("vim") // vim should exist in hardcoded data
	if err != nil {
		t.Errorf("LoadApp should succeed for hardcoded app: %v", err)
	}

	// Verify app exists
	app, exists := registry.Get("vim")
	if !exists {
		t.Error("hardcoded vim app should exist")
	}
	if app.Name != "vim" {
		t.Error("app should have correct name")
	}
}

func TestRegistry_LoadApp_FilePermissionError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory instead of a file (this will cause read error)
	appDir := filepath.Join(tmpDir, "badapp.yaml")
	os.Mkdir(appDir, 0755)

	registry := NewRegistry(tmpDir)
	err := registry.LoadApp("badapp")

	// Should fallback gracefully
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound for unreadable file, got %v", err)
	}
}

func TestRegistry_GetTableData_EdgeCases(t *testing.T) {
	registry := NewRegistry("")

	// Test with apps that don't exist
	tableData := registry.GetTableData([]string{"nonexistent1", "nonexistent2"})

	// Should still return header
	if len(tableData) == 0 {
		t.Error("should return at least header")
	}

	header := tableData[0]
	expectedLen := 3 // "Shortcut" + 2 app names
	if len(header) != expectedLen {
		t.Errorf("header should have %d columns, got %d", expectedLen, len(header))
	}

	if header[0] != "Shortcut" {
		t.Error("first column should be 'Shortcut'")
	}
}

func TestRegistry_loadAppFromFile_InvalidPath(t *testing.T) {
	registry := NewRegistry("")

	// Test with directory instead of file
	tmpDir := t.TempDir()

	_, err := registry.loadAppFromFile(tmpDir) // Directory, not file
	if err == nil {
		t.Error("should return error when trying to read directory as file")
	}

	// Test with completely invalid path
	_, err = registry.loadAppFromFile("/dev/null/invalid/path.yaml")
	if err == nil {
		t.Error("should return error for invalid path")
	}
}
