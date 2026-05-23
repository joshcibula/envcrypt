package vault

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// Get decrypts the vault and returns the value for the given key.
// Returns an error if the vault does not exist, the key is not found,
// or decryption fails.
func Get(vaultPath, keyPath, key string) (string, error) {
	identity, err := keystore.Load(keyPath)
	if err != nil {
		return "", fmt.Errorf("load key: %w", err)
	}

	plaintext, err := crypto.DecryptFile(vaultPath, identity)
	if err != nil {
		return "", fmt.Errorf("decrypt vault: %w", err)
	}

	envVars, err := envfile.Parse(string(plaintext))
	if err != nil {
		return "", fmt.Errorf("parse env: %w", err)
	}

	val, ok := envVars[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in vault", key)
	}

	return val, nil
}
