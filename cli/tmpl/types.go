/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-16 09:13
 * @ Modified time: 2025-11-16 09:13
 * @ Description: Type definitions for template generator
 */

package tmpl

// Config represents the configuration for project generation
type Config struct {
	TemplateName string
	OutputDir    string
	ProjectName  string
	Variables    map[string]any
	Force        bool
}

// TemplateInfo contains metadata about a template
type TemplateInfo struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Version     string        `yaml:"version"`
	Variables   []VariableDef `yaml:"variables"`
	Path        string        `yaml:"-"` // Internal: path to template directory
}

// VariableDef defines a template variable
type VariableDef struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     any    `yaml:"default"`
	Type        string `yaml:"type"` // string, int, bool, etc.
}

// Template represents a loaded template ready for rendering
type Template struct {
	Info     *TemplateInfo
	Files    []TemplateFile
	BasePath string
}

// TemplateFile represents a file in the template
type TemplateFile struct {
	SourcePath      string // Path in template directory
	TargetPath      string // Rendered path (with template variables)
	IsTemplate      bool   // Whether file contains template syntax
	Content         []byte // File content (if not template, use as-is)
	TemplateContent string // Template content (if IsTemplate is true)
}
