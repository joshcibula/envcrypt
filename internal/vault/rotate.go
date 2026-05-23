package vault

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// Rotate generates a new key pair, re-encrypts the vault with the new key,
// and replaces the old key on disk. The old identity is used to decrypt the
// current vault before re-encrypting with the freshly generated identity.
func Rotate(cfg *config.Config) error {
	// Load the existing identity so we can decrypt the current vault.
	oldIdentity, err := keystore.Load(cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("rotate: load old key: %w", err)
	}

	// Decrypt the current vault into a temporary buffer.
	plaintext, err := crypto.DecryptFile(cfg.VaultFile, oldIdentity)
	if err != nil {
		return fmt.Errorf("rotate: decrypt vault: %w", err)
	}

	// Generate a fresh identity.
	newIdentity, err := keystore.Generate()
	if err != nil {
		return fmt.Errorf("rotate: generate new key: %w", err)
	}

	// Re-encrypt the plaintext with the new identity's public key.
	recipient := newIdentity.Recipient()
	if err := crypto.EncryptFile(cfg.VaultFile, plaintext, recipient); err != nil {
		return fmt.Errorf("rotate: re-encrypt vault: %w", err)
	}

	// Persist the new identity, overwriting the old key file.
	if err := keystore.Save(cfg.KeyFile, newIdentity); err != nil {
		return fmt.Errorf("rotate: save new key: %w", err)
	}

	return nil
}
