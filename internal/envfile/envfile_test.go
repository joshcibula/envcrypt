package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envcrypt/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParse(t *testing.T) {
	path := writeTempEnv(t, `# comment
DB_HOST=localhost
DB_PORT=5432
SECRET="mysecret"
`)

	entries, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	cases := []struct{ key, value string }{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"SECRET", "mysecret"},
	}

	for i, c := range cases {
		if entries[i].Key != c.key || entries[i].Value != c.value {
			t.Errorf("entry %d: expected %s=%s, got %s=%s", i, c.key, c.value, entries[i].Key, entries[i].Value)
		}
	}
}

func TestParseMissingFile(t *testing.T) {
	_, err := envfile.Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseInvalidSyntax(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_WITHOUT_EQUALS\n")
	_, err := envfile.Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid syntax, got nil")
	}
}

func TestSerialize(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	out := envfile.Serialize(entries)
	expected := "FOO=bar\nBAZ=qux\n"
	if out != expected {
		t.Errorf("Serialize: expected %q, got %q", expected, out)
	}
}
