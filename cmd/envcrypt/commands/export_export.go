package commands

import "github.com/spf13/cobra"

// NewExportCmdForTest exposes newExportCmd for use in tests outside this package.
func NewExportCmdForTest() *cobra.Command {
	return newExportCmd()
}

func init() {
	rootCmd.AddCommand(newExportCmd())
}
