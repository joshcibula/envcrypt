package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourorg/envcrypt/internal/vault"
)

func newUnlockCmd() *cobra.Command {
	var envPath string
	var vaultPath string
	var keyPath string

	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Decrypt the vault into a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := vault.Open(vaultPath)
			if err != nil {
				return fmt.Errorf("open vault: %w", err)
			}

			if err := v.Unlock(envPath, keyPath); err != nil {
				return fmt.Errorf("unlock failed: %w", err)
			}

			fmt.Printf("Unlocked %s → %s\n", vaultPath, envPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&envPath, "env", "e", ".env", "Output path for decrypted .env file")
	cmd.Flags().StringVarP(&vaultPath, "vault", "v", ".env.age", "Path to vault file")
	cmd.Flags().StringVarP(&keyPath, "key", "k", "", "Path to age identity key file (default: ~/.config/envcrypt/key.age)")

	return cmd
}
