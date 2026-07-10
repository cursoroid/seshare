// Command seshare shares Claude Code sessions peer-to-peer so a recipient can
// continue them with `claude --resume`. See docs/superpowers/specs.
package main

import (
	"os"

	"github.com/cursoroid/seshare/internal/cli"
)

func main() { os.Exit(cli.Main(os.Args[1:])) }
