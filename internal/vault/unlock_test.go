package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envcrypt/internal/vault"
)

func TestUnlockRoundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".env.vault")
	keyPath := filepath.Join(dir, "identity.txt")

	// Write a temp env file and init the vault
	envContent := "DB_HOST=localhost\nDB_PORT=5432\nSECRET=supersecret\n"
	if err := os.WriteFile(envPath, []byte(envContent), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Remove the plaintext env file
	if err := os.Remove(envPath); err != nil {
		t.Fatalf("remove env: %v", err)
	}

	// Unlock should recreate it
	if err := vault.Unlock(vaultPath, envPath, keyPath); err != nil {
		t.Fatalf("unlock: %v", err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("read env: %v", err)
	}

	got := string(data)
	for _, line := range []string{"DB_HOST=localhost", "DB_PORT=5432", "SECRET=supersecret"} {
		if !containsLine(got, line) {
			t.Errorf("expected line %q in unlocked env, got:\n%s", line, got)
		}
	}
}

func TestUnlockMissingVault(t *testing.T) {
	dir := t.TempDir()
	err := vault.Unlock(
		filepath.Join(dir, "nonexistent.vault"),
		filepath.Join(dir, ".env"),
		filepath.Join(dir, "identity.txt"),
	)
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}

func TestUnlockWrongKey(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath := filepath.Join(dir, ".env.vault")
	keyPath := filepath.Join(dir, "identity.txt")
	wrongKeyPath := filepath.Join(dir, "wrong_identity.txt")

	if err := os.WriteFile(envPath, []byte("FOO=bar\n"), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Generate a different key
	if err := vault.Init(
		filepath.Join(dir, "other.env"),
		filepath.Join(dir, "other.vault"),
		wrongKeyPath,
		false,
	); err == nil {
		// ignore error from missing other.env, just need the key generated
	}

	// Write a dummy wrong key
	if err := os.WriteFile(wrongKeyPath, []byte("AGE-SECRET-KEY-1QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ\n"), 0o600); err != nil {
		t.Fatalf("write wrong key: %v", err)
	}

	err := vault.Unlock(vaultPath, envPath, wrongKeyPath)
	if err == nil {
		t.Fatal("expected error when unlocking with wrong key")
	}
}

func containsLine(s, line string) bool {
	for _, l := range splitLines(s) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
