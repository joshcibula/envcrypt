package keystore

import (
	"errors"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const (
	DefaultKeyFile = ".envcrypt_key"
)

// KeyPair holds an age X25519 identity (private key) and its recipient (public key).
type KeyPair struct {
	Identity  *age.X25519Identity
	Recipient *age.X25519Recipient
}

// Generate creates a new X25519 key pair.
func Generate() (*KeyPair, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		Identity:  identity,
		Recipient: identity.Recipient(),
	}, nil
}

// Save writes the private key to a file with restrictive permissions.
func Save(kp *KeyPair, path string) error {
	if path == "" {
		path = DefaultKeyFile
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(kp.Identity.String()+"\n"), 0o600)
}

// Load reads an age private key from the given file path.
func Load(path string) (*KeyPair, error) {
	if path == "" {
		path = DefaultKeyFile
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("key file not found: run 'envcrypt keygen' first")
		}
		return nil, err
	}
	identity, err := age.ParseX25519Identity(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		Identity:  identity,
		Recipient: identity.Recipient(),
	}, nil
}
