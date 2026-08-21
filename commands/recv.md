---
description: Receive a Claude Code session someone is sending over seshare
argument-hint: <contact-or-code>
allowed-tools: Bash
---

Receive a session with seshare. Argument: `$ARGUMENTS` — the name you saved the
*sender* under, or the raw code they gave you.

1. If `command -v seshare` fails, stop and tell the user:
   `brew install cursoroid/homebrew-tap/seshare`.
2. Run `seshare recv $ARGUMENTS` **in the background** from the directory the
   user wants to continue in. It waits for the sender (both sides must be online)
   and gives up after about two minutes.
3. On success it prints `cd <dir> && claude --resume <id>`. Give the user that
   line to run themselves — `claude --resume` cannot start inside this session.
   In Claude Code they can run it directly by typing `!` followed by the command.
4. If `--resume` later chokes on the sender's local file snapshots, re-receive
   with `seshare recv $ARGUMENTS --strip-snapshots`.

Never pass `-r` / `--resume` to `seshare recv` from here: it would try to launch
a nested Claude Code.
