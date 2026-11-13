#!/bin/sh
# Install a verified Spade release without Homebrew or a Go compiler.
set -eu

fail() {
    printf 'Error: %s\n' "$*" >&2
    exit 1
}

usage() {
    cat <<'USAGE'
Usage: sh installer.sh [--version VERSION] [--install-dir DIRECTORY]

Installs the latest stable Spade release to ~/.local/bin by default.
  --version VERSION        Install a specific stable version (e.g. v0.1.1)
  --install-dir DIRECTORY  Choose the directory for the spd executable
  -h, --help               Show this help
USAGE
}

main() {
    spade_version=latest
    spade_install_dir=${SPADE_INSTALL_DIR:-${HOME:?HOME must be set}/.local/bin}
    while [ "$#" -gt 0 ]; do
        case "$1" in
            --version|--install-dir)
                [ "$#" -ge 2 ] && [ -n "$2" ] || fail "$1 requires a value"
                case "$1" in
                    --version) spade_version=$2 ;;
                    --install-dir) spade_install_dir=$2 ;;
                esac
                shift 2
                ;;
            -h|--help) usage; return ;;
            *) fail "Unknown argument: $1 (use --help)" ;;
        esac
    done

    case "$(uname -s)" in
        Darwin) spade_os=darwin ;;
        Linux) spade_os=linux ;;
        *) fail "Supported operating systems are macOS and Linux" ;;
    esac
    case "$(uname -m)" in
        x86_64|amd64) spade_arch=amd64 ;;
        arm64|aarch64) spade_arch=arm64 ;;
        *) fail "Supported architectures are amd64 and arm64" ;;
    esac
    for spade_tool in curl tar install mktemp; do
        command -v "$spade_tool" >/dev/null 2>&1 || fail "Required command not found: $spade_tool"
    done
    if command -v sha256sum >/dev/null 2>&1; then
        spade_hash_tool=sha256sum
    elif command -v shasum >/dev/null 2>&1; then
        spade_hash_tool=shasum
    else
        fail "Install sha256sum or shasum to verify downloads"
    fi

    spade_repo=https://github.com/notenderdreams/spade
    if [ "$spade_version" = latest ]; then
        spade_latest_url=$(curl -fsSL --retry 3 --proto '=https' -o /dev/null -w '%{url_effective}' "$spade_repo/releases/latest") || fail "Could not resolve the latest release"
        spade_version=${spade_latest_url##*/}
    fi
    spade_version=${spade_version#v}
    printf '%s\n' "$spade_version" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' || fail "Version must be a stable version such as v0.1.1"

    spade_archive="spade_${spade_version}_${spade_os}_${spade_arch}.tar.gz"
    spade_base="$spade_repo/releases/download/v$spade_version"
    spade_tmp=$(mktemp -d)
    spade_staged=
    trap 'rm -rf "$spade_tmp"; if [ -n "$spade_staged" ]; then rm -f "$spade_staged"; fi' EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM

    printf 'Downloading Spade v%s for %s/%s...\n' "$spade_version" "$spade_os" "$spade_arch"
    curl -fsSL --retry 3 --proto '=https' -o "$spade_tmp/$spade_archive" "$spade_base/$spade_archive" || fail "Could not download the release archive"
    curl -fsSL --retry 3 --proto '=https' -o "$spade_tmp/checksums.txt" "$spade_base/spade_${spade_version}_checksums.txt" || fail "Could not download the release checksums"
    spade_expected=$(awk -v name="$spade_archive" '$2 == name {print $1}' "$spade_tmp/checksums.txt")
    printf '%s\n' "$spade_expected" | grep -Eq '^[[:xdigit:]]{64}$' || fail "Missing or invalid archive checksum"
    if [ "$spade_hash_tool" = sha256sum ]; then
        spade_actual=$(sha256sum "$spade_tmp/$spade_archive" | awk '{print $1}')
    else
        spade_actual=$(shasum -a 256 "$spade_tmp/$spade_archive" | awk '{print $1}')
    fi
    [ "$spade_actual" = "$spade_expected" ] || fail "Checksum mismatch; nothing was installed"

    tar -xzf "$spade_tmp/$spade_archive" -C "$spade_tmp" spd
    [ -f "$spade_tmp/spd" ] && [ ! -L "$spade_tmp/spd" ] || fail "Archive does not contain a regular spd executable"
    mkdir -p "$spade_install_dir"
    spade_staged=$(mktemp "$spade_install_dir/.spd.XXXXXX")
    install -m 755 "$spade_tmp/spd" "$spade_staged"
    mv -f "$spade_staged" "$spade_install_dir/spd"
    spade_staged=
    printf 'Installed Spade v%s to %s/spd\n' "$spade_version" "$spade_install_dir"
    case ":${PATH:-}:" in
        *":$spade_install_dir:"*) ;;
        *) printf 'Add %s to your PATH to run spd from any directory.\n' "$spade_install_dir" ;;
    esac
}

main "$@"
