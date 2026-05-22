package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envcrypt/internal/config"
	"github.com/yourusername/envcrypt/internal/crypto"
	"github.com/yourusername/envcrypt/internal/sharing"
)

func newShareCmd() *cobra.Command {
	var recipientsFile string
	var vaultFile string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "share",
		Short: "Re-encrypt the vault for additional recipients",
		Long: `Re-encrypt the current vault so that additional recipients
can decrypt it using their own age private keys.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadOrDefault(".envcrypt.toml")
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if vaultFile == "" {
				vaultFile = cfg.VaultFile
			}
			if outputFile == "" {
				outputFile = vaultFile
			}

			recipients, err := sharing.LoadRecipientsFile(recipientsFile)
			if err != nil {
				return fmt.Errorf("loading recipients: %w", err)
			}

			ageRecipients, err := sharing.ToAgeRecipients(recipients)
			if err != nil {
				return fmt.Errorf("parsing recipient keys: %w", err)
			}

			ciphertext, err := os.ReadFile(vaultFile)
			if err != nil {
				return fmt.Errorf("reading vault %q: %w", vaultFile, err)
			}

			identity, err := cfg.LoadIdentity()
			if err != nil {
				return fmt.Errorf("loading identity: %w", err)
			}

			plaintext, err := crypto.Decrypt(ciphertext, identity)
			if err != nil {
				return fmt.Errorf("decrypting vault: %w", err)
			}

			newCiphertext, err := crypto.Encrypt(plaintext, ageRecipients)
			if err != nil {
				return fmt.Errorf("re-encrypting for recipients: %w", err)
			}

			if err := os.WriteFile(outputFile, newCiphertext, 0o600); err != nil {
				return fmt.Errorf("writing shared vault %q: %w", outputFile, err)
			}

			fmt.Printf("Vault re-encrypted for %d recipient(s) → %s\n", len(recipients), outputFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&recipientsFile, "recipients", "r", "recipients.txt", "File containing age public keys (one per line)")
	cmd.Flags().StringVarP(&vaultFile, "vault", "v", "", "Vault file to re-encrypt (defaults to config value)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for shared vault (defaults to vault file)")

	return cmd
}
