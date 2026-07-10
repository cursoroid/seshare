// Package contacts stores per-person croc share codes in ~/.seshare/contacts.json.
package contacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// SeshareDir holds contacts.json. Exported so tests can redirect it.
var SeshareDir = defaultSeshareDir()

func defaultSeshareDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".seshare"
	}
	return filepath.Join(home, ".seshare")
}

func contactsFile() string { return filepath.Join(SeshareDir, "contacts.json") }

func load() (map[string]string, error) {
	data, err := os.ReadFile(contactsFile())
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("contacts file corrupt: %w", err)
	}
	return m, nil
}

// Add stores (or replaces) the code for a contact name.
func Add(name, code string) error {
	m, err := load()
	if err != nil {
		return err
	}
	m[name] = code
	if err := os.MkdirAll(SeshareDir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(contactsFile(), data, 0o600)
}

// Get returns a contact's code, or an error if unknown.
func Get(name string) (string, error) {
	m, err := load()
	if err != nil {
		return "", err
	}
	code, ok := m[name]
	if !ok {
		return "", fmt.Errorf("no contact named %q (run: seshare pair %s)", name, name)
	}
	return code, nil
}

// ListNames returns contact names, sorted.
func ListNames() ([]string, error) {
	m, err := load()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}
