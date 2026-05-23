package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/cmd/envcrypt/commands"
	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/crypto"
	"github.com/nicholasgasior/envcrypt/internal/envfile"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func setupEditVault(t *testing.T) (dir, vaultPath, keyPath string) {
	t.Helper()
	dir = t.TempDir()
	envPath := filepath.Join(dir, ".env")
	vaultPath = filepath.Join(dir, "vault.age")
	keyPath = filepath.Join(dir, "key.txt")

	_ = os.WriteFile(envPath, []byte("HELLO=world\n"), 0600)
	if err := vault.Init(envPath, vaultPath, keyPath, false); err != nil {
		t.Fatalf("init vault: %v", err)
	}

	cfg := config.Default()
	cfg.VaultPath = vaultPath
	cfg.KeyPath = keyPath
	_ = config.Save(cfg, filepath.Join(dir, ".envcrypt.json"))

	return dir, vaultPath, keyPath
}

func runEditCmd(t *testing.T, cfgPath, key, value string) (string, error) {
	t.Helper()
	cmd := commands.NewEditCmdForTest()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--config", cfgPath, key, value})
	err := cmd.Execute()
	return buf.String(), err
}

func TestEditCmdUpdatesKey(t *testing.T) {
	dir, vaultPath, keyPath := setupEditVault(t)
	cfgPath := filepath.Join(dir, ".envcrypt.json")

	_, err := runEditCmd(t, cfgPath, "HELLO", "envcrypt")
	if err != nil {
		t.Fatalf("edit cmd: %v", err)
	}

	identity, _ := keystore.Load(keyPath)
	v, _ := vault.Open(vaultPath)
	plaintext, _ := crypto.Decrypt(v.Ciphertext, identity)
	envMap, _ := envfile.Parse(string(plaintext))

	if envMap["HELLO"] != "envcrypt" {
		t.Errorf("expected HELLO=envcrypt, got %s", envMap["HELLO"])
	}
}
