#!/bin/sh
# Install seshare and make sure Go's bin dir is on PATH.
# Usage: curl -fsSL https://raw.githubusercontent.com/cursoroid/seshare/main/install.sh | sh
set -e

if ! command -v go >/dev/null 2>&1; then
	echo "Go is required: https://go.dev/dl/" >&2
	exit 1
fi

echo "installing seshare…"
go install github.com/cursoroid/seshare@latest

bin="$(go env GOBIN)"
[ -n "$bin" ] || bin="$(go env GOPATH)/bin"

# already on PATH? done.
case ":$PATH:" in
	*":$bin:"*) echo "seshare installed — $bin is already on PATH."; exit 0 ;;
esac

# pick the rc for the login shell. ponytail: zsh/bash cover ~all macOS/Linux;
# other shells fall through to the manual hint below.
case "$(basename "${SHELL:-}")" in
	zsh)  rc="$HOME/.zshrc" ;;
	bash) rc="$HOME/.bashrc" ;;
	*)    rc="" ;;
esac

line="export PATH=\"\$PATH:$bin\""
if [ -n "$rc" ] && ! { [ -f "$rc" ] && grep -qF "$bin" "$rc"; }; then
	printf '\n# added by seshare install\n%s\n' "$line" >>"$rc"
	echo "seshare installed. Added $bin to PATH in $rc."
	echo "Run:  source $rc   (or open a new terminal)"
else
	echo "seshare installed. Add this to your shell profile:"
	echo "  $line"
fi
