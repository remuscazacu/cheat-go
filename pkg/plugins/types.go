package plugins

import (
	"cheat-go/pkg/apps"
	"context"
)

type Plugin interface {
	Name() string
	Version() string
	Author() string
	Description() string
	Init(config map[string]interface{}) error
	Execute(ctx context.Context, args []string) error
	Cleanup() error
}

type AppProviderPlugin interface {
	Plugin
	GetApps() ([]apps.App, error)
	LoadApp(name string) (*apps.App, error)
}

type TransformPlugin interface {
	Plugin
	Transform(app *apps.App) (*apps.App, error)
}

type ExportPlugin interface {
	Plugin
	Export(apps []apps.App, format string) ([]byte, error)
	SupportedFormats() []string
}

type ImportPlugin interface {
	Plugin
	Import(data []byte, format string) ([]apps.App, error)
	SupportedFormats() []string
}

type Metadata struct {
	Name        string                 `json:"name" yaml:"name"`
	Version     string                 `json:"version" yaml:"version"`
	Author      string                 `json:"author" yaml:"author"`
	Description string                 `json:"description" yaml:"description"`
	Type        string                 `json:"type" yaml:"type"`
	Config      map[string]interface{} `json:"config" yaml:"config"`
}

type Registry struct {
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

func (r *Registry) Register(name string, plugin Plugin) error {
	if _, exists := r.plugins[name]; exists {
		return ErrPluginAlreadyRegistered
	}
	r.plugins[name] = plugin
	return nil
}

func (r *Registry) Get(name string) (Plugin, error) {
	plugin, exists := r.plugins[name]
	if !exists {
		return nil, ErrPluginNotFound
	}
	return plugin, nil
}

func (r *Registry) List() []string {
	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

func (r *Registry) Unregister(name string) error {
	if _, exists := r.plugins[name]; !exists {
		return ErrPluginNotFound
	}
	delete(r.plugins, name)
	return nil
}
