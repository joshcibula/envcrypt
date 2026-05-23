package commands

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var vaultPath string

	cmd := &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete a key from the encrypted vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyName := args[0]

			if err := vault.Delete(vaultPath, keyName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			fmt.Printf("Deleted key %q from vault.\n", keyName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&vaultPath, "vault", "v", "secrets.age", "Path to the encrypted vault file")

	return cmd
}

func init() {
	rootCmd.AddCommand(newDeleteCmd())
}
