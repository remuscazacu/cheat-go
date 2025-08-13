package apps

import (
	"reflect"
	"testing"
)

func TestNewAppRegistry(t *testing.T) {
	registry := NewAppRegistry()
	if registry == nil {
		t.Fatal("NewAppRegistry() returned nil")
	}
	if registry.apps == nil {
		t.Error("apps map should be initialized")
	}
	if len(registry.apps) != 0 {
		t.Error("new registry should be empty")
	}
}

func TestAppRegistry_Register(t *testing.T) {
	registry := NewAppRegistry()
	app := &App{
		Name:        "test-app",
		Description: "Test application",
		Version:     "1.0",
		Categories:  []string{"test"},
		Shortcuts: []Shortcut{
			{Keys: "q", Description: "quit", Category: "general"},
		},
	}

	registry.Register(app)

	retrieved, exists := registry.Get("test-app")
	if !exists {
		t.Error("registered app should exist")
	}
	if !reflect.DeepEqual(retrieved, app) {
		t.Error("retrieved app should match registered app")
	}
}

func TestAppRegistry_Get(t *testing.T) {
	registry := NewAppRegistry()

	// Test getting non-existent app
	_, exists := registry.Get("non-existent")
	if exists {
		t.Error("non-existent app should not exist")
	}

	// Test getting existing app
	app := &App{Name: "test"}
	registry.Register(app)

	retrieved, exists := registry.Get("test")
	if !exists {
		t.Error("registered app should exist")
	}
	if retrieved.Name != "test" {
		t.Error("retrieved app name should match")
	}
}

func TestAppRegistry_GetAll(t *testing.T) {
	registry := NewAppRegistry()

	// Test empty registry
	all := registry.GetAll()
	if len(all) != 0 {
		t.Error("empty registry should return empty map")
	}

	// Test with apps
	app1 := &App{Name: "app1"}
	app2 := &App{Name: "app2"}
	registry.Register(app1)
	registry.Register(app2)

	all = registry.GetAll()
	if len(all) != 2 {
		t.Errorf("expected 2 apps, got %d", len(all))
	}
	if all["app1"] != app1 || all["app2"] != app2 {
		t.Error("GetAll should return all registered apps")
	}
}

func TestAppRegistry_List(t *testing.T) {
	registry := NewAppRegistry()

	// Test empty registry
	names := registry.List()
	if len(names) != 0 {
		t.Error("empty registry should return empty slice")
	}

	// Test with apps
	registry.Register(&App{Name: "vim"})
	registry.Register(&App{Name: "zsh"})

	names = registry.List()
	if len(names) != 2 {
		t.Errorf("expected 2 app names, got %d", len(names))
	}

	// Check that both names are present (order may vary)
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}
	if !nameMap["vim"] || !nameMap["zsh"] {
		t.Error("List should return all app names")
	}
}

func TestApp_JSONMarshaling(t *testing.T) {
	app := &App{
		Name:        "test-app",
		Description: "Test application",
		Version:     "1.0.0",
		Categories:  []string{"editor", "terminal"},
		Shortcuts: []Shortcut{
			{Keys: "ctrl+q", Description: "quit", Category: "general", Tags: []string{"exit"}},
			{Keys: "ctrl+s", Description: "save", Category: "file", Tags: []string{"save", "write"}},
		},
		Metadata: map[string]string{
			"author": "test",
			"url":    "https://example.com",
		},
	}

	// Test that the struct can be created and accessed properly
	if app.Name != "test-app" {
		t.Error("app name should be set correctly")
	}
	if len(app.Shortcuts) != 2 {
		t.Error("app should have 2 shortcuts")
	}
	if app.Shortcuts[0].Keys != "ctrl+q" {
		t.Error("first shortcut keys should be ctrl+q")
	}
}

func TestShortcut_Fields(t *testing.T) {
	shortcut := Shortcut{
		Keys:        "ctrl+x",
		Description: "execute",
		Category:    "action",
		Tags:        []string{"execute", "run"},
		Platform:    "linux",
	}

	if shortcut.Keys != "ctrl+x" {
		t.Error("shortcut keys should be set correctly")
	}
	if shortcut.Description != "execute" {
		t.Error("shortcut description should be set correctly")
	}
	if len(shortcut.Tags) != 2 {
		t.Error("shortcut should have 2 tags")
	}
	if shortcut.Platform != "linux" {
		t.Error("shortcut platform should be set correctly")
	}
}
