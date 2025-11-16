/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-14 11:17:03
 * @ Modified time: 2025-11-14 12:00:00
 * @ Description: Scaffold command for generating projects from templates
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vixyninja/go-blocks/cli/tmpl"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate projects from templates",
	Long: `Scaffold generates new projects from predefined templates.
You can list available templates or generate a new project using a template.`,
}

var scaffoldNewCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Generate a new project from a template",
	Long: `Generate a new project from a template.
Example:
  cli scaffold new myproject --template go-service --output ./projects/myproject`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		templateName, _ := cmd.Flags().GetString("template")
		outputDir, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")

		if templateName == "" {
			return fmt.Errorf("template name is required (use --template)")
		}

		// Default output directory to current directory + project name
		if outputDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			outputDir = filepath.Join(cwd, projectName)
		}

		// Get template info to collect variables
		templateInfo, err := tmpl.GetTemplateInfo(templateName)
		if err != nil {
			return fmt.Errorf("failed to get template info: %w", err)
		}

		// Prepare variables
		variables := make(map[string]any)
		variables["ProjectName"] = projectName

		// Collect variables from flags or use defaults
		for _, varDef := range templateInfo.Variables {
			if varDef.Name == "ProjectName" {
				continue // Already set
			}

			// Try to get from flag
			flag := cmd.Flags().Lookup(varDef.Name)
			if flag != nil {
				value, _ := cmd.Flags().GetString(varDef.Name)
				if value != "" {
					variables[varDef.Name] = value
				} else if varDef.Default != nil {
					variables[varDef.Name] = varDef.Default
				}
			} else if varDef.Default != nil {
				variables[varDef.Name] = varDef.Default
			}
		}

		// Special handling for ModulePath - default to github.com/user/projectname
		if _, exists := variables["ModulePath"]; !exists || variables["ModulePath"] == "" {
			modulePath, _ := cmd.Flags().GetString("module")
			if modulePath == "" {
				// Default module path
				variables["ModulePath"] = fmt.Sprintf("github.com/%s/%s", "user", projectName)
			} else {
				variables["ModulePath"] = modulePath
			}
		}

		// Generate project
		config := &tmpl.Config{
			TemplateName: templateName,
			OutputDir:    outputDir,
			ProjectName:  projectName,
			Variables:    variables,
			Force:        force,
		}

		fmt.Printf("Generating project '%s' from template '%s'...\n", projectName, templateName)
		fmt.Printf("Output directory: %s\n", outputDir)

		if err := tmpl.GenerateProject(config); err != nil {
			return fmt.Errorf("failed to generate project: %w", err)
		}

		fmt.Printf("✓ Project '%s' generated successfully!\n", projectName)
		return nil
	},
}

var scaffoldListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all available templates that can be used to generate projects.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		templates, err := tmpl.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}

		if len(templates) == 0 {
			fmt.Println("No templates found.")
			return nil
		}

		fmt.Println("Available templates:")
		fmt.Println()
		for _, tmpl := range templates {
			fmt.Printf("  %s\n", tmpl.Name)
			if tmpl.Description != "" {
				fmt.Printf("    Description: %s\n", tmpl.Description)
			}
			if tmpl.Version != "" {
				fmt.Printf("    Version: %s\n", tmpl.Version)
			}
			if len(tmpl.Variables) > 0 {
				fmt.Printf("    Variables:\n")
				for _, v := range tmpl.Variables {
					required := ""
					if v.Required {
						required = " (required)"
					}
					fmt.Printf("      - %s%s: %s\n", v.Name, required, v.Description)
				}
			}
			fmt.Println()
		}

		return nil
	},
}
