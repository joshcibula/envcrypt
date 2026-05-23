package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envcrypt/internal/vault"
)

func TestExportRoundtrip(t *testing.T) {
	dir := t.TempDir()
	envContent := "APP_ENV=production\nSECRET=abc123\n"

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, vault.InitOptions{EnvFile: envFile}); err != nil {
		t.Fatalf("init: %v", err)
	}

	outPath := filepath.Join(dir, "exported.env")
	err := vault.Export(dir, vault.ExportOptions{OutputPath: outPath})
	if err != nil {
		t.Fatalf("export: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read exported: %v", err)
	}

	got := string(data)
	for _, want := range []string{"APP_ENV=production", "SECRET=abc123"} {
		if !containsLine(got, want) {
			t.Errorf("exported file missing line %q; got:\n%s", want, got)
		}
	}
}

func TestExportNoOverwrite(t *testing.T) {
	dir := t.TempDir()
	envContent := "KEY=value\n"

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, vault.InitOptions{EnvFile: envFile}); err != nil {
		t.Fatalf("init: %v", err)
	}

	outPath := filepath.Join(dir, "out.env")
	if err := os.WriteFile(outPath, []byte("existing"), 0600); err != nil {
		t.Fatalf("pre-create output: %v", err)
	}

	err := vault.Export(dir, vault.ExportOptions{OutputPath: outPath, Overwrite: false})
	if err == nil {
		t.Fatal("expected error when output file exists without --force")
	}
}

func TestExportOverwrite(t *testing.T) {
	dir := t.TempDir()
	envContent := "KEY=value\n"

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, vault.InitOptions{EnvFile: envFile}); err != nil {
		t.Fatalf("init: %v", err)
	}

	outPath := filepath.Join(dir, "out.env")
	if err := os.WriteFile(outPath, []byte("old content"), 0600); err != nil {
		t.Fatalf("pre-create output: %v", err)
	}

	err := vault.Export(dir, vault.ExportOptions{OutputPath: outPath, Overwrite: true})
	if err != nil {
		t.Fatalf("export with overwrite: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) == "old content" {
		t.Error("file was not overwritten")
	}
}

func TestExportMissingVault(t *testing.T) {
	dir := t.TempDir()
	err := vault.Export(dir, vault.ExportOptions{OutputPath: filepath.Join(dir, "out.env")})
	if err == nil {
		t.Fatal("expected error for missing vault")
	}
}
