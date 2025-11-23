package tmpl

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateConfig validates the generation configuration
func ValidateConfig(config *Config) error {
	if config.TemplateName == "" {
		return fmt.Errorf("template name is required")
	}

	if config.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if config.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Check if output directory exists and is not empty (unless Force is set)
	if !config.Force {
		if info, err := os.Stat(config.OutputDir); err == nil {
			if info.IsDir() {
				// Check if directory is empty
				entries, err := os.ReadDir(config.OutputDir)
				if err != nil {
					return fmt.Errorf("failed to read output directory: %w", err)
				}
				if len(entries) > 0 {
					return fmt.Errorf("output directory %s is not empty (use --force to overwrite)", config.OutputDir)
				}
			} else {
				return fmt.Errorf("output path %s exists and is not a directory", config.OutputDir)
			}
		}
	}

	return nil
}

// ValidateTemplate validates template structure
func ValidateTemplate(template *Template) error {
	if template.Info == nil {
		return fmt.Errorf("template info is missing")
	}

	if template.Info.Name == "" {
		return fmt.Errorf("template name is missing")
	}

	if len(template.Files) == 0 {
		return fmt.Errorf("template has no files")
	}

	return nil
}

// ValidateVariables validates that all required variables are provided
func ValidateVariables(template *Template, variables map[string]any) error {
	if template.Info == nil {
		return fmt.Errorf("template info is missing")
	}

	for _, varDef := range template.Info.Variables {
		if varDef.Required {
			value, exists := variables[varDef.Name]
			if !exists || value == nil || value == "" {
				return fmt.Errorf("required variable '%s' is missing or empty", varDef.Name)
			}
		}
	}

	return nil
}

// PrepareOutputDir prepares the output directory
func PrepareOutputDir(outputDir string, force bool) error {
	// Create parent directory if needed
	parentDir := filepath.Dir(outputDir)
	if parentDir != "." && parentDir != "/" {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
	}

	// Remove existing directory if force is set
	if force {
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}
