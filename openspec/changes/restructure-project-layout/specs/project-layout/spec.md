# Project Layout

## Purpose
Modern Go project structure separating CLI application, public library API, and internal implementation following Go best practices.

## ADDED Requirements

### Requirement: CMD Directory
The system SHALL use `cmd/` directory for CLI application entry point.

#### Scenario: Single binary in cmd/filesprawl
- **WHEN** the project has one main binary
- **THEN** place it in `cmd/filesprawl/` directory

#### Scenario: Main package in cmd
- **WHEN** building the application
- **THEN** `cmd/filesprawl/main.go` contains the main package

#### Scenario: Cobra commands in cmd
- **WHEN** implementing CLI commands
- **THEN** place command definitions in `cmd/filesprawl/` alongside main.go

#### Scenario: Build from cmd directory
- **WHEN** running `go build ./cmd/filesprawl`
- **THEN** produce the filesprawl binary

### Requirement: Public Library API
The system SHALL provide a public library API at the module root.

#### Scenario: Library entrypoint at root
- **WHEN** library users import the package
- **THEN** they use `import "github.com/haggishunk/filesprawl"`

#### Scenario: Exported types and functions
- **WHEN** defining public API in filesprawl.go
- **THEN** export high-level functions for common operations

#### Scenario: Clean import paths
- **WHEN** library users access functionality
- **THEN** they use `filesprawl.NewDuplicateDetector()`, not `pkg.NewDuplicateDetector()`

#### Scenario: Documentation at root
- **WHEN** library users run `go doc`
- **THEN** display package documentation from root-level filesprawl.go

### Requirement: Internal Packages
The system SHALL keep private implementation in `internal/` directory.

#### Scenario: Compiler enforcement
- **WHEN** external code tries to import `internal/` packages
- **THEN** Go compiler prevents the import

#### Scenario: Business logic in internal
- **WHEN** implementing core functionality
- **THEN** place it in appropriate `internal/` subdirectories

#### Scenario: Existing internal packages
- **WHEN** restructuring
- **THEN** keep existing internal packages (analysis, database, object, operation, rclone, remote, repository)

#### Scenario: CLI utilities in internal/cli
- **WHEN** implementing CLI-specific code
- **THEN** place it in `internal/cli/` package

### Requirement: No PKG Directory
The system SHALL NOT use `pkg/` directory.

#### Scenario: Public packages at root
- **WHEN** exposing library functionality
- **THEN** place public packages at module root, not in `pkg/`

#### Scenario: Clean import paths without pkg
- **WHEN** users import the library
- **THEN** import path is `github.com/haggishunk/filesprawl`, not `.../pkg/filesprawl`

### Requirement: Examples Directory
The system SHALL provide example programs demonstrating library usage.

#### Scenario: Examples in dedicated directory
- **WHEN** providing usage examples
- **THEN** place them in `examples/` directory

#### Scenario: Runnable examples
- **WHEN** examples are provided
- **THEN** each example is a complete, runnable program

#### Scenario: Example imports library
- **WHEN** example programs run
- **THEN** they import and use the public library API

### Requirement: CLI Layer Separation
The system SHALL separate CLI concerns from business logic.

#### Scenario: Thin CLI wrappers
- **WHEN** implementing commands in cmd/
- **THEN** they are thin wrappers calling library functions

#### Scenario: CLI-specific code isolated
- **WHEN** handling CLI concerns (flags, output formatting, env vars)
- **THEN** keep them in `cmd/` or `internal/cli/`, not in business logic

#### Scenario: Business logic reusable
- **WHEN** implementing core functionality
- **THEN** it can be used from both CLI and library API

### Requirement: Backward Compatibility
The system SHALL maintain backward compatibility during migration.

#### Scenario: Root main.go as symlink
- **WHEN** migrating to new structure
- **THEN** keep `main.go` at root as symlink to `cmd/filesprawl/main.go`

#### Scenario: go build still works
- **WHEN** running `go build` in root directory
- **THEN** build succeeds and produces binary

#### Scenario: Existing commands unchanged
- **WHEN** users run CLI commands
- **THEN** behavior is identical to before restructuring

### Requirement: Build and Installation
The system SHALL support standard Go build and install patterns.

#### Scenario: Build with go build
- **WHEN** running `go build ./cmd/filesprawl`
- **THEN** produce filesprawl binary in current directory

#### Scenario: Install with go install
- **WHEN** running `go install ./cmd/filesprawl`
- **THEN** install binary to $GOPATH/bin or $GOBIN

#### Scenario: Build from Makefile
- **WHEN** running `make build`
- **THEN** build from cmd/filesprawl to build/ directory

#### Scenario: Cross-compilation support
- **WHEN** building for different platforms
- **THEN** support GOOS and GOARCH environment variables

### Requirement: Documentation
The system SHALL document both CLI and library usage.

#### Scenario: Library godoc
- **WHEN** library users run `go doc github.com/haggishunk/filesprawl`
- **THEN** display package documentation with usage examples

#### Scenario: README covers both uses
- **WHEN** users read README.md
- **THEN** see documentation for both CLI usage and library usage

#### Scenario: Example programs documented
- **WHEN** users explore examples/
- **THEN** each example has clear comments explaining usage

### Requirement: Module Structure
The system SHALL maintain single module structure.

#### Scenario: Single go.mod
- **WHEN** managing dependencies
- **THEN** use one go.mod at repository root

#### Scenario: Internal packages importable within module
- **WHEN** cmd/ code imports internal/ packages
- **THEN** imports work (same module)

#### Scenario: No nested modules
- **WHEN** structuring the project
- **THEN** avoid creating go.mod in subdirectories

