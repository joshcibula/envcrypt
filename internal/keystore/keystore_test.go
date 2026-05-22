package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envcrypt/internal/keystore"
)

func TestGenerate(t *testing.T) {
	kp, err := keystore.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if kp.Identity == nil {
		t.Fatal("expected non-nil Identity")
	}
	if kp.Recipient == nil {
		t.Fatal("expected non-nil Recipient")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test.key")

	kp, err := keystore.Generate()
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if err := keystore.Save(kp, keyPath); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}

	loaded, err := keystore.Load(keyPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Identity.String() != kp.Identity.String() {
		t.Errorf("loaded identity mismatch: got %s, want %s", loaded.Identity, kp.Identity)
	}
	if loaded.Recipient.String() != kp.Recipient.String() {
		t.Errorf("loaded recipient mismatch: got %s, want %s", loaded.Recipient, kp.Recipient)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := keystore.Load("/nonexistent/path/key.txt")
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
}
