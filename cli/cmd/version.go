package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vixyninja/go-blocks/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	Run: func(cmd *cobra.Command, _ []string) {
		j, _ := cmd.Flags().GetBool("json")
		if j {
			fmt.Println(version.Get().JSON())
			return
		}
		fmt.Print(version.Get().String())
	},
}
