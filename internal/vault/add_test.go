package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestAddNewKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := vault.Init(dir, envPath, false); err != nil {
		t.Fatalf("Init: %v", err)
	}

	if err := vault.Add(dir, "NEW_KEY", "new_value"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	keys, err := vault.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	found := false
	for _, k := range keys {
		if k == "NEW_KEY" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected NEW_KEY in vault, got %v", keys)
	}
}

func TestAddOverwritesExistingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=original\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := vault.Init(dir, envPath, false); err != nil {
		t.Fatalf("Init: %v", err)
	}

	if err := vault.Add(dir, "FOO", "updated"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	tmpOut := filepath.Join(dir, "out.env")
	if err := vault.Export(dir, tmpOut, false); err != nil {
		t.Fatalf("Export: %v", err)
	}

	data, err := os.ReadFile(tmpOut)
	if err != nil {
		t.Fatal(err)
	}

	if !containsLine(string(data), "FOO=updated") {
		t.Errorf("expected FOO=updated in exported env, got:\n%s", string(data))
	}
}

func TestAddMissingVault(t *testing.T) {
	dir := t.TempDir()
	err := vault.Add(dir, "KEY", "val")
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}
