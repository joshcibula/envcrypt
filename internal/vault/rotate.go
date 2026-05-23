package vault

import (
	"fmt"
	"os"

	"github.com/user/envcrypt/internal/crypto"
	"github.com/user/envcrypt/internal/keystore"
)

// Rotate generates a new key pair, re-encrypts the vault with the new key,
// and replaces the old key on disk.
func Rotate(vaultPath, keyPath string) error {
	// Load existing identity to decrypt current vault
	oldIdentity, err := keystore.Load(keyPath)
	if err != nil {
		return fmt.Errorf("rotate: load existing key: %w", err)
	}

	// Decrypt vault contents with old key
	plaintext, err := crypto.DecryptFile(vaultPath, oldIdentity)
	if err != nil {
		return fmt.Errorf("rotate: decrypt vault: %w", err)
	}

	// Generate new identity
	newIdentity, err := keystore.Generate()
	if err != nil {
		return fmt.Errorf("rotate: generate new key: %w", err)
	}

	// Write new vault encrypted with new key
	tmpVault := vaultPath + ".tmp"
	recipient := newIdentity.Recipient()
	if err := crypto.EncryptFile(tmpVault, plaintext, recipient); err != nil {
		_ = os.Remove(tmpVault)
		return fmt.Errorf("rotate: encrypt vault with new key: %w", err)
	}

	// Atomically replace old vault
	if err := os.Rename(tmpVault, vaultPath); err != nil {
		_ = os.Remove(tmpVault)
		return fmt.Errorf("rotate: replace vault file: %w", err)
	}

	// Persist new key, overwriting old one
	if err := keystore.Save(newIdentity, keyPath); err != nil {
		return fmt.Errorf("rotate: save new key: %w", err)
	}

	return nil
}
