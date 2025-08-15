package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ScriptPlugin struct {
	metadata Metadata
	path     string
	config   map[string]interface{}
}

func NewScriptPlugin(metadata Metadata, path string) *ScriptPlugin {
	return &ScriptPlugin{
		metadata: metadata,
		path:     path,
		config:   metadata.Config,
	}
}

func (s *ScriptPlugin) Name() string {
	return s.metadata.Name
}

func (s *ScriptPlugin) Version() string {
	return s.metadata.Version
}

func (s *ScriptPlugin) Author() string {
	return s.metadata.Author
}

func (s *ScriptPlugin) Description() string {
	return s.metadata.Description
}

func (s *ScriptPlugin) Init(config map[string]interface{}) error {
	if config != nil {
		for k, v := range config {
			s.config[k] = v
		}
	}
	return nil
}

func (s *ScriptPlugin) Execute(ctx context.Context, args []string) error {
	interpreter := s.config["interpreter"]
	if interpreter == nil {
		interpreter = "sh"
	}

	cmd := exec.CommandContext(ctx, interpreter.(string), append([]string{s.path}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script execution failed: %w\nOutput: %s", err, output)
	}

	return nil
}

func (s *ScriptPlugin) Cleanup() error {
	return nil
}

type JSONExportPlugin struct {
	BasePlugin
}

type BasePlugin struct {
	name        string
	version     string
	author      string
	description string
}

func (b *BasePlugin) Name() string                                     { return b.name }
func (b *BasePlugin) Version() string                                  { return b.version }
func (b *BasePlugin) Author() string                                   { return b.author }
func (b *BasePlugin) Description() string                              { return b.description }
func (b *BasePlugin) Init(config map[string]interface{}) error         { return nil }
func (b *BasePlugin) Execute(ctx context.Context, args []string) error { return nil }
func (b *BasePlugin) Cleanup() error                                   { return nil }

func NewJSONExportPlugin() *JSONExportPlugin {
	return &JSONExportPlugin{
		BasePlugin: BasePlugin{
			name:        "json-export",
			version:     "1.0.0",
			author:      "cheat-go",
			description: "Export apps to JSON format",
		},
	}
}

func (j *JSONExportPlugin) Export(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func (j *JSONExportPlugin) SupportedFormats() []string {
	return []string{"json"}
}

type MarkdownExportPlugin struct {
	BasePlugin
}

func NewMarkdownExportPlugin() *MarkdownExportPlugin {
	return &MarkdownExportPlugin{
		BasePlugin: BasePlugin{
			name:        "markdown-export",
			version:     "1.0.0",
			author:      "cheat-go",
			description: "Export apps to Markdown format",
		},
	}
}

func (m *MarkdownExportPlugin) Export(apps interface{}) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("# Cheat Sheet\n\n")

	return []byte(sb.String()), nil
}

func (m *MarkdownExportPlugin) SupportedFormats() []string {
	return []string{"md", "markdown"}
}
