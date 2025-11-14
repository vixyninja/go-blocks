/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-14 09:45:16
 * @ Modified time: 2025-11-14 11:05:57
 * @ Description: Root command
 */

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-gen",
	Short: "cli-gen",
	Long:  `cli-gen`,
}

// Execute runs the root Cobra command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
