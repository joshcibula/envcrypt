package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
)

// UnlockOptions holds optional overrides for the Unlock operation.
type UnlockOptions struct {
	// OutputFile overrides the destination .env path from config.
	OutputFile string
}

// Unlock decrypts the vault file and writes the plaintext .env to disk.
func Unlock(cfg *config.Config, opts UnlockOptions) error {
	identity, err := keystore.Load(cfg.IdentityFile)
	if err != nil {
		return fmt.Errorf("load identity: %w", err)
	}

	vaultPath := cfg.VaultFile
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault file not found: %s", vaultPath)
	}

	plaintext, err := crypto.DecryptFile(vaultPath, identity)
	if err != nil {
		return fmt.Errorf("decrypt vault: %w", err)
	}

	dest := cfg.EnvFile
	if opts.OutputFile != "" {
		dest = opts.OutputFile
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.WriteFile(dest, plaintext, 0600); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	return nil
}
