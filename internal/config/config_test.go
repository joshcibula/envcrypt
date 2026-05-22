package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.EnvFile != ".env" {
		t.Errorf("expected EnvFile '.env', got %q", cfg.EnvFile)
	}
	if cfg.EncryptedFile != ".env.age" {
		t.Errorf("expected EncryptedFile '.env.age', got %q", cfg.EncryptedFile)
	}
	if cfg.KeyFile == "" {
		t.Error("expected non-empty KeyFile")
	}
	if cfg.RecipientsFile == "" {
		t.Error("expected non-empty RecipientsFile")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envcrypt.json")

	orig := &config.Config{
		KeyFile:        "mykey",
		RecipientsFile: "myrecipients",
		EnvFile:        "staging.env",
		EncryptedFile:  "staging.env.age",
	}

	if err := config.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.KeyFile != orig.KeyFile {
		t.Errorf("KeyFile mismatch: got %q, want %q", loaded.KeyFile, orig.KeyFile)
	}
	if loaded.EnvFile != orig.EnvFile {
		t.Errorf("EnvFile mismatch: got %q, want %q", loaded.EnvFile, orig.EnvFile)
	}
	if loaded.EncryptedFile != orig.EncryptedFile {
		t.Errorf("EncryptedFile mismatch: got %q, want %q", loaded.EncryptedFile, orig.EncryptedFile)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/.envcrypt.json")
	if err != os.ErrNotExist {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestLoadOrDefault_Missing(t *testing.T) {
	cfg, err := config.LoadOrDefault("/nonexistent/.envcrypt.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvFile != ".env" {
		t.Errorf("expected default EnvFile, got %q", cfg.EnvFile)
	}
}

func TestLoadOrDefault_Existing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envcrypt.json")

	orig := config.Default()
	orig.EnvFile = "production.env"
	if err := config.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	cfg, err := config.LoadOrDefault(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.EnvFile != "production.env" {
		t.Errorf("expected 'production.env', got %q", cfg.EnvFile)
	}
}
