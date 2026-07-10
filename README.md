# seshare

Share a Claude Code session peer-to-peer so someone else can **continue** it
with `claude --resume`. Transfers go directly between the two machines over
[croc](https://github.com/schollz/croc) — no server, no account.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/cursoroid/seshare/main/install.sh | sh
```

Runs `go install` and adds Go's bin dir to your `PATH` if it isn't already.
croc is embedded — no separate install needed.

Prefer to do it by hand? `go install github.com/cursoroid/seshare@latest`, then
ensure `$(go env GOPATH)/bin` is on your `PATH`.

## Update

Re-run the install command; `@latest` pulls the newest release:

```sh
go install github.com/cursoroid/seshare@latest
```

Pin a specific version with `@v1.2.3` instead of `@latest`.

## Use

Pair once with a person (exchange the code a single time, any channel):

```sh
# you
seshare pair alice
# -> prints a code; send it to alice once

# alice
seshare pair you --code <that-code>
```

Then send and receive by name — no code typing:

```sh
# sender (in the repo dir the session belongs to)
seshare send @alice              # newest session in this dir
seshare send <session-id> @alice # a specific one

# recipient (in the dir they want to continue from)
seshare recv @you
# -> cd <dir> && claude --resume <new-id>
```

Both sides must be online at the same time (it's a live P2P handoff).

One-off without pairing: `seshare send` prints a one-time code; recipient runs
`seshare recv <code>`.

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
