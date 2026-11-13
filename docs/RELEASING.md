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
   SPADE_TEST_BIN=/absolute/path/to/spd SPADE_TEST_VERSION=0.1.0 sh test/smoke.sh
   ```

5. Commit the release preparation changes and check that the worktree is clean.

## Publish

After reviewing the prepared commit, push the branch and an annotated tag:

```sh
git push origin main
git tag -a v0.1.0 -m 'Spade v0.1.0'
git push origin v0.1.0
```

The tag workflow runs the shared CI checks on Linux and macOS, builds snapshot
packages for all supported targets, and only then publishes the GitHub release.
It uses the repository's built-in `GITHUB_TOKEN` with `contents: write` for the
publishing job. No separate release secret is required. GitHub Actions must be
enabled for the repository.

The release includes four versioned `.tar.gz` archives (Linux/macOS ×
amd64/arm64) and `spade_VERSION_checksums.txt` with SHA-256 hashes. Windows is
excluded because process management uses Unix APIs. macOS binaries are unsigned
and are not notarized. Prerelease tags are marked as prereleases automatically.

Confirm the Actions run succeeds, download the native archive, verify its
checksum, and run `spd --version` and the smoke suite against the downloaded
binary. The CLI version must match the tag without its leading `v`.

If publishing fails, inspect the failed job before retrying it. Do not move a
published release tag; fix the problem in a new version if necessary.
