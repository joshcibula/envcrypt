package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/envcrypt/internal/config"
	"github.com/yourusername/envcrypt/internal/vault"
)

func newUnlockCmd() *cobra.Command {
	var (
		vaultPath string
		envPath   string
		keyPath   string
	)

	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Decrypt the vault and restore the .env file",
		Long: `Decrypt the encrypted vault file using your age identity key
and write the plaintext variables back to the .env file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadOrDefault("")
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			if vaultPath == "" {
				vaultPath = cfg.VaultFile
			}
			if envPath == "" {
				envPath = cfg.EnvFile
			}
			if keyPath == "" {
				keyPath = cfg.IdentityFile
			}

			if err := vault.Unlock(vaultPath, envPath, keyPath); err != nil {
				return fmt.Errorf("unlock: %w", err)
			}

			fmt.Printf("Unlocked: %s → %s\n", vaultPath, envPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&vaultPath, "vault", "", "path to encrypted vault file (default from config)")
	cmd.Flags().StringVar(&envPath, "env", "", "path to output .env file (default from config)")
	cmd.Flags().StringVar(&keyPath, "key", "", "path to age identity key file (default from config)")

	return cmd
}
