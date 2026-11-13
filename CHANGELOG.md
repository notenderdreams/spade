# Changelog

## 0.1.1

- Install prebuilt macOS and Linux binaries with `brew install notenderdreams/tap/spade`.
- Generate `spade.rb` during releases, upload it as a release asset, and automatically publish it to the Homebrew tap.
- Install Bash, Zsh, and Fish completions through Homebrew.
- Ad-hoc sign and verify both Intel and Apple Silicon macOS binaries before packaging.

See [the release notes](docs/releases/v0.1.1.md) for details.

## 0.1.0

Initial release of Spade (`spd`).

- Reusable command shortcuts with argument templates and defaults.
- Script management, import/export, and local script backups.
- Sequential chains, pipelines, workspace selection, and stop-on-error support.
- Detached execution and background process listing, stopping, log attachment, and cleanup.
- Bash, Zsh, and Fish shell completion.
- Linux and macOS binaries for amd64 and arm64, with SHA-256 checksums.

See [the release notes](docs/releases/v0.1.0.md) for details.
