#!/bin/sh

set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
GO_PATH=$(go env GOPATH)
GO_CACHE=$(go env GOCACHE)
GO_MOD_CACHE=$(go env GOMODCACHE)
TMP_HOME=$(mktemp -d)
BIN="$TMP_HOME/spd"
WORKSPACE="$TMP_HOME/workspace"
MARKER="$TMP_HOME/should-not-run"
export HOME="$TMP_HOME"
export GOPATH="$GO_PATH"
export GOCACHE="$GO_CACHE"
export GOMODCACHE="$GO_MOD_CACHE"
export XDG_CONFIG_HOME="$TMP_HOME/.config"

cleanup() {
	rm -rf "$TMP_HOME"
}
trap cleanup EXIT INT TERM

mkdir -p "$WORKSPACE"
if [ -n "${SPADE_TEST_BIN:-}" ]; then
	cp "$SPADE_TEST_BIN" "$BIN"
else
	go build -o "$BIN" "$ROOT"
fi

"$BIN" --version | grep -Fx "spade version ${SPADE_TEST_VERSION:-0.1.1}"

"$BIN" add smoke-source printf 'hello\nworld\n'
"$BIN" add smoke-upper tr 'a-z' 'A-Z'
"$BIN" add smoke-count wc -- -c
"$BIN" chain add smoke-pipeline smoke-source then smoke-upper then smoke-count

pipeline_output=$("$BIN" chain smoke-pipeline | awk 'NF { value=$NF } END { print value }')
test "$pipeline_output" = "12"

"$BIN" add smoke-workspace pwd --workspace "$WORKSPACE"
workspace_output=$("$BIN" smoke-workspace | tail -n 1)
test "$workspace_output" = "$WORKSPACE"

"$BIN" add smoke-fail false
"$BIN" add smoke-marker touch "$MARKER"
"$BIN" chain add smoke-stop smoke-fail! smoke-marker
"$BIN" chain smoke-stop >/dev/null 2>&1 || true
test ! -e "$MARKER"

"$BIN" add smoke-bg echo "detached test"
"$BIN" smoke-bg -d
sleep 1
"$BIN" bg clean

"$BIN" add smoke-bg-attach sleep 1
"$BIN" smoke-bg-attach -d
attach_out=$("$BIN" bg attach smoke-bg-attach)
case "$attach_out" in
  *"Attaching to"* ) ;;
  * ) echo "unexpected attach output: $attach_out"; exit 1 ;;
esac
"$BIN" bg clean

echo "smoke tests passed"
