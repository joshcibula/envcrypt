package vault

import (
	"fmt"

	"filippo.io/age"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// Delete removes a key-value pair from the encrypted vault.
// Returns an error if the vault does not exist or the key is not present.
func Delete(vaultPath, keyName string) error {
	ident, err := keystore.Load(vaultPath + ".key")
	if err != nil {
		return fmt.Errorf("delete: load key: %w", err)
	}

	envBytes, err := crypto.DecryptFile(vaultPath, ident)
	if err != nil {
		return fmt.Errorf("delete: decrypt vault: %w", err)
	}

	parsed, err := envfile.Parse(string(envBytes))
	if err != nil {
		return fmt.Errorf("delete: parse env: %w", err)
	}

	if _, exists := parsed[keyName]; !exists {
		return fmt.Errorf("delete: key %q not found in vault", keyName)
	}

	delete(parsed, keyName)

	updated := envfile.Serialize(parsed)

	recipient, ok := ident.Recipient().(*age.X25519Recipient)
	if !ok {
		return fmt.Errorf("delete: unexpected recipient type")
	}

	if err := crypto.EncryptFile([]byte(updated), vaultPath, recipient); err != nil {
		return fmt.Errorf("delete: re-encrypt vault: %w", err)
	}

	return nil
}
