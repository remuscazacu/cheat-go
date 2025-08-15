package plugins

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader(t *testing.T) {
	tests := []struct {
		name string
		dirs []string
		want int
	}{
		{
			name: "Default directories",
			dirs: nil,
			want: 3,
		},
		{
			name: "Custom directories",
			dirs: []string{"/tmp/plugins", "/opt/plugins"},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewLoader(tt.dirs...)
			if len(loader.pluginDirs) != tt.want {
				t.Errorf("NewLoader() dirs = %v, want %v", len(loader.pluginDirs), tt.want)
			}
		})
	}
}

func TestLoader_LoadScriptPlugin(t *testing.T) {
	tempDir := t.TempDir()

	yamlContent := `
name: test-plugin
version: 1.0.0
author: test
description: Test plugin
type: script
config:
  interpreter: sh
`

	pluginPath := filepath.Join(tempDir, "test.yaml")
	if err := os.WriteFile(pluginPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader(tempDir)

	err := loader.LoadScriptPlugin(pluginPath)
	if err != nil {
		t.Fatalf("LoadScriptPlugin() error = %v", err)
	}

	plugin, err := loader.GetPlugin("test-plugin")
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}

	if plugin.Name() != "test-plugin" {
		t.Errorf("Plugin name = %v, want %v", plugin.Name(), "test-plugin")
	}
}

func TestLoader_LoadFromDirectory(t *testing.T) {
	tempDir := t.TempDir()

	yamlContent1 := `
name: plugin1
version: 1.0.0
author: test
description: Plugin 1
`

	yamlContent2 := `
name: plugin2
version: 1.0.0
author: test
description: Plugin 2
`

	os.WriteFile(filepath.Join(tempDir, "plugin1.yaml"), []byte(yamlContent1), 0644)
	os.WriteFile(filepath.Join(tempDir, "plugin2.yaml"), []byte(yamlContent2), 0644)
	os.WriteFile(filepath.Join(tempDir, "notaplugin.txt"), []byte("text"), 0644)

	loader := NewLoader()

	err := loader.LoadFromDirectory(tempDir)
	if err != nil {
		t.Fatalf("LoadFromDirectory() error = %v", err)
	}

	plugins := loader.ListPlugins()
	if len(plugins) != 2 {
		t.Errorf("ListPlugins() = %v plugins, want 2", len(plugins))
	}
}

func TestLoader_UnloadPlugin(t *testing.T) {
	tempDir := t.TempDir()

	yamlContent := `
name: unload-test
version: 1.0.0
author: test
description: Unload test
`

	pluginPath := filepath.Join(tempDir, "unload.yaml")
	os.WriteFile(pluginPath, []byte(yamlContent), 0644)

	loader := NewLoader(tempDir)
	loader.LoadScriptPlugin(pluginPath)

	if err := loader.UnloadPlugin("unload-test"); err != nil {
		t.Fatalf("UnloadPlugin() error = %v", err)
	}

	_, err := loader.GetPlugin("unload-test")
	if err != ErrPluginNotFound {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}
}

func TestLoader_ExportPluginInfo(t *testing.T) {
	tempDir := t.TempDir()

	yamlContent := `
name: export-test
version: 1.0.0
author: test
description: Export test
`

	pluginPath := filepath.Join(tempDir, "export.yaml")
	os.WriteFile(pluginPath, []byte(yamlContent), 0644)

	loader := NewLoader(tempDir)
	loader.LoadScriptPlugin(pluginPath)

	var buf bytes.Buffer
	if err := loader.ExportPluginInfo(&buf); err != nil {
		t.Fatalf("ExportPluginInfo() error = %v", err)
	}

	var info map[string]*Metadata
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal exported info: %v", err)
	}

	if _, exists := info["export-test"]; !exists {
		t.Error("Exported info doesn't contain export-test plugin")
	}
}

func TestScriptPlugin(t *testing.T) {
	metadata := Metadata{
		Name:        "test-script",
		Version:     "1.0.0",
		Author:      "test",
		Description: "Test script plugin",
		Config: map[string]interface{}{
			"interpreter": "sh",
		},
	}

	plugin := NewScriptPlugin(metadata, "/path/to/script.sh")

	if plugin.Name() != "test-script" {
		t.Errorf("Name() = %v, want %v", plugin.Name(), "test-script")
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Version() = %v, want %v", plugin.Version(), "1.0.0")
	}

	if plugin.Author() != "test" {
		t.Errorf("Author() = %v, want %v", plugin.Author(), "test")
	}

	err := plugin.Init(map[string]interface{}{"new_key": "value"})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if plugin.config["new_key"] != "value" {
		t.Error("Init() didn't update config")
	}
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	plugin := &BasePlugin{
		name:        "test",
		version:     "1.0.0",
		author:      "test",
		description: "Test plugin",
	}

	// Test Register
	err := registry.Register("test", plugin)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Test duplicate registration
	err = registry.Register("test", plugin)
	if err != ErrPluginAlreadyRegistered {
		t.Errorf("Expected ErrPluginAlreadyRegistered, got %v", err)
	}

	// Test Get
	got, err := registry.Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != plugin {
		t.Error("Get() returned wrong plugin")
	}

	// Test Get non-existent
	_, err = registry.Get("nonexistent")
	if err != ErrPluginNotFound {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}

	// Test List
	names := registry.List()
	if len(names) != 1 || names[0] != "test" {
		t.Errorf("List() = %v, want [test]", names)
	}

	// Test Unregister
	err = registry.Unregister("test")
	if err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	// Test Unregister non-existent
	err = registry.Unregister("nonexistent")
	if err != ErrPluginNotFound {
		t.Errorf("Expected ErrPluginNotFound, got %v", err)
	}
}
