# Releasing Spade

Releases use semantic version tags (`vMAJOR.MINOR.PATCH`). The tag is the
source of the packaged CLI version. `cmd.Version` is the fallback for local
source builds; update it when preparing a new version.

## Prepare

1. Update `cmd.Version`, `CHANGELOG.md`, and versioned installation examples in
   `README.md`.
2. Write `docs/releases/vX.Y.Z.md`. The publishing workflow requires this file
   for the exact tag being released.
3. Run the checks below with the Go version specified by `go.mod` or later:

   ```sh
   go mod tidy -diff
   go mod verify
   go vet ./...
   go test -race ./...
   sh test/smoke.sh
   python3 test/installer_test.py
   ```

4. Install GoReleaser **v2.20.2**, matching the workflows, then validate all four
   platform builds without uploading anything:

   ```sh
   goreleaser check
   goreleaser release --snapshot --clean
   ```

   Snapshot packages are written under the ignored `dist/` directory. They
   contain the `spd` binary and documentation. The version inside a snapshot
   differs from the final tagged release. The smoke suite can exercise a
   prebuilt native binary:

   ```sh
   SPADE_TEST_BIN=/absolute/path/to/spd SPADE_TEST_VERSION=0.1.1 sh test/smoke.sh
   ```

5. Commit the release preparation changes and check that the worktree is clean.

## Publish

After reviewing the prepared commit, push the branch and an annotated tag:

```sh
git push origin main
git tag -a v0.1.1 -m 'Spade v0.1.1'
git push origin v0.1.1
```

The tag workflow runs the shared CI checks on Linux and macOS, builds snapshot
packages for all supported targets, and only then publishes the GitHub release.
It uses the repository's built-in `GITHUB_TOKEN` with `contents: write` for the
publishing job. No separate release secret is required for GitHub releases.
GitHub Actions must be enabled for the repository.

## Homebrew tap

After a stable release is published, the workflow updates
[`notenderdreams/homebrew-tap`](https://github.com/notenderdreams/homebrew-tap)
with a generated binary formula at `Formula/spade.rb`. The identical `spade.rb`
is uploaded as a GitHub release asset, matching Dryve's generate-and-publish flow.
Prereleases do not update the tap.
Users install with `brew install notenderdreams/tap/spade`.

The Spade repository must have the Actions secret `HOMEBREW_TAP_TOKEN`. Use a
fine-grained GitHub token with **Contents: Read and write** access to the
`homebrew-tap` repository. The default `GITHUB_TOKEN` cannot write to the tap.
The token is used only when checking out the tap for publication; it is not
available to pull request checks.

The formula selects the prebuilt archive for macOS/Linux and Intel/ARM64, pins
its SHA-256 checksum, and installs `spd` plus shell completions. No Go compiler
is required on the user's machine. This matches the binary-formula approach used
by Dryve in the same tap.

Published releases run on a macOS runner. A post-build hook ad-hoc signs and
verifies both macOS architectures before archiving and checksumming them. Local
and CI snapshot builds on Linux skip this hook. Ad-hoc signing does not identify
a trusted developer and does not provide notarization. Developer ID signing
and notarization require Apple Developer Program membership. The original
v0.1.0 assets are unchanged; its ARM64 binary already has Go's linker-generated
ad-hoc signature, while its Intel binary has no signature.

The template is `packaging/homebrew/spade.rb.in`. The release action renders it
using the checksums of that release's final archives, uploads `spade.rb`, and
copies the same file into the tap before committing and pushing it. No manual
formula edits or tap commits are required.

The reusable `Homebrew` workflow tests installation and completions with
Homebrew before publishing. It also accepts a manual dispatch with an existing
release tag, so formula fixes can be published without changing release binaries
or moving a tag. If a tap update fails, inspect the Actions log and rerun the
Homebrew workflow after correcting the problem.

## Verify published artifacts

The release includes four versioned `.tar.gz` archives (Linux/macOS ×
amd64/arm64), `installer.sh`, and `spade_VERSION_checksums.txt` with SHA-256 hashes. Windows is
excluded because process management uses Unix APIs. macOS binaries have no
Developer ID signature and are not notarized. Prerelease tags are marked as prereleases automatically.

The standalone `installer.sh` detects the platform, downloads the selected
stable release, verifies its archive against the checksum manifest, and installs
`spd` to `~/.local/bin` by default. It accepts `--version` and `--install-dir`.
The README uses the script on `main`, so it also works with existing releases
that predate the installer release asset.

Confirm the Actions run succeeds, download the native archive, verify its
checksum, and run `spd --version` and the smoke suite against the downloaded
binary. The CLI version must match the tag without its leading `v`.

If publishing fails, inspect the failed job before retrying it. Do not move a
published release tag; fix the problem in a new version if necessary.
