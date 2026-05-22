// Package sharing provides utilities for sharing encrypted .env files
// with other recipients using their age public keys.
package sharing

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"filippo.io/age"
)

// Recipient represents a named age recipient (public key holder).
type Recipient struct {
	Name      string
	PublicKey string
}

// ParseRecipient parses a recipient from a "name:publickey" string.
func ParseRecipient(s string) (Recipient, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Recipient{}, fmt.Errorf("invalid recipient format %q: expected name:publickey", s)
	}
	return Recipient{Name: parts[0], PublicKey: parts[1]}, nil
}

// LoadRecipientsFile reads a recipients file where each non-comment line
// is in the format "name:publickey".
func LoadRecipientsFile(path string) ([]Recipient, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("recipients file not found: %s", path)
		}
		return nil, fmt.Errorf("reading recipients file: %w", err)
	}

	var recipients []Recipient
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r, err := ParseRecipient(line)
		if err != nil {
			return nil, err
		}
		recipients = append(recipients, r)
	}
	return recipients, nil
}

// ToAgeRecipients converts a slice of Recipient to age.Recipient instances.
func ToAgeRecipients(recipients []Recipient) ([]age.Recipient, error) {
	if len(recipients) == 0 {
		return nil, errors.New("at least one recipient is required")
	}
	var ageRecipients []age.Recipient
	for _, r := range recipients {
		parsed, err := age.ParseX25519Recipient(r.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid public key for recipient %q: %w", r.Name, err)
		}
		ageRecipients = append(ageRecipients, parsed)
	}
	return ageRecipients, nil
}
