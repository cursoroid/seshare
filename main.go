// seshare — move a Claude Code session to another machine over croc so the
// recipient can continue it with `claude --resume`. See docs/superpowers/specs.
package main

import (
	"fmt"
	"os"
)

const usage = `seshare — share Claude Code sessions peer-to-peer

usage:
  seshare pair <name> [--code <code>]   pair with someone (exchange a code once)
  seshare pair --list                   list contacts
  seshare send [session-id] [@name]     send newest (or given) session
  seshare recv <@name | code>           receive and stage a session

run a command with no valid args for its own help.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "pair":
		err = cmdPair(os.Args[2:])
	case "send":
		err = cmdSend(os.Args[2:])
	case "recv":
		err = cmdRecv(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Println(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", os.Args[1], usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
