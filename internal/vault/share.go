package vault

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/sharing"
	age "filippo.io/age"
)

// ShareOptions holds parameters for re-encrypting a vault for additional recipients.
type ShareOptions struct {
	// VaultPath is the path to the encrypted vault file.
	VaultPath string
	// RecipientsFile is the path to the file listing recipient public keys.
	RecipientsFile string
	// OutputPath is where the re-encrypted vault will be written.
	// If empty, VaultPath is overwritten.
	OutputPath string
}

// Share decrypts the vault using the local identity and re-encrypts it for
// the given recipients plus the local identity (so the owner retains access).
func Share(cfg *config.Config, opts ShareOptions) error {
	identity, err := keystore.Load(cfg.IdentityFile)
	if err != nil {
		return fmt.Errorf("load identity: %w", err)
	}

	recipients, err := sharing.LoadRecipientsFile(opts.RecipientsFile)
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}

	ageRecipients, err := sharing.ToAgeRecipients(recipients)
	if err != nil {
		return fmt.Errorf("parse recipient keys: %w", err)
	}

	// Always include the local identity so the owner can still decrypt.
	ownRecipient := identity.Recipient()
	allRecipients := append([]age.Recipient{ownRecipient}, ageRecipients...)

	vaultPath := opts.VaultPath
	if vaultPath == "" {
		vaultPath = cfg.VaultFile
	}

	plaintext, err := crypto.DecryptFile(vaultPath, identity)
	if err != nil {
		return fmt.Errorf("decrypt vault: %w", err)
	}

	out := opts.OutputPath
	if out == "" {
		out = vaultPath
	}

	tmp := out + ".tmp"
	if err := os.WriteFile(tmp, plaintext, 0600); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	defer os.Remove(tmp)

	if err := crypto.EncryptFile(tmp, out, allRecipients...); err != nil {
		return fmt.Errorf("re-encrypt vault: %w", err)
	}

	return nil
}
