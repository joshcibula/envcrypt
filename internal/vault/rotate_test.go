package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envcrypt/internal/keystore"
	"github.com/yourusername/envcrypt/internal/vault"
)

func TestRotateReplacesKey(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("SECRET=hello\nTOKEN=world\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, envFile, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	cfgPath := filepath.Join(dir, ".envcrypt.toml")
	keyPath := filepath.Join(dir, ".envcrypt.key")
	_ = cfgPath

	// Capture the public key before rotation.
	identityBefore, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load key before rotate: %v", err)
	}
	pubBefore := identityBefore.Recipient().String()

	if err := vault.Rotate(dir); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	// The key on disk must have changed.
	identityAfter, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load key after rotate: %v", err)
	}
	pubAfter := identityAfter.Recipient().String()

	if pubBefore == pubAfter {
		t.Error("expected public key to change after rotation, but it did not")
	}

	// The vault must still be unlockable with the new key.
	outDir := t.TempDir()
	outPath := filepath.Join(outDir, ".env")
	if err := vault.Unlock(dir, outPath); err != nil {
		t.Fatalf("unlock after rotate: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read unlocked env: %v", err)
	}
	if got := string(data); got != "SECRET=hello\nTOKEN=world\n" {
		t.Errorf("unexpected env contents after rotate: %q", got)
	}
}

func TestRotateMissingVault(t *testing.T) {
	dir := t.TempDir()
	err := vault.Rotate(dir)
	if err == nil {
		t.Fatal("expected error rotating non-existent vault, got nil")
	}
}
