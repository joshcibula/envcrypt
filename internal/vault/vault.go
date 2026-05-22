// Package vault provides high-level operations for managing encrypted .env vaults.
// A vault ties together key management, env file parsing, and encryption.
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envcrypt/internal/config"
	"github.com/user/envcrypt/internal/crypto"
	"github.com/user/envcrypt/internal/envfile"
	"github.com/user/envcrypt/internal/keystore"
)

// Vault represents an envcrypt vault bound to a config and key.
type Vault struct {
	Cfg      *config.Config
	identity *keystore.Identity
}

// Open loads or initialises a vault using the given config directory.
func Open(cfgDir string) (*Vault, error) {
	cfg, err := config.LoadOrDefault(filepath.Join(cfgDir, "config.toml"))
	if err != nil {
		return nil, fmt.Errorf("vault: load config: %w", err)
	}

	id, err := keystore.Load(cfg.IdentityFile)
	if err != nil {
		return nil, fmt.Errorf("vault: load identity: %w", err)
	}

	return &Vault{Cfg: cfg, identity: id}, nil
}

// Lock encrypts the plaintext env file, writing ciphertext to the configured path.
func (v *Vault) Lock(plainPath string) error {
	if _, err := os.Stat(plainPath); err != nil {
		return fmt.Errorf("vault: lock: source file not found: %w", err)
	}

	recipient, err := v.identity.Recipient()
	if err != nil {
		return fmt.Errorf("vault: lock: get recipient: %w", err)
	}

	if err := crypto.EncryptFile(plainPath, v.Cfg.EncryptedFile, recipient); err != nil {
		return fmt.Errorf("vault: lock: %w", err)
	}
	return nil
}

// Unlock decrypts the vault's encrypted file and returns the parsed env map.
func (v *Vault) Unlock() (map[string]string, error) {
	tmp, err := os.CreateTemp("", "envcrypt-unlock-*.env")
	if err != nil {
		return nil, fmt.Errorf("vault: unlock: create temp: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	if err := crypto.DecryptFile(v.Cfg.EncryptedFile, tmpPath, v.identity); err != nil {
		return nil, fmt.Errorf("vault: unlock: %w", err)
	}

	vars, err := envfile.Parse(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("vault: unlock: parse env: %w", err)
	}
	return vars, nil
}
