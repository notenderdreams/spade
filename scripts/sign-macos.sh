#!/bin/sh
set -eu

# Linux snapshot builds cannot run codesign. Published releases run on macOS.
[ "$1" = darwin ] || exit 0
if [ "$(uname -s)" != Darwin ]; then
    echo "Skipping macOS ad-hoc signing on this snapshot build host"
    exit 0
fi

codesign --force --sign - --identifier dev.notenderdreams.spade --timestamp=none "$2"
codesign --verify --strict "$2"
