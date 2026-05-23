package vault

import (
	"fmt"
	"sort"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// ListResult holds the result of listing vault keys.
type ListResult struct {
	Keys []string
}

// Count returns the number of keys in the vault.
func (r *ListResult) Count() int {
	return len(r.Keys)
}

// List decrypts the vault and returns the list of environment variable keys
// without exposing their values.
func List(vaultPath string) (*ListResult, error) {
	cfg, err := config.LoadOrDefault(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("list: load config: %w", err)
	}

	identity, err := keystore.Load(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("list: load key: %w", err)
	}

	plaintext, err := crypto.DecryptFile(cfg.VaultFile, identity)
	if err != nil {
		return nil, fmt.Errorf("list: decrypt vault: %w", err)
	}

	pairs, err := envfile.Parse(string(plaintext))
	if err != nil {
		return nil, fmt.Errorf("list: parse env: %w", err)
	}

	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &ListResult{Keys: keys}, nil
}
