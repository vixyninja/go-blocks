/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-16 09:06
 * @ Modified time: 2025-11-16 09:10
 * @ Description: Template renderer using Go text/template
 */

package tmpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

// RenderTemplate renders a template with the given variables
func RenderTemplate(tmpl *Template, variables map[string]any) ([]RenderedFile, error) {
	var renderedFiles []RenderedFile

	// Prepare template functions
	funcMap := template.FuncMap{
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"title":      toTitle,
		"trim":       strings.TrimSpace,
		"replace":    strings.ReplaceAll,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"join":       strings.Join,
		"split":      strings.Split,
		"filepath":   filepath.Join,
		"base":       filepath.Base,
		"dir":        filepath.Dir,
		"ext":        filepath.Ext,
		"noext":      func(s string) string { return strings.TrimSuffix(s, filepath.Ext(s)) },
		"snakecase":  toSnakeCase,
		"camelcase":  toCamelCase,
		"pascalcase": toPascalCase,
		"kebabcase":  toKebabCase,
	}

	for _, file := range tmpl.Files {
		renderedFile := RenderedFile{
			TargetPath: file.TargetPath,
		}

		if file.IsTemplate {
			// Render template content
			tmpl, err := template.New(file.SourcePath).Funcs(funcMap).Parse(file.TemplateContent)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template %s: %w", file.SourcePath, err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, variables); err != nil {
				return nil, fmt.Errorf("failed to render template %s: %w", file.SourcePath, err)
			}

			// Also render the target path if it contains template variables
			if strings.Contains(file.TargetPath, "{{") {
				pathTmpl, err := template.New("path").Funcs(funcMap).Parse(file.TargetPath)
				if err != nil {
					return nil, fmt.Errorf("failed to parse path template %s: %w", file.TargetPath, err)
				}
				var pathBuf bytes.Buffer
				if err := pathTmpl.Execute(&pathBuf, variables); err != nil {
					return nil, fmt.Errorf("failed to render path template %s: %w", file.TargetPath, err)
				}
				renderedFile.TargetPath = pathBuf.String()
			}

			renderedFile.Content = buf.Bytes()
		} else {
			// Use content as-is, but still render path if needed
			if strings.Contains(file.TargetPath, "{{") {
				pathTmpl, err := template.New("path").Funcs(funcMap).Parse(file.TargetPath)
				if err != nil {
					return nil, fmt.Errorf("failed to parse path template %s: %w", file.TargetPath, err)
				}
				var pathBuf bytes.Buffer
				if err := pathTmpl.Execute(&pathBuf, variables); err != nil {
					return nil, fmt.Errorf("failed to render path template %s: %w", file.TargetPath, err)
				}
				renderedFile.TargetPath = pathBuf.String()
			}
			renderedFile.Content = file.Content
		}

		renderedFiles = append(renderedFiles, renderedFile)
	}

	return renderedFiles, nil
}

// RenderedFile represents a file after template rendering
type RenderedFile struct {
	TargetPath string
	Content    []byte
}

// toTitle capitalizes the first letter of each word (replacement for deprecated strings.Title)
func toTitle(s string) string {
	if s == "" {
		return s
	}
	// Split into words and capitalize first letter of each
	words := strings.Fields(s)
	result := make([]string, len(words))
	for i, word := range words {
		result[i] = toTitleWord(word)
	}
	return strings.Join(result, " ")
}

// toUpperRune converts a rune to uppercase
func toUpperRune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}
	return r
}

// toTitleWord capitalizes the first letter of a word
func toTitleWord(word string) string {
	if word == "" {
		return word
	}
	runes := []rune(word)
	if len(runes) == 0 {
		return word
	}
	runes[0] = toUpperRune(runes[0])
	return string(runes)
}

// String case conversion helpers
func toSnakeCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}

func toCamelCase(s string) string {
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	if len(words) == 0 {
		return ""
	}
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += toTitleWord(strings.ToLower(words[i]))
	}
	return result
}

func toPascalCase(s string) string {
	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	result := ""
	for _, word := range words {
		result += toTitleWord(strings.ToLower(word))
	}
	return result
}

func toKebabCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}
