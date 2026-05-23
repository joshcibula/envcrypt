package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestGetExistingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".env.enc")
	keyPath := filepath.Join(dir, "key.age")

	_ = os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	if err := Init(vaultPath, keyPath, envPath, false); err != nil {
		t.Fatalf("Init: %v", err)
	}

	val, err := vault.Get(vaultPath, keyPath, "FOO")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "bar" {
		t.Errorf("expected \"bar\", got %q", val)
	}
}

func TestGetMissingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".env.enc")
	keyPath := filepath.Join(dir, "key.age")

	_ = os.WriteFile(envPath, []byte("FOO=bar\n"), 0600)

	if err := Init(vaultPath, keyPath, envPath, false); err != nil {
		t.Fatalf("Init: %v", err)
	}

	_, err := vault.Get(vaultPath, keyPath, "DOES_NOT_EXIST")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestGetMissingVault(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.age")

	_, err := vault.Get(filepath.Join(dir, ".env.enc"), keyPath, "FOO")
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}
