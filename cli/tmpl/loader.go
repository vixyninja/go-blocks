package tmpl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	templateYAML = "template.yaml"
)

// LoadTemplate loads a template from the file system
func LoadTemplate(templatePath string) (*Template, error) {
	// Check if template directory exists
	info, err := os.Stat(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("template path is not a directory: %s", templatePath)
	}

	// Load template.yaml
	templateYAMLPath := filepath.Join(templatePath, templateYAML)
	templateInfo, err := loadTemplateInfo(templateYAMLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template info: %w", err)
	}
	templateInfo.Path = templatePath

	// Load template files
	files, err := loadTemplateFiles(templatePath, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template files: %w", err)
	}

	return &Template{
		Info:     templateInfo,
		Files:    files,
		BasePath: templatePath,
	}, nil
}

// loadTemplateInfo loads template metadata from template.yaml
func loadTemplateInfo(yamlPath string) (*TemplateInfo, error) {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var info TemplateInfo
	if err := yaml.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse template.yaml: %w", err)
	}

	return &info, nil
}

// loadTemplateFiles recursively loads all files from template directory
func loadTemplateFiles(templatePath, basePath string) ([]TemplateFile, error) {
	var files []TemplateFile

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip template.yaml itself
		if info.Name() == templateYAML {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path from base
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Determine if file is a template
		// Files with .tmpl extension or containing {{ }} are templates
		isTemplate := strings.HasSuffix(path, ".tmpl") || strings.Contains(string(content), "{{")

		templateFile := TemplateFile{
			SourcePath: relPath,
			TargetPath: relPath,
			IsTemplate: isTemplate,
		}

		if isTemplate {
			// Remove .tmpl extension from target path
			if strings.HasSuffix(relPath, ".tmpl") {
				templateFile.TargetPath = strings.TrimSuffix(relPath, ".tmpl")
			}
			templateFile.TemplateContent = string(content)
		} else {
			templateFile.Content = content
		}

		files = append(files, templateFile)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// FindTemplate searches for a template in common locations
func FindTemplate(templateName string) (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Get executable directory (if available)
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	// Search in multiple locations
	possiblePaths := []string{
		// Relative to current working directory
		filepath.Join(cwd, "cli", "tmpl", "templates", templateName),
		filepath.Join(cwd, "tmpl", "templates", templateName),
		filepath.Join(cwd, "templates", templateName),
		// Relative to executable
		filepath.Join(execDir, "cli", "tmpl", "templates", templateName),
		filepath.Join(execDir, "tmpl", "templates", templateName),
		filepath.Join(execDir, "templates", templateName),
		// Direct path (absolute or relative)
		templateName,
		filepath.Join(cwd, templateName),
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			// Check if it has template.yaml
			if _, err := os.Stat(filepath.Join(path, templateYAML)); err == nil {
				absPath, err := filepath.Abs(path)
				if err == nil {
					return absPath, nil
				}
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("template '%s' not found. Searched in: %v", templateName, possiblePaths)
}
