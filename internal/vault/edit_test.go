package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestEditUpdatesExistingKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, "vault.age")
	keyPath := filepath.Join(dir, "key.txt")

	_ = os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := vault.Edit(vaultPath, keyPath, "FOO", "updated"); err != nil {
		t.Fatalf("edit: %v", err)
	}

	identity, _ := keystore.Load(keyPath)
	v, _ := vault.Open(vaultPath)
	plaintext, _ := crypto.Decrypt(v.Ciphertext, identity)
	envMap, _ := envfile.Parse(string(plaintext))

	if envMap["FOO"] != "updated" {
		t.Errorf("expected FOO=updated, got FOO=%s", envMap["FOO"])
	}
	if envMap["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got BAZ=%s", envMap["BAZ"])
	}
}

func TestEditAddsNewKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, "vault.age")
	keyPath := filepath.Join(dir, "key.txt")

	_ = os.WriteFile(envPath, []byte("FOO=bar\n"), 0600)

	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := vault.Edit(vaultPath, keyPath, "NEW_KEY", "new_value"); err != nil {
		t.Fatalf("edit: %v", err)
	}

	identity, _ := keystore.Load(keyPath)
	v, _ := vault.Open(vaultPath)
	plaintext, _ := crypto.Decrypt(v.Ciphertext, identity)
	envMap, _ := envfile.Parse(string(plaintext))

	if envMap["NEW_KEY"] != "new_value" {
		t.Errorf("expected NEW_KEY=new_value, got NEW_KEY=%s", envMap["NEW_KEY"])
	}
}

func TestEditMissingVault(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.txt")
	_ = keystore.Generate(keyPath)

	err := vault.Edit(filepath.Join(dir, "nonexistent.age"), keyPath, "K", "V")
	if err == nil {
		t.Fatal("expected error for missing vault")
	}
}
