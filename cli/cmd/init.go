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
	// Add commands to root
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(scaffoldCmd)

	// Add subcommands
	scaffoldCmd.AddCommand(scaffoldNewCmd)
	scaffoldCmd.AddCommand(scaffoldListCmd)

	// Flags for root
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Flags for version
	versionCmd.Flags().BoolP("json", "j", false, "Output version in JSON format")

	// Flags for scaffold new
	scaffoldNewCmd.Flags().StringP("template", "t", "", "Template name to use")
	scaffoldNewCmd.Flags().StringP("output", "o", "", "Output directory (default: ./<project-name>)")
	scaffoldNewCmd.Flags().StringP("module", "m", "", "Go module path (e.g., github.com/user/project)")
	scaffoldNewCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	scaffoldNewCmd.Flags().StringArray("var", []string{}, "Template variables in format key=value (can be used multiple times)")

}
