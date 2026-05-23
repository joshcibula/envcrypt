package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestRotateReplacesKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, "vault.age")
	keyPath := filepath.Join(dir, "key.txt")

	// Write a minimal .env file.
	if err := os.WriteFile(envPath, []byte("SECRET=hello\n"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		VaultFile: vaultPath,
		KeyFile:   keyPath,
	}

	// Initialise the vault so there is something to rotate.
	if err := vault.Init(cfg, envPath, false); err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Capture the public key before rotation.
	oldIdentity, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load old identity: %v", err)
	}
	oldPub := oldIdentity.Recipient().String()

	// Rotate.
	if err := vault.Rotate(cfg); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	// The key file must now contain a different public key.
	newIdentity, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load new identity: %v", err)
	}
	newPub := newIdentity.Recipient().String()

	if oldPub == newPub {
		t.Error("expected public key to change after rotation, but it did not")
	}

	// The vault must still be openable with the new key.
	env, err := vault.Open(cfg)
	if err != nil {
		t.Fatalf("Open after rotate: %v", err)
	}
	if env["SECRET"] != "hello" {
		t.Errorf("expected SECRET=hello, got %q", env["SECRET"])
	}
}

func TestRotateMissingVault(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		VaultFile: filepath.Join(dir, "nonexistent.age"),
		KeyFile:   filepath.Join(dir, "key.txt"),
	}

	// Generate a key so the load step does not fail first.
	id, _ := keystore.Generate()
	_ = keystore.Save(cfg.KeyFile, id)

	if err := vault.Rotate(cfg); err == nil {
		t.Error("expected error rotating missing vault, got nil")
	}
}
