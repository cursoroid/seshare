package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// seshareDir holds contacts.json. Overridable in tests.
var seshareDir = defaultSeshareDir()

func defaultSeshareDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".seshare"
	}
	return filepath.Join(home, ".seshare")
}

func contactsFile() string { return filepath.Join(seshareDir, "contacts.json") }

func loadContacts() (map[string]string, error) {
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

func addContact(name, code string) error {
	m, err := loadContacts()
	if err != nil {
		return err
	}
	m[name] = code
	if err := os.MkdirAll(seshareDir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(contactsFile(), data, 0o600)
}

func getContact(name string) (string, error) {
	m, err := loadContacts()
	if err != nil {
		return "", err
	}
	code, ok := m[name]
	if !ok {
		return "", fmt.Errorf("no contact named %q (run: seshare pair %s)", name, name)
	}
	return code, nil
}

func listContactNames() ([]string, error) {
	m, err := loadContacts()
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
