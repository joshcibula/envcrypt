package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/yourorg/envcrypt/internal/crypto"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envcrypt-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestEncryptDecryptFile(t *testing.T) {
	id := generateTestIdentity(t)
	dir := t.TempDir()

	src := writeTempFile(t, "API_KEY=secret123\nDEBUG=false")
	enc := filepath.Join(dir, "env.age")
	dec := filepath.Join(dir, "env.decrypted")

	if err := crypto.EncryptFile(src, enc, []age.Recipient{id.Recipient()}); err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	encData, _ := os.ReadFile(enc)
	srcData, _ := os.ReadFile(src)
	if string(encData) == string(srcData) {
		t.Error("encrypted file should differ from source")
	}

	if err := crypto.DecryptFile(enc, dec, []age.Identity{id}); err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}

	decData, _ := os.ReadFile(dec)
	if string(decData) != string(srcData) {
		t.Errorf("decrypted content mismatch: got %q, want %q", decData, srcData)
	}
}

func TestEncryptFileMissing(t *testing.T) {
	id := generateTestIdentity(t)
	err := crypto.EncryptFile("/nonexistent/path.env", "/tmp/out.age", []age.Recipient{id.Recipient()})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestDecryptFileMissing(t *testing.T) {
	id := generateTestIdentity(t)
	err := crypto.DecryptFile("/nonexistent/path.age", "/tmp/out.env", []age.Identity{id})
	if err == nil {
		t.Fatal("expected error for missing encrypted file")
	}
}
