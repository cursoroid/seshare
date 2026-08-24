---
name: seshare
description: Use when the user wants to hand off, share, send, or pick up a Claude Code session — "send this session to bob", "share this conversation with a teammate", "continue the session alice sent me", "pair with someone on seshare", or when listing/rotating seshare contacts.
---

# seshare

`seshare` moves a Claude Code session transcript between two machines over croc
(direct P2P, no server, no account) so the other person can continue it with
`claude --resume`.

## Before anything

`command -v seshare`. If missing:
`brew install cursoroid/homebrew-tap/seshare`, or
`curl -fsSL https://raw.githubusercontent.com/cursoroid/seshare/main/install.sh | sh`.

## Pairing (once per person)

Names are **local and asymmetric** — you save them as `bob`, they save you as
`alice`. Get this wrong and `recv` just times out.

```sh
seshare pair bob            # prints a code; user shares it with bob once
seshare pair bob --code X   # the other side, saving the code they were sent
seshare pair --list         # contacts
seshare pair bob --rotate   # new code; re-share once
```

Contacts live in `~/.seshare/contacts.json` (0600). A per-contact code is a
permanent shared secret.

## Sending

Ask for confirmation in chat first — a transcript can hold secrets, file
contents and absolute paths — then run with `--yes`, because the tool call has
no TTY for seshare's own prompt.

```sh
seshare send --yes bob                # newest session in this directory
seshare send --yes <session-id> bob   # a specific one
seshare send --yes                    # no contact: prints a one-time code
```

A contact's code is a **permanent** shared secret and `send bob` prints it. Never
repeat it back — say "sending to bob" and stop. Only the one-time code from a
contactless send is safe to relay.

An unpaired or misspelled name is parsed as a session id, so `send bpb <id>`
quietly falls back to a one-time code. Read the output before reporting success:
`sending to "bob"` reached the contact, `one-time code` did not.

Run it in the background: it blocks until the recipient connects, with no
timeout of its own. Give up after a few minutes and kill the shell, or the croc
process and its temp copy of the transcript outlive the turn. Both sides must be
online at the same time. The on-disk transcript trails the live conversation, so
the newest turn or two may not travel.

## Receiving

```sh
seshare recv alice                    # name the user saved the SENDER under
seshare recv <code>                   # raw code, always works
seshare recv alice --strip-snapshots  # if --resume chokes on file snapshots
```

The session is staged for whatever directory the command runs in, so ask which
one the user wants and pass it: `cd <dir> && seshare recv alice`.

Background it too. The whole receive — finding the sender and moving the file —
must finish inside about two minutes. It ends with `cd <dir> && claude --resume
<id>`: hand that to the user, do not run it, and never pass `-r`/`--resume` to
`recv` from a tool call. Both exec `claude --resume` on the caller's stdio, which
starts a nested Claude Code that never exits and hangs the call. Suggest the user
type `!` followed by the command.

`--strip-snapshots` applies on receive, so it needs a whole new live transfer:
the sender has to send again (and a spent one-time code needs replacing).

The conversation resumes anywhere, but tool results pointing at the sender's
absolute paths won't re-resolve unless the recipient has the same code checked
out.

## Browsing

`seshare tui` browses sessions across all projects with a preview pane and sends
the selected one. It needs a real terminal — tell the user to run it themselves
rather than calling it from a tool.
