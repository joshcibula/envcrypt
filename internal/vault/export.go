package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envcrypt/internal/config"
	"github.com/yourusername/envcrypt/internal/crypto"
	"github.com/yourusername/envcrypt/internal/envfile"
	"github.com/yourusername/envcrypt/internal/keystore"
)

// ExportOptions controls how the vault is exported.
type ExportOptions struct {
	// OutputPath is the destination file. If empty, defaults to ".env".
	OutputPath string
	// Overwrite allows replacing an existing output file.
	Overwrite bool
}

// Export decrypts the vault and writes the plaintext env file to disk.
// It differs from Unlock in that it targets an arbitrary output path and
// does not modify vault state.
func Export(vaultDir string, opts ExportOptions) error {
	cfg, err := config.LoadOrDefault(filepath.Join(vaultDir, ".envcrypt.toml"))
	if err != nil {
		return fmt.Errorf("export: load config: %w", err)
	}

	identity, err := keystore.Load(cfg.KeyPath)
	if err != nil {
		return fmt.Errorf("export: load key: %w", err)
	}

	vaultPath := filepath.Join(vaultDir, cfg.VaultFile)
	plaintext, err := crypto.DecryptFile(vaultPath, identity)
	if err != nil {
		return fmt.Errorf("export: decrypt vault: %w", err)
	}

	envVars, err := envfile.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("export: parse env: %w", err)
	}

	outPath := opts.OutputPath
	if outPath == "" {
		outPath = ".env"
	}

	if !opts.Overwrite {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf("export: output file %q already exists (use --force to overwrite)", outPath)
		}
	}

	contents := envfile.Serialize(envVars)
	if err := os.WriteFile(outPath, []byte(contents), 0600); err != nil {
		return fmt.Errorf("export: write output file: %w", err)
	}

	return nil
}
