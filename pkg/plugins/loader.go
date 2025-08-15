package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrPluginNotFound          = errors.New("plugin not found")
	ErrPluginAlreadyRegistered = errors.New("plugin already registered")
	ErrInvalidPlugin           = errors.New("invalid plugin")
	ErrPluginLoadFailed        = errors.New("plugin load failed")
)

type Loader struct {
	registry      *Registry
	pluginDirs    []string
	loadedPlugins map[string]*LoadedPlugin
}

type LoadedPlugin struct {
	Plugin   Plugin
	Metadata *Metadata
	Path     string
}

func NewLoader(dirs ...string) *Loader {
	if len(dirs) == 0 {
		dirs = getDefaultPluginDirs()
	}

	return &Loader{
		registry:      NewRegistry(),
		pluginDirs:    dirs,
		loadedPlugins: make(map[string]*LoadedPlugin),
	}
}

func getDefaultPluginDirs() []string {
	dirs := []string{}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		dirs = append(dirs, filepath.Join(homeDir, ".config", "cheat-go", "plugins"))
	}

	dirs = append(dirs, "/usr/local/share/cheat-go/plugins")
	dirs = append(dirs, "./plugins")

	return dirs
}

func (l *Loader) LoadAll() error {
	for _, dir := range l.pluginDirs {
		if err := l.LoadFromDirectory(dir); err != nil {
			continue
		}
	}
	return nil
}

func (l *Loader) LoadFromDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory does not exist: %s", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		if strings.HasSuffix(entry.Name(), ".so") {
			if err := l.LoadNativePlugin(path); err != nil {
				continue
			}
		} else if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {
			if err := l.LoadScriptPlugin(path); err != nil {
				continue
			}
		}
	}

	return nil
}

func (l *Loader) LoadNativePlugin(path string) error {
	// Native plugin loading is disabled for now due to Go plugin limitations
	// We'll use script-based plugins instead
	return fmt.Errorf("native plugin loading not yet implemented")
}

func (l *Loader) LoadScriptPlugin(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin file %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read plugin file %s: %w", path, err)
	}

	var metadata Metadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse plugin metadata %s: %w", path, err)
	}

	scriptPlugin := NewScriptPlugin(metadata, path)

	l.loadedPlugins[metadata.Name] = &LoadedPlugin{
		Plugin:   scriptPlugin,
		Metadata: &metadata,
		Path:     path,
	}

	return l.registry.Register(metadata.Name, scriptPlugin)
}

func (l *Loader) GetPlugin(name string) (Plugin, error) {
	return l.registry.Get(name)
}

func (l *Loader) ListPlugins() []*LoadedPlugin {
	plugins := make([]*LoadedPlugin, 0, len(l.loadedPlugins))
	for _, p := range l.loadedPlugins {
		plugins = append(plugins, p)
	}
	return plugins
}

func (l *Loader) UnloadPlugin(name string) error {
	if loaded, exists := l.loadedPlugins[name]; exists {
		if err := loaded.Plugin.Cleanup(); err != nil {
			return fmt.Errorf("failed to cleanup plugin %s: %w", name, err)
		}
		delete(l.loadedPlugins, name)
	}

	return l.registry.Unregister(name)
}

func (l *Loader) ExportPluginInfo(w io.Writer) error {
	info := make(map[string]*Metadata)
	for name, loaded := range l.loadedPlugins {
		info[name] = loaded.Metadata
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}
