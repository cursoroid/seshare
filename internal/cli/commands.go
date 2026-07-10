package cli

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cursoroid/seshare/internal/contacts"
	"github.com/cursoroid/seshare/internal/ids"
	"github.com/cursoroid/seshare/internal/session"
	"github.com/cursoroid/seshare/internal/transcript"
	"github.com/cursoroid/seshare/internal/transport"
)

// ---- pair ---------------------------------------------------------------

func cmdPair(args []string) error {
	if len(args) == 1 && args[0] == "--list" {
		names, err := contacts.ListNames()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("no contacts yet")
			return nil
		}
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	}

	var name, code string
	var rotate bool
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--code":
			if i+1 >= len(args) {
				return fmt.Errorf("--code needs a value")
			}
			i++
			code = args[i]
		case "--rotate":
			rotate = true
		default:
			if name != "" {
				return fmt.Errorf("unexpected argument %q", args[i])
			}
			name = args[i]
		}
	}
	if name == "" {
		return fmt.Errorf("usage: seshare pair <name> [--code <code> | --rotate]")
	}

	if code != "" { // receiving side of the pairing (or accepting a rotated code)
		if err := contacts.Add(name, code); err != nil {
			return err
		}
		fmt.Printf("paired with %q\n", name)
		return nil
	}

	// initiating side: generate the shared code
	if _, err := contacts.Get(name); err == nil && !rotate {
		return fmt.Errorf("already paired with %q; use --rotate to generate a new code", name)
	}
	code = ids.NewCode()
	if err := contacts.Add(name, code); err != nil {
		return err
	}
	verb := "paired with"
	if rotate {
		verb = "rotated code for"
	}
	fmt.Printf("%s %q.\nSend them this code once (any channel):\n\n    %s\n\nThey run:  seshare pair <your-name> --code %s\n", verb, name, code, code)
	return nil
}

// ---- send ---------------------------------------------------------------

func cmdSend(args []string) error {
	var id, name string
	var yes bool
	for _, a := range args {
		switch {
		case a == "--yes" || a == "-y":
			yes = true
		case strings.HasPrefix(a, "@"):
			name = a[1:]
		case strings.HasPrefix(a, "-"):
			return fmt.Errorf("unknown flag %q", a)
		default:
			id = a
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path, err := session.Resolve(cwd, id)
	if err != nil {
		return err
	}

	if !yes {
		fmt.Printf("About to send %s\n", filepath.Base(path))
		fmt.Print("This transcript may contain secrets, file contents and absolute paths.\nOnly send to someone you trust. Continue? [y/N] ")
		if !confirm() {
			return fmt.Errorf("aborted")
		}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	gz, err := gzipToTemp(raw, "seshare-"+strings.TrimSuffix(filepath.Base(path), ".jsonl")+".jsonl.gz")
	if err != nil {
		return err
	}
	defer os.Remove(gz)

	code := ""
	if name != "" {
		if code, err = contacts.Get(name); err != nil {
			return err
		}
		fmt.Printf("sending to %q…\n", name)
	} else {
		code = ids.NewCode()
		fmt.Printf("one-time code (share it): %s\nrecipient runs:  seshare recv %s\n", code, code)
	}
	return transport.Send(gz, code)
}

// ---- recv ---------------------------------------------------------------

func cmdRecv(args []string) error {
	var target string
	var strip bool
	for _, a := range args {
		switch {
		case a == "--strip-snapshots":
			strip = true
		case strings.HasPrefix(a, "-") && a != "-":
			return fmt.Errorf("unknown flag %q", a)
		default:
			if target != "" {
				return fmt.Errorf("unexpected argument %q", a)
			}
			target = a
		}
	}
	if target == "" {
		return fmt.Errorf("usage: seshare recv <@name | code> [--strip-snapshots]")
	}

	code := target
	if strings.HasPrefix(target, "@") {
		c, err := contacts.Get(target[1:])
		if err != nil {
			return err
		}
		code = c
	}

	gzPath, err := transport.Recv(code)
	if err != nil {
		return err
	}
	defer os.Remove(gzPath)

	data, err := gunzipFile(gzPath)
	if err != nil {
		return err
	}
	if strip {
		data = transcript.StripSnapshots(data)
	}

	origCwd, origVer := transcript.PeekCwdVersion(data)
	if local := session.LocalClaudeVersion(); origVer != "" && local != "" && origVer != local {
		fmt.Fprintf(os.Stderr, "warning: session was made with Claude Code %s, you have %s\n", origVer, local)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	newID := ids.NewUUID()
	out := transcript.Rewrite(data, cwd, newID)

	dir := session.Dir(cwd)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	dest := filepath.Join(dir, newID+".jsonl")
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		return err
	}

	fmt.Printf("\nsession staged (originally from %s)\ncontinue it with:\n\n    cd %s && claude --resume %s\n", origCwd, cwd, newID)
	return nil
}

// ---- helpers ------------------------------------------------------------

func confirm() bool {
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

func gzipToTemp(data []byte, name string) (string, error) {
	f, err := os.CreateTemp("", "*-"+name)
	if err != nil {
		return "", err
	}
	zw := gzip.NewWriter(f)
	if _, err := zw.Write(data); err != nil {
		f.Close()
		return "", err
	}
	if err := zw.Close(); err != nil {
		f.Close()
		return "", err
	}
	return f.Name(), f.Close()
}

func gunzipFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("received file is not gzip (%w)", err)
	}
	defer zr.Close()
	return io.ReadAll(zr)
}
