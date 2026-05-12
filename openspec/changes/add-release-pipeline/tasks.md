# Implementation Tasks

## 1. Makefile Version Plumbing
- [x] 1.1 Add `VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)` to Makefile
- [x] 1.2 Add `GO_VERSION := $(VERSION:v%=%)` to strip a leading `v` for the embedded value
- [x] 1.3 Define `LDFLAGS := -s -w -X main.version=$(GO_VERSION)`
- [x] 1.4 Update `build` target to pass `-ldflags '$(LDFLAGS)'` and emit `build/filesprawl`
- [x] 1.5 Mark `build` as `.PHONY` so it always rebuilds despite the existing `build/` directory

## 2. Install Target
- [x] 2.1 Add `install` target that runs `go install -ldflags '$(LDFLAGS)' ./cmd/filesprawl`
- [x] 2.2 Add `install` to `.PHONY`
- [x] 2.3 Verify with `make -n install` that the recipe expands correctly

## 3. CLI Version Flag
- [x] 3.1 Add `var version = "dev"` to `cmd/filesprawl/root.go`
- [x] 3.2 Set `rootCmd.Version = version` so Cobra serves `--version`
- [x] 3.3 Add a regression test in `cmd/filesprawl/cmd_test.go` asserting `rootCmd.Version` is non-empty

## 4. Release Workflow
- [x] 4.1 Create `.github/workflows/release.yml` triggered on `push.tags: ['v*']`
- [x] 4.2 Grant `contents: write` permission for the release job
- [x] 4.3 Define a build matrix for linux/{amd64,arm64} and darwin/{amd64,arm64}
- [x] 4.4 Set `VERSION: ${{ github.ref_name }}` so the v-prefixed tag flows into the Makefile
- [x] 4.5 Build via `make build` and rename the artifact to `filesprawl-${GOOS}-${GOARCH}`
- [x] 4.6 Upload each artifact with `actions/upload-artifact@v4`
- [x] 4.7 In a downstream `release` job, download all artifacts and emit `checksums.txt` via `shasum -a 256`
- [x] 4.8 Publish a GitHub Release with all binaries plus checksums and `generate_release_notes: true`

## 5. Verification
- [x] 5.1 `make VERSION=v1.2.3 build` and confirm `build/filesprawl --version` reports `1.2.3`
- [x] 5.2 `make VERSION=2.0.0 build` and confirm bare semver passes through unchanged
- [x] 5.3 `env -u VERSION make build` and confirm the `git describe` fallback is embedded
- [x] 5.4 `go test ./...` and `go vet ./...` pass after the changes

## 6. Validation
- [ ] 6.1 Run `openspec validate add-release-pipeline --strict`
