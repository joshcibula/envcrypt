package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Root returns the root cobra command for envcrypt.
func Root() *cobra.Command {
	if rootCmd != nil {
		return rootCmd
	}

	rootCmd = &cobra.Command{
		Use:   "envcrypt",
		Short: "Lightweight .env file encryption and sharing tool",
		Long: `envcrypt encrypts and decrypts .env files using age encryption.

Use it to safely store and share environment secrets with teammates.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newLockCmd())
	rootCmd.AddCommand(newUnlockCmd())

	return rootCmd
}
