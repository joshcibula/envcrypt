package crypto

import (
	"bytes"
	"fmt"
	"io"

	"filippo.io/age"
)

// Encrypt encrypts plaintext using the provided age recipients.
func Encrypt(plaintext []byte, recipients []age.Recipient) ([]byte, error) {
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipients...)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create age encryptor: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("crypto: failed to write plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("crypto: failed to finalize encryption: %w", err)
	}
	return buf.Bytes(), nil
}

// Decrypt decrypts ciphertext using the provided age identities.
func Decrypt(ciphertext []byte, identities []age.Identity) ([]byte, error) {
	r, err := age.Decrypt(bytes.NewReader(ciphertext), identities...)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create age decryptor: %w", err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to read decrypted data: %w", err)
	}
	return plaintext, nil
}
