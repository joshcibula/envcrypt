package commands

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

// NewListCmdForTest returns a list command instance suitable for use in tests.
func NewListCmdForTest() *cobra.Command {
	var vaultDir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List environment variable keys stored in the vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := vault.List(vaultDir)
			if err != nil {
				return fmt.Errorf("list: %w", err)
			}
			if len(result.Keys) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no keys found)")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(result.Keys, "\n"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&vaultDir, "dir", "d", ".", "directory containing the vault")
	return cmd
}
