---
description: Send this Claude Code session to a paired seshare contact
argument-hint: [contact] [session-id]
allowed-tools: Bash
---

Send the current session with seshare. Arguments: `$ARGUMENTS` (a contact name,
optionally a session id; empty means "newest session in this directory", and no
contact means "print a one-time code").

1. If `command -v seshare` fails, stop and tell the user:
   `brew install cursoroid/homebrew-tap/seshare` (or the `install.sh` one-liner).
2. Confirm before uploading. `seshare send` normally prompts about this itself,
   but it cannot prompt through a tool call, so **you** ask instead: name the
   session file (`seshare send` picks the newest `.jsonl` for this directory)
   and warn that a transcript may contain secrets, file contents and absolute
   paths. Wait for an explicit yes. Do not run the command before that.
3. Run `seshare send --yes $ARGUMENTS` **in the background** — it blocks on a
   live croc transfer until the recipient connects.
4. Report the `seshare recv <code>` line it prints so the user can pass it on,
   then watch the background output and say whether the transfer completed.

Notes for the user, if relevant:
- The transcript on disk lags the live conversation slightly, so the last turn
  or two may not make it across.
- Contact names are local. Whatever you call them, they receive with the name
  *they* saved you under — the raw code always works either way.
