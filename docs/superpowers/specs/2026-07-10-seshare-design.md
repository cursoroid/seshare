# seshare — design

**Date:** 2026-07-10
**Status:** approved, pre-implementation

## Purpose

Move a Claude Code session from one machine to another so the recipient can
**continue** it with `claude --resume`. Transport is peer-to-peer over `croc`
(code-phrase, encrypted, relay fallback). A one-time pairing step lets known
people exchange the code only once and reuse it forever.

## Why this shape

A Claude Code session is a JSONL transcript at
`~/.claude/projects/<munged-cwd>/<session-id>.jsonl`, where `<munged-cwd>` is
the working directory with `/` replaced by `-` (no real hash). Each line
carries `cwd`, `version`, `sessionId`, and message `uuid`/`parentUuid`. Import
therefore needs no separate manifest — everything is in the transcript.

Rejected alternatives (during brainstorming): gist (async but content sits in a
URL-addressable store), Mongo (needs an embedded credential or a fronting API —
sank the "share with anyone" story), rotating derived codes (more crypto than a
personal tool needs), hand-rolled WebRTC (needs a signaling server anyway).
P2P via `croc` kills the credential/access problem with zero infra; the only
cost is sender and recipient must be online at the same time.

## Storage

`~/.seshare/contacts.json` = `{ "<name>": "<code>" }`, file perms `0600`
(codes are transfer secrets). Created on first `pair`.

## Commands

### `seshare pair <name> [--code <code>]`
- No `--code`: generate a random code (`crypto/rand` → hex), store under
  `<name>`, print it. User sends that code **once** to the other person.
- `--code <code>`: store the supplied code under `<name>` (the receiving side
  of the pairing).
- `seshare pair --list`: print stored contact names (never print codes).

### `seshare send [session-id] [@name]`
1. Resolve session: default = newest `*.jsonl` in
   `~/.claude/projects/<munged-$PWD>/`; or the explicit `session-id`.
2. Warn + require confirmation: transcript may contain secrets / absolute
   paths; only send to someone trusted. (P2P is direct + encrypted to one
   recipient, so no content scrubbing in v1.)
3. gzip the transcript.
4. Transfer via `croc send`:
   - `@name` given → use that contact's stored code.
   - no `@name` → let croc emit a fresh one-time code and print it.

### `seshare recv <@name | code>`
1. `croc recv` using the contact's stored code (`@name`) or the given
   one-time code. Ungzip.
2. Read `cwd` and `version` from the transcript's own lines. Warn if `version`
   differs from the locally installed Claude Code version.
3. Rewrite every line:
   - `cwd` → recipient's `$PWD`.
   - `sessionId` (field and the output filename) → a fresh UUID
     (`crypto/rand`), to avoid collisions with the recipient's existing
     sessions. Message `uuid`/`parentUuid` chains are left intact.
4. Write to `~/.claude/projects/<munged-$PWD>/<new-id>.jsonl`.
5. Print: `cd <dir> && claude --resume <new-id>`.
6. `--strip-snapshots` (optional): drop `file-history-snapshot` /
   `isSnapshotUpdate` lines before writing — fallback if resume chokes on
   snapshot references.

## Dependencies

- `croc` on `PATH` (v1 shells out via `os/exec`). Upgrade path: import
  `github.com/schollz/croc` as a library for a single self-contained binary.
- Go stdlib only otherwise: `compress/gzip`, `encoding/json`, `crypto/rand`,
  `os/exec`, `path/filepath`.

## Known risks (accepted for v1)

- `file-history-snapshot` lines reference the sender's local files/snapshots.
  The **conversation** resumes fine (context lives in the transcript);
  snapshot/checkpoint *restore* may not. Fallback: `--strip-snapshots`.
- Tool results referencing absolute paths or a specific repo won't re-resolve
  unless the recipient has the same code checked out.
- Static per-contact code is a **permanent shared secret**: anyone who ever
  learns it can intercept future transfers between that pair. Accepted for a
  known-people tool. Mark in code with a `ponytail:` comment; upgrade path is
  `HMAC(secret, time-bucket)` rotating derivation (TOTP-style, ±1 window on
  receive for clock skew).

## Out of scope

Secret scrubbing/redaction, central storage (Mongo/gist), rotating codes,
live/streaming view, GUI, any config beyond `contacts.json`.

## Checks

Non-trivial logic gets one runnable check each (Go `testing`, no framework):
- munged-cwd derivation round-trips a known path.
- import rewrite: given a sample transcript, `cwd` → `$PWD` on every line and
  `sessionId` is a single fresh UUID across all lines.
- `--strip-snapshots` removes exactly the snapshot line types.
