package vault

import (
	"fmt"

	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// Edit decrypts the vault, updates the value for the given key, and re-encrypts it.
// If the key does not exist, it is created.
func Edit(vaultPath, keyPath, key, value string) error {
	identity, err := keystore.Load(keyPath)
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	v, err := Open(vaultPath)
	if err != nil {
		return fmt.Errorf("open vault: %w", err)
	}

	plaintext, err := crypto.Decrypt(v.Ciphertext, identity)
	if err != nil {
		return fmt.Errorf("decrypt vault: %w", err)
	}

	envMap, err := envfile.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	envMap[key] = value

	updated := envfile.Serialize(envMap)

	recipient, err := identity.Recipient()
	if err != nil {
		return fmt.Errorf("derive recipient: %w", err)
	}

	ciphertext, err := crypto.Encrypt([]byte(updated), recipient)
	if err != nil {
		return fmt.Errorf("encrypt vault: %w", err)
	}

	v.Ciphertext = ciphertext
	return v.save(vaultPath)
}
