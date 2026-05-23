package commands

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"
)

// newRotateCmd returns the "rotate" sub-command which generates a new
// encryption key and re-encrypts the vault without touching the plaintext.
func newRotateCmd() *cobra.Command {
	var cfgPath string

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate the encryption key and re-encrypt the vault",
		Long: `Generate a fresh age identity, decrypt the current vault with the
old key, re-encrypt it with the new key, and replace the key file on disk.

The plaintext .env file is never written to disk during rotation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadOrDefault(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if err := vault.Rotate(cfg); err != nil {
				return fmt.Errorf("rotate: %w", err)
			}

			fmt.Println("Key rotated and vault re-encrypted successfully.")
			return nil
		},
	}

	cmd.Flags().StringVar(&cfgPath, "config", "", "path to envcrypt config file (default: .envcrypt.yaml)")
	return cmd
}
