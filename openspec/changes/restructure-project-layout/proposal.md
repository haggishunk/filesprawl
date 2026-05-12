# Change: Restructure Project Layout

## Why
The current project structure has all code in a single `main.go` file (168 lines) with business logic mixed into the CLI layer. This creates several issues:

**Current Problems:**
- Business logic is not reusable as a library
- CLI code is tightly coupled to main.go
- Cannot use filesprawl functionality from other Go programs
- Testing requires testing through the CLI layer
- Adding new commands requires modifying main.go
- No clear separation between CLI, application logic, and domain logic

**Industry Standards:**
Based on official Go documentation and community consensus (2024):
- `cmd/` for multiple binaries or to separate CLI from library
- `internal/` for private packages (compiler-enforced)
- Library code at root for clean import paths (NOT in `pkg/`)
- Thin CLI wrappers around reusable library code
- Start simple, add structure when needed

## What Changes

### Proposed Structure
```
filesprawl/
├── cmd/
│   └── filesprawl/          # CLI application entry point
│       ├── main.go          # Cobra root command setup
│       ├── root.go          # Root command definition
│       ├── list_remotes.go  # list-remotes command
│       ├── index.go         # index command
│       └── duplicates.go    # duplicates command
│
├── internal/                # Private packages (existing + new)
│   ├── analysis/            # Duplicate detection logic (existing)
│   ├── database/            # Database interface (existing)
│   ├── object/              # Domain models (existing)
│   ├── operation/           # Scanning operations (existing)
│   ├── rclone/              # rclone RC client (existing)
│   ├── remote/              # Remote configuration (existing)
│   ├── repository/          # Data persistence (existing)
│   └── cli/                 # NEW: CLI-specific utilities
│       ├── config.go        # Database connection, env vars
│       └── output.go        # Output formatting helpers
│
├── filesprawl.go            # NEW: Public library API
├── examples/                # NEW: Example programs using the library
│   └── find_duplicates/
│       └── main.go
│
├── main.go                  # DEPRECATED: Symlink to cmd/filesprawl/main.go
├── go.mod
├── go.sum
├── README.md
└── ...
```

### Key Changes

1. **Extract CLI to `cmd/filesprawl/`**
   - Move command definitions from main.go to cmd/filesprawl/
   - Use Cobra command structure (combines with migrate-to-cobra-cli)
   - CLI becomes a thin wrapper around library functions

2. **Create Public Library API**
   - New `filesprawl.go` at root with exported functions
   - Clean import path: `github.com/haggishunk/filesprawl`
   - Enable programmatic use of filesprawl functionality

3. **Add `internal/cli` Package**
   - Move CLI-specific code (database connection, output formatting)
   - Keep internal/ packages focused on business logic
   - CLI utilities stay private (not part of public API)

4. **Maintain Backward Compatibility**
   - Keep `main.go` as symlink to `cmd/filesprawl/main.go`
   - Binary still built with `go build`
   - All existing commands work unchanged

## Impact

### Benefits for Library Users
```go
import "github.com/haggishunk/filesprawl"

// Use filesprawl programmatically
detector := filesprawl.NewDuplicateDetector(repo)
duplicates, err := detector.FindAcrossRemotes(ctx, opts)
```

### Benefits for CLI Development
- Each command in its own file
- Easier to add new commands
- Better testability
- Clear separation of concerns

### Benefits for Testing
- Test business logic independently of CLI
- Test CLI separately from business logic
- Integration tests can use library directly

### Affected Code
- `main.go` → moves to `cmd/filesprawl/main.go`
- Command functions → split into `cmd/filesprawl/*.go`
- Database connection → moves to `internal/cli/config.go`
- Internal packages → remain in `internal/` (no changes)
- New `filesprawl.go` → public API exports

### Migration Path
1. Phase 1: Create structure (no breaking changes)
2. Phase 2: Move code gradually
3. Phase 3: Deprecate old main.go (keep as symlink)
4. Phase 4: Document library API

