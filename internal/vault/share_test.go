package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/config"
	"github.com/nicholasgasior/envcrypt/internal/keystore"
	"github.com/nicholasgasior/envcrypt/internal/vault"
)

func TestShareReEncryptsForRecipients(t *testing.T) {
	dir := t.TempDir()

	// Create owner identity.
	ownerIdentityFile := filepath.Join(dir, "owner.key")
	ownerIdentity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate owner identity: %v", err)
	}
	if err := keystore.Save(ownerIdentity, ownerIdentityFile); err != nil {
		t.Fatalf("save owner identity: %v", err)
	}

	// Create a second identity to share with.
	recipientIdentity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate recipient identity: %v", err)
	}

	// Write a recipients file with the second identity's public key.
	recsFile := filepath.Join(dir, "recipients.txt")
	pubKey := recipientIdentity.Recipient().String()
	if err := os.WriteFile(recsFile, []byte(pubKey+"\n"), 0600); err != nil {
		t.Fatalf("write recipients file: %v", err)
	}

	// Initialise a vault.
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("SECRET=hello\n"), 0600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	vaultFile := filepath.Join(dir, "vault.age")
	cfg := &config.Config{
		IdentityFile: ownerIdentityFile,
		VaultFile:    vaultFile,
		EnvFile:      envFile,
	}
	if err := vault.Init(cfg, vault.InitOptions{EnvFile: envFile, Force: false}); err != nil {
		t.Fatalf("vault init: %v", err)
	}

	// Share the vault.
	err = vault.Share(cfg, vault.ShareOptions{
		VaultPath:      vaultFile,
		RecipientsFile: recsFile,
		OutputPath:     vaultFile,
	})
	if err != nil {
		t.Fatalf("share: %v", err)
	}

	// The recipient should now be able to decrypt the vault.
	outDir := t.TempDir()
	recipientCfg := &config.Config{
		IdentityFile: filepath.Join(outDir, "recipient.key"),
		VaultFile:    vaultFile,
		EnvFile:      filepath.Join(outDir, ".env"),
	}
	if err := keystore.Save(recipientIdentity, recipientCfg.IdentityFile); err != nil {
		t.Fatalf("save recipient identity: %v", err)
	}
	if err := vault.Unlock(recipientCfg, vault.UnlockOptions{}); err != nil {
		t.Fatalf("recipient unlock: %v", err)
	}

	data, err := os.ReadFile(recipientCfg.EnvFile)
	if err != nil {
		t.Fatalf("read decrypted env: %v", err)
	}
	if string(data) != "SECRET=hello\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestShareMissingVault(t *testing.T) {
	dir := t.TempDir()

	identityFile := filepath.Join(dir, "owner.key")
	identity, err := keystore.Generate()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	if err := keystore.Save(identity, identityFile); err != nil {
		t.Fatalf("save identity: %v", err)
	}

	recsFile := filepath.Join(dir, "recipients.txt")
	if err := os.WriteFile(recsFile, []byte(identity.Recipient().String()+"\n"), 0600); err != nil {
		t.Fatalf("write recipients file: %v", err)
	}

	cfg := &config.Config{
		IdentityFile: identityFile,
		VaultFile:    filepath.Join(dir, "nonexistent.age"),
	}

	err = vault.Share(cfg, vault.ShareOptions{
		VaultPath:      cfg.VaultFile,
		RecipientsFile: recsFile,
	})
	if err == nil {
		t.Fatal("expected error for missing vault, got nil")
	}
}
