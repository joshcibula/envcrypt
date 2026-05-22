// Package crypto provides encryption and decryption utilities for envcrypt
// using the age encryption library (filippo.io/age).
//
// It exposes two layers of API:
//
//   - Low-level byte-slice operations via [Encrypt] and [Decrypt], which
//     accept age recipients/identities and operate on in-memory data.
//
//   - File-level helpers via [EncryptFile] and [DecryptFile], which read from
//     and write to disk paths, suitable for use in CLI commands.
//
// Example usage:
//
//	id, _ := age.GenerateX25519Identity()
//	ciphertext, _ := crypto.Encrypt([]byte("SECRET=42"), []age.Recipient{id.Recipient()})
//	plaintext, _ := crypto.Decrypt(ciphertext, []age.Identity{id})
package crypto
