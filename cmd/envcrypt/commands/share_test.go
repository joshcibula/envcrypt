package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/yourusername/envcrypt/internal/crypto"
	"github.com/yourusername/envcrypt/internal/keystore"
)

func writeTempVault(t *testing.T, dir string, plaintext []byte, recipient age.Recipient) string {
	t.Helper()
	ciphertext, err := crypto.Encrypt(plaintext, []age.Recipient{recipient})
	if err != nil {
		t.Fatalf("encrypt vault: %v", err)
	}
	p := filepath.Join(dir, "env.age")
	if err := os.WriteFile(p, ciphertext, 0o600); err != nil {
		t.Fatalf("write vault: %v", err)
	}
	return p
}

func TestShareReEncryptsForRecipients(t *testing.T) {
	dir := t.TempDir()

	// Original identity that locked the vault.
	origIdentity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	plaintext := []byte("SECRET=hello\nTOKEN=world\n")
	vaultPath := writeTempVault(t, dir, plaintext, origIdentity.Recipient())

	// New recipient who should gain access.
	newIdentity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate new identity: %v", err)
	}

	// Write a recipients file with both original and new public keys.
	recipientsPath := filepath.Join(dir, "recipients.txt")
	recipientsContent := origIdentity.Recipient().String() + "\n" + newIdentity.Recipient().String() + "\n"
	if err := os.WriteFile(recipientsPath, []byte(recipientsContent), 0o600); err != nil {
		t.Fatalf("write recipients file: %v", err)
	}

	outputPath := filepath.Join(dir, "shared.age")

	// Re-encrypt: decrypt with origIdentity, re-encrypt for both recipients.
	ciphertext, err := os.ReadFile(vaultPath)
	if err != nil {
		t.Fatalf("read vault: %v", err)
	}
	decrypted, err := crypto.Decrypt(ciphertext, origIdentity)
	if err != nil {
		t.Fatalf("decrypt with orig identity: %v", err)
	}

	recipients := []age.Recipient{origIdentity.Recipient(), newIdentity.Recipient()}
	newCiphertext, err := crypto.Encrypt(decrypted, recipients)
	if err != nil {
		t.Fatalf("re-encrypt: %v", err)
	}
	if err := os.WriteFile(outputPath, newCiphertext, 0o600); err != nil {
		t.Fatalf("write shared vault: %v", err)
	}

	// Verify the new identity can decrypt the shared vault.
	sharedCiphertext, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read shared vault: %v", err)
	}
	got, err := crypto.Decrypt(sharedCiphertext, newIdentity)
	if err != nil {
		t.Fatalf("new identity cannot decrypt shared vault: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("decrypted content mismatch: got %q, want %q", got, plaintext)
	}
}
