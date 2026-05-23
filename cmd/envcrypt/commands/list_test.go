package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/vault"
	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envcrypt/cmd/envcrypt/commands"
)

func setupListVault(t *testing.T, envContent string) string {
	t.Helper()
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	if err := vault.Init(dir, envFile, false); err != nil {
		t.Fatalf("init vault: %v", err)
	}
	return dir
}

func runListCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envcrypt"}
	root.AddCommand(commands.NewListCmdForTest())
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs(append([]string{"list"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestListCmdOutputsKeys(t *testing.T) {
	dir := setupListVault(t, "FOO=bar\nBAZ=qux\n")
	out, err := runListCmd(t, "--dir", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAZ") {
		t.Errorf("expected FOO and BAZ in output, got: %q", out)
	}
	if strings.Contains(out, "bar") || strings.Contains(out, "qux") {
		t.Errorf("output should not contain values, got: %q", out)
	}
}

func TestListCmdEmptyVault(t *testing.T) {
	dir := setupListVault(t, "")
	out, err := runListCmd(t, "--dir", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no keys found") {
		t.Errorf("expected empty message, got: %q", out)
	}
}
