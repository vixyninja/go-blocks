/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-14 09:48:00
 * @ Modified time: 2025-11-14 11:14:02
 * @ Description: Version information
 */

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
