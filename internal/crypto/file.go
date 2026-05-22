package crypto

import (
	"fmt"
	"os"

	"filippo.io/age"
)

const encryptedFilePerm = 0600

// EncryptFile reads plaintext from srcPath, encrypts it, and writes the
// ciphertext to dstPath.
func EncryptFile(srcPath, dstPath string, recipients []age.Recipient) error {
	plaintext, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("crypto: read source file %q: %w", srcPath, err)
	}
	ciphertext, err := Encrypt(plaintext, recipients)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dstPath, ciphertext, encryptedFilePerm); err != nil {
		return fmt.Errorf("crypto: write encrypted file %q: %w", dstPath, err)
	}
	return nil
}

// DecryptFile reads ciphertext from srcPath, decrypts it, and writes the
// plaintext to dstPath.
func DecryptFile(srcPath, dstPath string, identities []age.Identity) error {
	ciphertext, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("crypto: read encrypted file %q: %w", srcPath, err)
	}
	plaintext, err := Decrypt(ciphertext, identities)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dstPath, plaintext, encryptedFilePerm); err != nil {
		return fmt.Errorf("crypto: write decrypted file %q: %w", dstPath, err)
	}
	return nil
}
