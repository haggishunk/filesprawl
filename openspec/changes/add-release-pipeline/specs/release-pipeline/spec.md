# Release Pipeline

## Purpose
Reproducible build and release process that embeds version metadata into the `filesprawl` binary, produces stripped artifacts via the Makefile, and publishes multi-platform binaries with checksums on every `v*` git tag.

## ADDED Requirements

### Requirement: Version Override via Environment
The Makefile SHALL allow the release version to be supplied externally and SHALL fall back to `git describe` when none is provided.

#### Scenario: Explicit version on the command line
- **WHEN** running `make VERSION=v1.2.3 build`
- **THEN** the resulting binary embeds `1.2.3` and `build/filesprawl --version` prints `filesprawl version 1.2.3`

#### Scenario: Bare semver passes through unchanged
- **WHEN** running `make VERSION=2.0.0 build`
- **THEN** the embedded version is `2.0.0` (no transformation applied beyond the `v` strip)

#### Scenario: Local fallback when VERSION is unset
- **WHEN** running `make build` with no `VERSION` in the environment
- **THEN** the embedded version is the output of `git describe --tags --always --dirty`, or `dev` if git is unavailable

### Requirement: Leading v Stripped for Embedded Version
The Makefile SHALL strip a single leading `v` from `VERSION` before injecting it into the Go binary so that embedded versions are bare semver.

#### Scenario: v-prefixed git tag
- **WHEN** `VERSION=v1.2.3` is supplied
- **THEN** the value passed to `-X main.version=...` is `1.2.3`

#### Scenario: Already-bare version
- **WHEN** `VERSION=1.2.3` is supplied
- **THEN** the value passed to `-X main.version=...` is `1.2.3` (unchanged)

### Requirement: Stripped Release Binaries
Release builds produced by the Makefile SHALL omit the symbol table and DWARF debug information.

#### Scenario: ldflags include strip directives
- **WHEN** the `build` or `install` target runs
- **THEN** the Go linker is invoked with `-s -w` in addition to the version injection

### Requirement: Install Target
The Makefile SHALL provide an `install` target that installs the CLI with the same version metadata as `build`.

#### Scenario: make install uses the same ldflags as make build
- **WHEN** running `make install`
- **THEN** `go install` is invoked against `./cmd/filesprawl` with the same `-ldflags` value used by the `build` target

### Requirement: CLI Version Flag
The `filesprawl` CLI SHALL surface the embedded version through Cobra's `--version` flag.

#### Scenario: --version reports the embedded value
- **WHEN** running `filesprawl --version`
- **THEN** the CLI prints `filesprawl version <embedded-version>` and exits 0

#### Scenario: Default version when not injected
- **WHEN** the binary is built without `-ldflags '-X main.version=...'` (e.g. `go build`)
- **THEN** `--version` reports `dev`

### Requirement: Tagged Release Workflow
The repository SHALL include a GitHub Actions workflow that publishes a multi-platform GitHub Release for every `v*` git tag.

#### Scenario: Triggered by v-prefixed tag
- **WHEN** a tag matching `v*` is pushed to the repository
- **THEN** the `release` workflow runs

#### Scenario: Build matrix covers Linux and macOS on amd64 and arm64
- **WHEN** the workflow's build job runs
- **THEN** binaries are produced for `linux/amd64`, `linux/arm64`, `darwin/amd64`, and `darwin/arm64`

#### Scenario: Tag flows into the Makefile as VERSION
- **WHEN** the workflow builds a binary
- **THEN** `VERSION=${{ github.ref_name }}` is exported so `make build` embeds the v-prefixed tag (with the leading `v` stripped at injection time)

#### Scenario: Artifacts named per platform
- **WHEN** the workflow uploads a build
- **THEN** the artifact filename is `filesprawl-${GOOS}-${GOARCH}`

### Requirement: Release Checksums
The release workflow SHALL publish SHA-256 checksums for every binary it ships.

#### Scenario: checksums.txt accompanies the binaries
- **WHEN** the release job runs
- **THEN** a `checksums.txt` file produced via `shasum -a 256 filesprawl-*` is uploaded alongside the binaries

#### Scenario: GitHub release attaches binaries and checksums
- **WHEN** the release is published
- **THEN** the GitHub Release contains every `filesprawl-*` artifact plus `checksums.txt` and uses auto-generated release notes
