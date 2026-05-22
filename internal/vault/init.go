package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envcrypt/internal/config"
	"github.com/user/envcrypt/internal/keystore"
)

// InitOptions controls vault initialisation behaviour.
type InitOptions struct {
	// CfgDir is the directory where config and keys are stored.
	CfgDir string
	// EnvFile is the path to the plaintext .env file to encrypt on init.
	EnvFile string
	// Force overwrites an existing vault without prompting.
	Force bool
}

// Init creates a new vault: generates a key pair, saves config, and locks the
// provided env file if one is given.
func Init(opts InitOptions) (*Vault, error) {
	if err := os.MkdirAll(opts.CfgDir, 0o700); err != nil {
		return nil, fmt.Errorf("vault: init: mkdir: %w", err)
	}

	cfgPath := filepath.Join(opts.CfgDir, "config.toml")
	if !opts.Force {
		if _, err := os.Stat(cfgPath); err == nil {
			return nil, fmt.Errorf("vault: init: vault already exists at %s (use --force to overwrite)", opts.CfgDir)
		}
	}

	id, err := keystore.Generate()
	if err != nil {
		return nil, fmt.Errorf("vault: init: generate key: %w", err)
	}

	keyPath := filepath.Join(opts.CfgDir, "identity.age")
	if err := keystore.Save(id, keyPath); err != nil {
		return nil, fmt.Errorf("vault: init: save key: %w", err)
	}

	cfg := config.Default()
	cfg.IdentityFile = keyPath
	if err := config.Save(cfg, cfgPath); err != nil {
		return nil, fmt.Errorf("vault: init: save config: %w", err)
	}

	v := &Vault{Cfg: cfg, identity: id}

	if opts.EnvFile != "" {
		if err := v.Lock(opts.EnvFile); err != nil {
			return nil, fmt.Errorf("vault: init: initial lock: %w", err)
		}
	}

	return v, nil
}
