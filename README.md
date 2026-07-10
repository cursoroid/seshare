# seshare

Share a Claude Code session peer-to-peer so someone else can **continue** it
with `claude --resume`. Transfers go directly between the two machines over
[croc](https://github.com/schollz/croc) — no server, no account.

## Install

Homebrew (macOS/Linux) — handles PATH for you:

```sh
brew install cursoroid/homebrew-tap/seshare
```

Or with Go (installs to `$(go env GOPATH)/bin`; the script also fixes PATH):

```sh
curl -fsSL https://raw.githubusercontent.com/cursoroid/seshare/main/install.sh | sh
# or by hand: go install github.com/cursoroid/seshare@latest
```

croc is embedded — no separate install needed.

## Update

```sh
brew upgrade seshare                                  # Homebrew
go install github.com/cursoroid/seshare@latest        # Go (@vX.Y.Z to pin)
```

## Use

Pair once — each of you saves the *other* person under a local name of your
choosing. The two names are independent (they don't have to match); the code is
exchanged a single time:

```sh
# on your machine — you save your friend as "bob"
seshare pair bob
# -> prints a code; send it to bob once

# on bob's machine — bob saves you as "alice"
seshare pair alice --code <that-code>
```

Then send and receive by name — no code typing:

```sh
# you (sender), in the repo dir the session belongs to — use the name YOU saved
seshare send bob                # newest session in this dir
seshare send <session-id> bob   # a specific one

# bob (recipient), in the dir he wants to continue from — the name HE saved for you
seshare recv alice        # prints: cd <dir> && claude --resume <new-id>
seshare recv alice -r     # ...or jump straight into `claude --resume`
```

The `@` prefix is optional — `send bob` and `send @bob` are the same. Use `@`
only if you ever have a contact whose name collides with a session id.

> **Names are local.** On `recv`, use the name *you* saved the **sender** under
> — not your own name. Contact names are per-machine and asymmetric: you call
> your friend `bob`, he calls you `alice`. A bare word that isn't one of your
> contacts is treated as a raw code, so a wrong name just fails (after a 2-min
> timeout). The raw code always works regardless of names: `seshare recv <code>`.

Both sides must be online at the same time (it's a live P2P handoff).

One-off without pairing: `seshare send` prints a one-time code; recipient runs
`seshare recv <code>`.

Prefer to pick visually? `seshare tui` browses your sessions across all
projects with a preview pane; hit enter to send the selected one to a paired
contact.

## Notes

- The transcript may contain secrets, file contents and absolute paths. `send`
  warns before uploading; only send to people you trust.
- The **conversation** resumes anywhere; tool results pointing at the sender's
  absolute paths or a specific repo won't re-resolve unless the recipient has
  the same code checked out.
- If `--resume` trips on the sender's local file snapshots, receive with
  `seshare recv @you --strip-snapshots`.
- Contacts live in `~/.seshare/contacts.json` (`0600`). A per-contact code is a
  permanent shared secret; rotate it with `seshare pair <name> --rotate` (then
  re-share the new code once).

See `docs/superpowers/specs/` for the design.
