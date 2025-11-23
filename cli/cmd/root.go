package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vixyninja/go-blocks/logx"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Internal development CLI",
	Long: `cli is the internal command-line tool that provides helper commands
for working with this project. It includes tools for generation,
formatting, utilities, migrations, and other developer tasks.`,
	Example: `
  # Show all commands
  cli --help
`,
}

// Execute runs the root Cobra command.
func Execute() {
	fmt.Println(logo)
	err := rootCmd.Execute()
	if err != nil {
		logx.NewDefaultLogger().Fatal(rootCmd.Context(), "Error: "+err.Error())
	}
}
