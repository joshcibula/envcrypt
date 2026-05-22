package crypto_test

import (
	"testing"

	"filippo.io/age"

	"github.com/yourorg/envcrypt/internal/crypto"
)

func generateTestIdentity(t *testing.T) *age.X25519Identity {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}
	return id
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	id := generateTestIdentity(t)
	plaintext := []byte("SECRET=hello\nDB_PASS=world")

	ciphertext, err := crypto.Encrypt(plaintext, []age.Recipient{id.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Fatal("expected non-empty ciphertext")
	}

	got, err := crypto.Decrypt(ciphertext, []age.Identity{id})
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Errorf("roundtrip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	id1 := generateTestIdentity(t)
	id2 := generateTestIdentity(t)

	ciphertext, err := crypto.Encrypt([]byte("KEY=value"), []age.Recipient{id1.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(ciphertext, []age.Identity{id2})
	if err == nil {
		t.Fatal("expected error decrypting with wrong key, got nil")
	}
}

func TestEncryptMultipleRecipients(t *testing.T) {
	id1 := generateTestIdentity(t)
	id2 := generateTestIdentity(t)
	plaintext := []byte("MULTI=true")

	ciphertext, err := crypto.Encrypt(plaintext, []age.Recipient{id1.Recipient(), id2.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	for _, id := range []*age.X25519Identity{id1, id2} {
		got, err := crypto.Decrypt(ciphertext, []age.Identity{id})
		if err != nil {
			t.Errorf("Decrypt with recipient %v failed: %v", id.Recipient(), err)
		}
		if string(got) != string(plaintext) {
			t.Errorf("mismatch: got %q, want %q", got, plaintext)
		}
	}
}
