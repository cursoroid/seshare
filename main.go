// Command seshare shares Claude Code sessions peer-to-peer so a recipient can
// continue them with `claude --resume`. See docs/superpowers/specs.
package main

import (
	"os"
	"runtime/debug"

	"github.com/cursoroid/seshare/internal/cli"
)

// version is set via -ldflags at release time; "dev" otherwise.
var version = "dev"

func main() { os.Exit(cli.Main(os.Args[1:], resolveVersion())) }

// resolveVersion prefers the ldflags value, falling back to the module version
// recorded by `go install github.com/cursoroid/seshare@vX.Y.Z`.
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return version
}
