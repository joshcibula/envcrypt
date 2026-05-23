package vault

import (
	"fmt"

	"filippo.io/age"
	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// Add inserts or updates a key-value pair in the vault and re-encrypts it.
// If the key already exists its value is overwritten.
func Add(vaultPath, key, value string) error {
	cfg, err := config.LoadOrDefault(vaultPath)
	if err != nil {
		return fmt.Errorf("add: load config: %w", err)
	}

	identity, err := keystore.Load(cfg.KeyPath)
	if err != nil {
		return fmt.Errorf("add: load key: %w", err)
	}

	plaintext, err := crypto.DecryptFile(cfg.VaultFile, identity)
	if err != nil {
		return fmt.Errorf("add: decrypt vault: %w", err)
	}

	envMap, err := envfile.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("add: parse env: %w", err)
	}

	envMap[key] = value

	serialised := envfile.Serialize(envMap)

	recipient, ok := identity.Recipient().(*age.X25519Recipient)
	if !ok {
		return fmt.Errorf("add: unexpected recipient type")
	}

	if err := crypto.EncryptFile(cfg.VaultFile, []byte(serialised), recipient); err != nil {
		return fmt.Errorf("add: encrypt vault: %w", err)
	}

	return nil
}
