package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestDeleteExistingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, "vault.age")

	_ = os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	if err := vault.Init(vaultPath, envPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := vault.Delete(vaultPath, "FOO"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	keys, err := vault.List(vaultPath)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	for _, k := range keys {
		if k == "FOO" {
			t.Error("expected FOO to be deleted but it still exists")
		}
	}

	found := false
	for _, k := range keys {
		if k == "BAZ" {
			found = true
		}
	}
	if !found {
		t.Error("expected BAZ to remain after deleting FOO")
	}
}

func TestDeleteMissingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, "vault.age")

	_ = os.WriteFile(envPath, []byte("FOO=bar\n"), 0600)

	if err := vault.Init(vaultPath, envPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	err := vault.Delete(vaultPath, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error when deleting non-existent key, got nil")
	}
}

func TestDeleteMissingVault(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "vault.age")

	err := vault.Delete(vaultPath, "FOO")
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}
