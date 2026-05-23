package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourusername/envcrypt/cmd/envcrypt/commands"
	"github.com/yourusername/envcrypt/internal/vault"
)

func setupExportVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("EXPORT_KEY=hello\nOTHER=world\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	if err := vault.Init(dir, vault.InitOptions{EnvFile: envFile}); err != nil {
		t.Fatalf("init vault: %v", err)
	}
	return dir
}

func runExportCmd(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(old) })

	root := &cobra.Command{Use: "envcrypt"}
	// Re-register export command via the exported constructor.
	_ = root

	var buf bytes.Buffer
	cmd := commands.NewExportCmdForTest()
	cmd.SetOut(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestExportCmdCreatesFile(t *testing.T) {
	dir := setupExportVault(t)
	out := filepath.Join(dir, "result.env")

	_, err := runExportCmd(t, dir, "--output", out)
	if err != nil {
		t.Fatalf("export cmd: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if !strings.Contains(string(data), "EXPORT_KEY=hello") {
		t.Errorf("expected EXPORT_KEY in output, got: %s", data)
	}
}

func TestExportCmdForceFlag(t *testing.T) {
	dir := setupExportVault(t)
	out := filepath.Join(dir, "result.env")

	if err := os.WriteFile(out, []byte("old"), 0600); err != nil {
		t.Fatalf("pre-create: %v", err)
	}

	_, err := runExportCmd(t, dir, "--output", out, "--force")
	if err != nil {
		t.Fatalf("export with --force: %v", err)
	}

	data, _ := os.ReadFile(out)
	if string(data) == "old" {
		t.Error("file was not overwritten with --force")
	}
}
