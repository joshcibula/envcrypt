// Package config handles loading and saving envcrypt project configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultConfigFile = ".envcrypt.json"

// Config holds project-level envcrypt settings.
type Config struct {
	// KeyFile is the path to the age identity (private key) file.
	KeyFile string `json:"key_file"`
	// RecipientsFile is the path to the shared recipients list.
	RecipientsFile string `json:"recipients_file"`
	// EnvFile is the plaintext .env file path.
	EnvFile string `json:"env_file"`
	// EncryptedFile is the encrypted output file path.
	EncryptedFile string `json:"encrypted_file"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		KeyFile:        ".envcrypt_key",
		RecipientsFile: ".envcrypt_recipients",
		EnvFile:        ".env",
		EncryptedFile:  ".env.age",
	}
}

// Validate checks that all required fields in the Config are non-empty.
func (c *Config) Validate() error {
	if c.KeyFile == "" {
		return fmt.Errorf("config: key_file must not be empty")
	}
	if c.RecipientsFile == "" {
		return fmt.Errorf("config: recipients_file must not be empty")
	}
	if c.EnvFile == "" {
		return fmt.Errorf("config: env_file must not be empty")
	}
	if c.EncryptedFile == "" {
		return fmt.Errorf("config: encrypted_file must not be empty")
	}
	return nil
}

// Save writes the config as JSON to the given path.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// Load reads and parses a config file from the given path.
// Returns os.ErrNotExist if the file does not exist.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadOrDefault attempts to load the config at path; if not found, returns Default().
func LoadOrDefault(path string) (*Config, error) {
	cfg, err := Load(path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	return cfg, err
}
