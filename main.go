// Command seshare shares Claude Code sessions peer-to-peer so a recipient can
// continue them with `claude --resume`. See docs/superpowers/specs.
package main

import (
	"os"

	"github.com/cursoroid/seshare/internal/app"
)

func main() { os.Exit(app.Main(os.Args[1:])) }
