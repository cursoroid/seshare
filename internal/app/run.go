// seshare — move a Claude Code session to another machine over croc so the
// recipient can continue it with `claude --resume`. See docs/superpowers/specs.
package app

import (
	"fmt"
	"os"
)

const usage = `seshare — share Claude Code sessions peer-to-peer

usage:
  seshare pair <name> [--code <code>]   pair with someone (exchange a code once)
  seshare pair <name> --rotate          replace a contact's code (re-share once)
  seshare pair --list                   list contacts
  seshare send [session-id] [@name]     send newest (or given) session
  seshare recv <@name | code>           receive and stage a session

run a command with no valid args for its own help.`

// Main runs the CLI and returns the process exit code.
func Main(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, usage)
		return 2
	}

	var err error
	switch args[0] {
	case "pair":
		err = cmdPair(args[1:])
	case "send":
		err = cmdSend(args[1:])
	case "recv":
		err = cmdRecv(args[1:])
	case "-h", "--help", "help":
		fmt.Println(usage)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", args[0], usage)
		return 2
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	return 0
}
