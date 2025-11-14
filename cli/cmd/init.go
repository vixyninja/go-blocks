/**
 * @ Author: vixyninja
 * @ Create Time: 2025-11-14 10:32:28
 * @ Modified time: 2025-11-14 11:36:29
 * @ Description: Init entry point
 * @ This file is specifically for registering commands
 */

package cmd

var (
	logo string = ` 
 ██████╗  ██████╗       ██████╗ ██╗      ██████╗  ██████╗██╗  ██╗███████╗
██╔════╝ ██╔═══██╗      ██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝██╔════╝
██║  ███╗██║   ██║█████╗██████╔╝██║     ██║   ██║██║     █████╔╝ ███████╗
██║   ██║██║   ██║╚════╝██╔══██╗██║     ██║   ██║██║     ██╔═██╗ ╚════██║
╚██████╔╝╚██████╔╝      ██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗███████║
 ╚═════╝  ╚═════╝       ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝
                                                                         `
)

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("json", "j", false, "Output version in JSON format")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
