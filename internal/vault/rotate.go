package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envcrypt/internal/config"
	"github.com/yourusername/envcrypt/internal/crypto"
	"github.com/yourusername/envcrypt/internal/keystore"
)

// Rotate re-encrypts the vault with a freshly generated key, replacing the
// existing identity file. The old key is discarded after a successful rotation.
func Rotate(dir string) error {
	cfg, err := config.LoadOrDefault(filepath.Join(dir, ".envcrypt.toml"))
	if err != nil {
		return fmt.Errorf("rotate: load config: %w", err)
	}

	vaultPath := filepath.Join(dir, cfg.VaultFile)
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("rotate: vault file not found: %s", vaultPath)
	}

	// Load existing identity so we can decrypt the current vault.
	keyPath := filepath.Join(dir, cfg.KeyFile)
	oldIdentity, err := keystore.Load(keyPath)
	if err != nil {
		return fmt.Errorf("rotate: load existing key: %w", err)
	}

	// Decrypt vault contents with the old key.
	plaintext, err := crypto.DecryptFile(vaultPath, oldIdentity)
	if err != nil {
		return fmt.Errorf("rotate: decrypt vault: %w", err)
	}

	// Generate a new identity.
	newIdentity, err := keystore.Generate()
	if err != nil {
		return fmt.Errorf("rotate: generate new key: %w", err)
	}

	// Re-encrypt with the new recipient.
	recipient := newIdentity.Recipient()
	if err := crypto.EncryptFile(plaintext, vaultPath, recipient); err != nil {
		return fmt.Errorf("rotate: encrypt vault with new key: %w", err)
	}

	// Persist the new identity, overwriting the old one.
	if err := keystore.Save(newIdentity, keyPath); err != nil {
		return fmt.Errorf("rotate: save new key: %w", err)
	}

	return nil
}
