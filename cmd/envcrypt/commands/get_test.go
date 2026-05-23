package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
)

func setupGetVault(t *testing.T) (vaultPath, keyPath string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath = filepath.Join(dir, ".env.enc")
	keyPath = filepath.Join(dir, "key.age")
	envPath := filepath.Join(dir, ".env")

	_ = os.WriteFile(envPath, []byte("GREETING=hello\nNAME=world\n"), 0600)

	identity, _ := keystore.Generate()
	_ = keystore.Save(identity, keyPath)
	recipient, _ := identity.Recipient()
	data := envfile.Serialize(map[string]string{"GREETING": "hello", "NAME": "world"})
	encrypted, _ := crypto.Encrypt([]byte(data), []age.Recipient{recipient})
	_ = os.WriteFile(vaultPath, encrypted, 0600)

	t.Setenv("ENVCRYPT_VAULT", vaultPath)
	t.Setenv("ENVCRYPT_KEY", keyPath)

	_ = config.Save(config.Config{VaultPath: vaultPath, KeyPath: keyPath}, filepath.Join(dir, "config.toml"))
	return
}

func runGetCmd(t *testing.T, vaultPath, keyPath, key string) (string, error) {
	t.Helper()
	cmd := NewGetCmdForTest(vaultPath, keyPath)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{key})
	err := cmd.Execute()
	return strings.TrimSpace(buf.String()), err
}

func TestGetCmdReturnsValue(t *testing.T) {
	vaultPath, keyPath := setupGetVault(t)
	out, err := runGetCmd(t, vaultPath, keyPath, "GREETING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello" {
		t.Errorf("expected \"hello\", got %q", out)
	}
}

func TestGetCmdMissingKey(t *testing.T) {
	vaultPath, keyPath := setupGetVault(t)
	_, err := runGetCmd(t, vaultPath, keyPath, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
