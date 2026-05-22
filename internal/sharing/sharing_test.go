package sharing_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/user/envcrypt/internal/sharing"
)

func TestParseRecipient(t *testing.T) {
	r, err := sharing.ParseRecipient("alice:age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "alice" {
		t.Errorf("expected name alice, got %q", r.Name)
	}
}

func TestParseRecipientInvalid(t *testing.T) {
	cases := []string{
		"",
		"nocoroln",
		":missingname",
		"missingkey:",
	}
	for _, c := range cases {
		_, err := sharing.ParseRecipient(c)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", c)
		}
	}
}

func writeTempRecipientsFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "recipients")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoadRecipientsFile(t *testing.T) {
	// Generate a real age identity to get a valid public key.
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating identity: %v", err)
	}
	pubKey := identity.Recipient().String()

	content := "# comment line\nalice:" + pubKey + "\n\nbob:" + pubKey + "\n"
	path := writeTempRecipientsFile(t, content)

	recipients, err := sharing.LoadRecipientsFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recipients) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(recipients))
	}
	if recipients[0].Name != "alice" || recipients[1].Name != "bob" {
		t.Errorf("unexpected recipient names: %v", recipients)
	}
}

func TestLoadRecipientsFileMissing(t *testing.T) {
	_, err := sharing.LoadRecipientsFile("/nonexistent/path/recipients")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestToAgeRecipients(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating identity: %v", err)
	}
	recipients := []sharing.Recipient{
		{Name: "alice", PublicKey: identity.Recipient().String()},
	}
	ageRecipients, err := sharing.ToAgeRecipients(recipients)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ageRecipients) != 1 {
		t.Errorf("expected 1 age recipient, got %d", len(ageRecipients))
	}
}

func TestToAgeRecipientsEmpty(t *testing.T) {
	_, err := sharing.ToAgeRecipients(nil)
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
}
