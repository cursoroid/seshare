package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func crocPath() (string, error) {
	p, err := exec.LookPath("croc")
	if err != nil {
		return "", fmt.Errorf("croc not found on PATH — install it: https://github.com/schollz/croc")
	}
	return p, nil
}

// crocSend blocks until the peer receives `file` using the shared `code`.
func crocSend(file, code string) error {
	croc, err := crocPath()
	if err != nil {
		return err
	}
	// croc v10 forbids the code on the command line; it must come from the env.
	cmd := exec.Command(croc, "send", file)
	cmd.Env = append(os.Environ(), "CROC_SECRET="+code)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// crocRecv receives a single file into a fresh dir and returns its path.
func crocRecv(code string) (string, error) {
	croc, err := crocPath()
	if err != nil {
		return "", err
	}
	out, err := os.MkdirTemp("", "seshare-recv-")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(croc, "--yes", "--out", out)
	cmd.Env = append(os.Environ(), "CROC_SECRET="+code)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	files, _ := filepath.Glob(filepath.Join(out, "*"))
	if len(files) != 1 {
		return "", fmt.Errorf("expected one received file in %s, got %d", out, len(files))
	}
	return files[0], nil
}
