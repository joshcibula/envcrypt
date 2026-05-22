// Package vault provides a high-level abstraction over envcrypt's core
// primitives — key management, env file parsing, and age encryption.
//
// A Vault ties together a Config (paths, options) and an age Identity
// (public/private key pair) to offer two primary operations:
//
//   - Lock: encrypt a plaintext .env file and persist the ciphertext.
//   - Unlock: decrypt the ciphertext and return the parsed key/value map.
//
// # Initialisation
//
// Use [Init] to create a brand-new vault in a given directory. Init generates
// a fresh age identity, writes it alongside a default config, and optionally
// performs an initial Lock of a provided .env file.
//
// Use [Open] to load an existing vault from disk, resolving the config and
// identity files automatically.
package vault
