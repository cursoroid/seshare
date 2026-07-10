package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/schollz/croc/v10/src/croc"
	"github.com/schollz/croc/v10/src/models"
)

// relayPorts mirrors croc's default transfer ports (base 9009 + 4).
func relayPorts() []string {
	base := 9009
	ports := make([]string, 5)
	for i := range ports {
		ports[i] = strconv.Itoa(base + i)
	}
	return ports
}

func baseOptions(code string, sender bool) croc.Options {
	return croc.Options{
		IsSender:         sender,
		SharedSecret:     code,
		NoPrompt:         true,
		DisableClipboard: true,
		RelayAddress:     models.DEFAULT_RELAY,
		RelayAddress6:    models.DEFAULT_RELAY6,
		RelayPassword:    models.DEFAULT_PASSPHRASE,
		Curve:            "p256",
		HashAlgorithm:    "xxhash",
	}
}

// crocSend blocks until the peer receives `file`, authenticated by `code`.
func crocSend(file, code string) error {
	opts := baseOptions(code, true)
	opts.RelayPorts = relayPorts()
	c, err := croc.New(opts)
	if err != nil {
		return err
	}
	filesInfo, emptyFolders, totalFolders, err := croc.GetFilesInfo([]string{file}, false, false, nil)
	if err != nil {
		return err
	}
	return c.Send(filesInfo, emptyFolders, totalFolders)
}

// crocRecv receives one file into a fresh dir and returns its path. croc writes
// to the working directory, so we chdir into a temp dir and restore after.
func crocRecv(code string) (string, error) {
	out, err := os.MkdirTemp("", "seshare-recv-")
	if err != nil {
		return "", err
	}
	saved, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err := os.Chdir(out); err != nil {
		return "", err
	}
	defer os.Chdir(saved)

	c, err := croc.New(baseOptions(code, false))
	if err != nil {
		return "", err
	}
	if err := c.Receive(); err != nil {
		return "", err
	}

	files, _ := filepath.Glob(filepath.Join(out, "*"))
	if len(files) != 1 {
		return "", fmt.Errorf("expected one received file in %s, got %d", out, len(files))
	}
	return files[0], nil
}
