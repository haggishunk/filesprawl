# Implementation Tasks

## 1. Create Directory Structure
- [x] 1.1 Create `cmd/filesprawl/` directory
- [x] 1.2 Create `internal/cli/` directory
- [ ] 1.3 Create `examples/` directory
- [x] 1.4 Verify existing `internal/` subdirectories remain unchanged

## 2. Move CLI Entry Point
- [x] 2.1 Copy `main.go` to `cmd/filesprawl/main.go`
- [x] 2.2 Update package declaration to `package main`
- [x] 2.3 Update imports in cmd/filesprawl/main.go (add github.com/haggishunk/filesprawl/internal/...)
- [x] 2.4 Test build with `go build ./cmd/filesprawl`
- [-] 2.5 Replace root `main.go` with symlink: `ln -s cmd/filesprawl/main.go main.go` (superseded: root main.go removed entirely for cross-platform cleanliness)

## 3. Create Cobra Command Structure (if combined with migrate-to-cobra-cli)
- [x] 3.1 Create `cmd/filesprawl/root.go` with root command
- [-] 3.2 Create `cmd/filesprawl/list_remotes.go` for list-remotes command (list-remotes lives inside root.go)
- [x] 3.3 Create `cmd/filesprawl/index.go` for index command
- [x] 3.4 Create `cmd/filesprawl/duplicates.go` for duplicates command
- [x] 3.5 Update `cmd/filesprawl/main.go` to execute root command
- [x] 3.6 Move command handler functions to appropriate files

## 4. Extract CLI Utilities
- [x] 4.1 Create `internal/cli/config.go`
- [x] 4.2 Move `openRepository()` to `internal/cli/config.go` as `OpenRepository()`
- [x] 4.3 Add database connection helper functions to config.go
- [ ] 4.4 Create `internal/cli/output.go` for output formatting utilities
- [ ] 4.5 Add any CLI-specific helper functions to internal/cli

## 5. Create Public Library API
- [x] 5.1 Create `filesprawl.go` at repository root
- [x] 5.2 Add package documentation comment
- [x] 5.3 Export `NewDuplicateDetector()` function
- [x] 5.4 Export `DuplicateDetector` type (re-export from internal/analysis)
- [x] 5.5 Export `FilterOptions` type (re-export from internal/analysis)
- [x] 5.6 Export `DuplicateGroup` type (re-export from internal/analysis)
- [ ] 5.7 Export `Scanner` functionality for indexing
- [ ] 5.8 Add high-level convenience functions
- [x] 5.9 Document all exported types and functions

## 6. Create Example Programs
- [ ] 6.1 Create `examples/find_duplicates/main.go`
- [ ] 6.2 Implement example showing duplicate detection via library
- [ ] 6.3 Create `examples/list_files/main.go` (optional)
- [ ] 6.4 Add README.md in examples/ explaining how to run them
- [ ] 6.5 Test examples compile and run

## 7. Update Build Configuration
- [x] 7.1 Update Makefile `build` target to use `./cmd/filesprawl`
- [x] 7.2 Update Makefile to output binary to `build/filesprawl`
- [x] 7.3 Add `install` target to Makefile
- [x] 7.4 Test `make build` produces correct binary
- [-] 7.5 Test `go build` still works from root (via symlink) (n/a: root main.go removed; use `go build ./cmd/filesprawl`)
- [ ] 7.6 Test `go install ./cmd/filesprawl` installs to GOPATH/bin

## 8. Update Import Paths
- [x] 8.1 Update imports in cmd/filesprawl/*.go to reference internal packages
- [x] 8.2 Update imports to use internal/cli where appropriate
- [x] 8.3 Ensure all imports use full module path (github.com/haggishunk/filesprawl/internal/...)
- [x] 8.4 Run `go mod tidy`
- [x] 8.5 Verify no broken imports

## 9. Update Tests
- [x] 9.1 Move `main_test.go` to `cmd/filesprawl/main_test.go` if applicable
- [x] 9.2 Update test imports
- [ ] 9.3 Add tests for public library API
- [x] 9.4 Run `go test ./...` and verify all tests pass
- [ ] 9.5 Add integration tests using library API

## 10. Documentation
- [x] 10.1 Add package documentation to `filesprawl.go`
- [ ] 10.2 Update README.md with "Using as a Library" section
- [ ] 10.3 Add library usage examples to README
- [ ] 10.4 Document example programs
- [ ] 10.5 Add godoc examples (Example functions in tests)
- [ ] 10.6 Update installation instructions in README
- [ ] 10.7 Document both CLI and library usage clearly

## 11. Update Project Metadata
- [ ] 11.1 Update openspec/project.md with new structure
- [ ] 11.2 Add library API section to project.md
- [ ] 11.3 Update .gitignore if needed (ensure build/ is ignored)
- [ ] 11.4 Update any CI/CD configuration for new build path

## 12. Testing and Validation
- [-] 12.1 Test `go build` from root directory (n/a: root main.go removed)
- [x] 12.2 Test `go build ./cmd/filesprawl`
- [ ] 12.3 Test `go install ./cmd/filesprawl`
- [x] 12.4 Test `make build`
- [ ] 12.5 Run all CLI commands to verify functionality
- [ ] 12.6 Test library usage in example programs
- [x] 12.7 Run `go test ./...`
- [x] 12.8 Run `go vet ./...`
- [ ] 12.9 Run `golint ./...` or `golangci-lint run`
- [ ] 12.10 Verify godoc output looks correct

## 13. Cleanup
- [x] 13.1 Remove old main.go if symlink works correctly
- [x] 13.2 Run `go mod tidy`
- [ ] 13.3 Run `go fmt ./...`
- [x] 13.4 Ensure no unused code remains
- [ ] 13.5 Update CHANGELOG if project has one

## 14. Final Validation
- [ ] 14.1 Build and test locally with real database
- [ ] 14.2 Test all CLI commands work identically
- [ ] 14.3 Test library examples compile and run
- [x] 14.4 Verify import paths are clean
- [ ] 14.5 Run `openspec validate restructure-project-layout --strict`

