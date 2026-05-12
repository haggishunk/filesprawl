# Change: Add Release Pipeline

## Why
The project lacks a reproducible release process. Builds have no embedded version, the Makefile only produces a single local binary, and there is no automation for publishing tagged releases. This makes it impossible for users to know which build they are running, blocks publishing pre-built binaries for common platforms, and forces every consumer to build from source. A release pipeline is needed to:

- Embed a build-time version into the binary so `--version` is meaningful in the field.
- Produce small, stripped release artifacts via the Makefile so local and CI builds are byte-identical.
- Publish multi-platform binaries with checksums on every `v*` git tag.

## What Changes
- Introduce a `VERSION` variable in the Makefile that defaults to `git describe --tags --always --dirty` and may be overridden via environment or CLI (`make VERSION=v1.2.3 build`).
- Strip a leading `v` from `VERSION` when injecting it into the binary so the embedded value is bare semver, matching standard release-artifact conventions.
- Inject the version into the CLI via `-ldflags '-X main.version=...'` and surface it through Cobra's `--version` flag.
- Strip symbol tables and DWARF info from release builds (`-s -w`) to reduce binary size.
- Add a `make install` target that uses the same ldflags so `go install`-style installations also pick up version metadata.
- Add a GitHub Actions workflow (`.github/workflows/release.yml`) that triggers on `v*` tags, builds a matrix of `linux/{amd64,arm64}` and `darwin/{amd64,arm64}` binaries via `make build`, generates SHA-256 checksums, and publishes them to a GitHub Release with auto-generated notes.

## Impact
- Affected specs:
  - New capability `release-pipeline`
- Affected code:
  - `Makefile` — `VERSION`, `GO_VERSION`, `LDFLAGS`, `build`, `install`, `.PHONY`
  - `cmd/filesprawl/root.go` — `var version` and `rootCmd.Version`
  - `cmd/filesprawl/cmd_test.go` — regression test asserting `rootCmd.Version` is wired
  - `.github/workflows/release.yml` — new release workflow
- Users can now:
  - Run `filesprawl --version` and see the actual release tag they installed.
  - Download pre-built, checksummed binaries for Linux and macOS on amd64 and arm64 from GitHub Releases.
  - Reproduce the published binary locally with `make VERSION=vX.Y.Z build`.
