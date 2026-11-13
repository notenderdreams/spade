#!/usr/bin/env python3
"""Generate spade.rb from the release's SHA-256 manifest."""

import argparse
from pathlib import Path
import re


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True)
    parser.add_argument("--checksums", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    version = args.version.removeprefix("v")
    if not re.fullmatch(r"(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)", version):
        parser.error("version must be a stable semantic version, e.g. v0.1.0")
    checksums = {}
    for line in args.checksums.read_text().splitlines():
        if not line.strip():
            continue
        checksum, filename = line.split()
        if not re.fullmatch(r"[0-9a-fA-F]{64}", checksum):
            parser.error(f"invalid SHA-256 checksum for {filename}")
        if filename in checksums:
            parser.error(f"duplicate checksum for {filename}")
        checksums[filename] = checksum.lower()
    template = Path(__file__).resolve().parents[1] / "packaging/homebrew/spade.rb.in"
    formula = template.read_text().replace("@VERSION@", version)
    for os_name in ("darwin", "linux"):
        for arch in ("amd64", "arm64"):
            name = f"spade_{version}_{os_name}_{arch}.tar.gz"
            if name not in checksums:
                parser.error(f"missing checksum for {name}")
            formula = formula.replace(f"@{os_name.upper()}_{arch.upper()}_SHA256@", checksums[name])
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(formula)
    print(f"Generated {args.output} for v{version}")


if __name__ == "__main__":
    main()
