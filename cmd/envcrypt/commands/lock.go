package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/envcrypt/internal/vault"
)

func newLockCmd() *cobra.Command {
	var envPath string
	var vaultPath string
	var recipientsFile string

	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Encrypt a .env file into the vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.Open(vaultPath)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}

			if err := v.Lock(envPath, recipientsFile); err != nil {
				return fmt.Errorf("lock failed: %w", err)
			}

			fmt.Printf("Locked %s → %s\n", envPath, vaultPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&envPath, "env", "e", ".env", "Path to .env file to encrypt")
	cmd.Flags().StringVarP(&vaultPath, "vault", "v", ".env.age", "Path to vault file")
	cmd.Flags().StringVarP(&recipientsFile, "recipients", "r", "", "Optional recipients file for multi-key encryption")

	return cmd
}
