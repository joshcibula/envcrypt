package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envcrypt/internal/crypto"
	"github.com/yourusername/envcrypt/internal/envfile"
	"github.com/yourusername/envcrypt/internal/keystore"
)

// Unlock decrypts the vault file and writes the plaintext .env file to disk.
// It loads the identity from the keystore, decrypts the vault, and writes
// the result to the configured env file path.
func Unlock(vaultPath, envPath, keyPath string) error {
	identity, err := keystore.Load(keyPath)
	if err != nil {
		return fmt.Errorf("unlock: load identity: %w", err)
	}

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("unlock: vault file not found: %s", vaultPath)
	}

	envVars, err := crypto.DecryptFile(vaultPath, identity)
	if err != nil {
		return fmt.Errorf("unlock: decrypt vault: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(envPath), 0o755); err != nil {
		return fmt.Errorf("unlock: create env dir: %w", err)
	}

	contents := envfile.Serialize(envVars)
	if err := os.WriteFile(envPath, []byte(contents), 0o600); err != nil {
		return fmt.Errorf("unlock: write env file: %w", err)
	}

	return nil
}
