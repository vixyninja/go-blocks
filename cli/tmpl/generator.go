package tmpl

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateProject generates a new project from a template
func GenerateProject(config *Config) error {
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Find and load template
	templatePath, err := FindTemplate(config.TemplateName)
	if err != nil {
		return err
	}

	template, err := LoadTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Validate template structure
	if err := ValidateTemplate(template); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Prepare variables with defaults
	variables := prepareVariables(template, config)

	// Validate required variables
	if err := ValidateVariables(template, variables); err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	// Prepare output directory
	if err := PrepareOutputDir(config.OutputDir, config.Force); err != nil {
		return fmt.Errorf("failed to prepare output directory: %w", err)
	}

	// Render template
	renderedFiles, err := RenderTemplate(template, variables)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Write rendered files
	if err := writeFiles(config.OutputDir, renderedFiles); err != nil {
		return fmt.Errorf("failed to write files: %w", err)
	}

	return nil
}

// prepareVariables merges user variables with template defaults
func prepareVariables(tmpl *Template, config *Config) map[string]any {
	variables := make(map[string]any)

	// Set defaults from template
	if tmpl.Info != nil {
		for _, varDef := range tmpl.Info.Variables {
			if varDef.Default != nil {
				variables[varDef.Name] = varDef.Default
			}
		}
	}

	// Override with user-provided variables
	for k, v := range config.Variables {
		variables[k] = v
	}

	// Always set ProjectName
	variables["ProjectName"] = config.ProjectName

	// Set CreateTime (current timestamp)
	variables["CreateTime"] = getCurrentTime()

	return variables
}

// getCurrentTime returns current time in a readable format
func getCurrentTime() string {
	// Format: 2025-11-14 12:00:00
	return time.Now().Format("2006-01-02 15:04:05")
}

// writeFiles writes rendered files to the output directory
func writeFiles(outputDir string, files []RenderedFile) error {
	for _, file := range files {
		targetPath := filepath.Join(outputDir, file.TargetPath)

		// Create directory if needed
		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(targetPath, file.Content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
	}

	return nil
}

// ListTemplates returns available templates
func ListTemplates() ([]TemplateInfo, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Get executable directory (if available)
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	// Try alternative paths
	possiblePaths := []string{
		filepath.Join(cwd, "cli", "tmpl", "templates"),
		filepath.Join(cwd, "tmpl", "templates"),
		filepath.Join(cwd, "templates"),
		filepath.Join(execDir, "cli", "tmpl", "templates"),
		filepath.Join(execDir, "tmpl", "templates"),
		filepath.Join(execDir, "templates"),
	}

	var templates []TemplateInfo

	for _, basePath := range possiblePaths {
		if info, err := os.Stat(basePath); err == nil && info.IsDir() {
			entries, err := os.ReadDir(basePath)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					templatePath := filepath.Join(basePath, entry.Name())
					template, err := LoadTemplate(templatePath)
					if err == nil {
						templates = append(templates, *template.Info)
					}
				}
			}
			if len(templates) > 0 {
				break // Found templates, stop searching
			}
		}
	}

	return templates, nil
}

// GetTemplateInfo returns information about a specific template
func GetTemplateInfo(templateName string) (*TemplateInfo, error) {
	templatePath, err := FindTemplate(templateName)
	if err != nil {
		return nil, err
	}

	template, err := LoadTemplate(templatePath)
	if err != nil {
		return nil, err
	}

	return template.Info, nil
}
