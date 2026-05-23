package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestListKeys(t *testing.T) {
	dir := t.TempDir()
	envContent := "FOO=bar\nBAZ=qux\nALPHA=1\n"
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, envFile, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	result, err := vault.List(dir)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(result.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result.Keys))
	}

	// Keys should be sorted
	expected := []string{"ALPHA", "BAZ", "FOO"}
	for i, k := range result.Keys {
		if k != expected[i] {
			t.Errorf("key[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestListMissingVault(t *testing.T) {
	dir := t.TempDir()
	_, err := vault.List(dir)
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}

func TestListDoesNotExposeValues(t *testing.T) {
	dir := t.TempDir()
	envContent := "SECRET=supersecret\n"
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, envFile, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	result, err := vault.List(dir)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(result.Keys) != 1 || result.Keys[0] != "SECRET" {
		t.Errorf("unexpected keys: %v", result.Keys)
	}
}
