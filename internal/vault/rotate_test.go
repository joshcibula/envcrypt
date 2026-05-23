package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/internal/keystore"
	"github.com/user/envcrypt/internal/vault"
)

func TestRotateReplacesKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".env.age")
	keyPath := filepath.Join(dir, "key.txt")

	// Write a sample .env file and initialise vault
	if err := os.WriteFile(envPath, []byte("SECRET=hello\nTOKEN=world\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Capture original public key
	origIdentity, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load original key: %v", err)
	}
	origPub := origIdentity.Recipient().String()

	// Rotate
	if err := vault.Rotate(vaultPath, keyPath); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	// New key should differ from original
	newIdentity, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("load new key: %v", err)
	}
	newPub := newIdentity.Recipient().String()
	if origPub == newPub {
		t.Error("expected new public key after rotate, got same key")
	}

	// Vault should still be decryptable with new key
	if _, err := vault.Open(vaultPath, keyPath); err != nil {
		t.Fatalf("open after rotate: %v", err)
	}
}

func TestRotateMissingVault(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "missing.age")
	keyPath := filepath.Join(dir, "key.txt")

	// Generate a key so Load succeeds but vault is absent
	identity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if err := keystore.Save(identity, keyPath); err != nil {
		t.Fatalf("save key: %v", err)
	}

	if err := vault.Rotate(vaultPath, keyPath); err == nil {
		t.Fatal("expected error rotating missing vault, got nil")
	}
}
