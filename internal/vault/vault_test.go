package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/internal/vault"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestInitAndOpen(t *testing.T) {
	dir := t.TempDir()
	envPath := writeTempEnv(t, "APP_KEY=secret\nDEBUG=true\n")

	v, err := vault.Init(vault.InitOptions{
		CfgDir:  dir,
		EnvFile: envPath,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil vault")
	}

	v2, err := vault.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	vars, err := v2.Unlock()
	if err != nil {
		t.Fatalf("Unlock: %v", err)
	}
	if vars["APP_KEY"] != "secret" {
		t.Errorf("expected APP_KEY=secret, got %q", vars["APP_KEY"])
	}
	if vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", vars["DEBUG"])
	}
}

func TestInitForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	envPath := writeTempEnv(t, "FOO=bar\n")

	if _, err := vault.Init(vault.InitOptions{CfgDir: dir, EnvFile: envPath}); err != nil {
		t.Fatalf("first Init: %v", err)
	}
	if _, err := vault.Init(vault.InitOptions{CfgDir: dir, EnvFile: envPath}); err == nil {
		t.Fatal("expected error on second Init without force")
	}
	if _, err := vault.Init(vault.InitOptions{CfgDir: dir, EnvFile: envPath, Force: true}); err != nil {
		t.Fatalf("forced re-init: %v", err)
	}
}

func TestOpenMissingVault(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	if _, err := vault.Open(dir); err == nil {
		t.Fatal("expected error opening missing vault")
	}
}

func TestLockMissingEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := writeTempEnv(t, "X=1\n")
	v, err := vault.Init(vault.InitOptions{CfgDir: dir, EnvFile: envPath})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := v.Lock("/nonexistent/.env"); err == nil {
		t.Fatal("expected error locking missing file")
	}
}
