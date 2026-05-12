# Design: Restructure Project Layout

## Overview
Restructure filesprawl to follow modern Go best practices with clear separation between CLI application and reusable library, enabling both command-line and programmatic usage.

## Current Structure (Before)

```
filesprawl/
├── main.go                  # 168 lines, CLI + business logic mixed
├── main_test.go             # Tests for main.go
├── internal/                # Business logic packages
│   ├── analysis/
│   ├── database/
│   ├── object/
│   ├── operation/
│   ├── rclone/
│   ├── remote/
│   └── repository/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

**Problems:**
- All CLI code in one file
- Business logic not usable as library
- No separation of concerns
- Cannot import filesprawl from other programs

## Proposed Structure (After)

```
filesprawl/
├── filesprawl.go            # NEW: Public library API
│
├── cmd/                     # NEW: CLI applications
│   └── filesprawl/
│       ├── main.go          # CLI entry point
│       ├── root.go          # Cobra root command
│       ├── list_remotes.go  # list-remotes command
│       ├── index.go         # index command
│       └── duplicates.go    # duplicates command
│
├── internal/                # Private packages
│   ├── analysis/            # Existing: duplicate detection
│   ├── database/            # Existing: database interface
│   ├── object/              # Existing: domain models
│   ├── operation/           # Existing: scanning operations
│   ├── rclone/              # Existing: rclone RC client
│   ├── remote/              # Existing: remote config
│   ├── repository/          # Existing: data persistence
│   └── cli/                 # NEW: CLI utilities
│       ├── config.go        # Database connection, env vars
│       └── output.go        # Output formatting
│
├── examples/                # NEW: Library usage examples
│   ├── find_duplicates/
│   │   └── main.go
│   └── README.md
│
├── main.go                  # Symlink to cmd/filesprawl/main.go (backward compat)
├── go.mod
├── go.sum
├── Makefile                 # Updated: build from cmd/
└── README.md                # Updated: add library section
```

## Detailed Design

### 1. Public Library API (`filesprawl.go`)

```go
// Package filesprawl provides file indexing and duplicate detection
// across multiple remote storage locations.
//
// Example usage:
//
//	import "github.com/haggishunk/filesprawl"
//
//	// Create repository
//	repo := filesprawl.NewRepository(db)
//
//	// Find duplicates across remotes
//	detector := filesprawl.NewDuplicateDetector(repo)
//	groups, err := detector.FindAcrossRemotes(ctx, filesprawl.FilterOptions{
//		HashType: "sha256",
//		MinSize:  1048576, // 1MB
//	})
package filesprawl

import (
	"context"
	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/repository"
)

// DuplicateDetector finds duplicate files using content hashes.
// Re-exported from internal/analysis for public use.
type DuplicateDetector = analysis.DuplicateDetector

// FilterOptions controls which duplicates are returned.
type FilterOptions = analysis.FilterOptions

// DuplicateGroup is a set of files that share the same content hash.
type DuplicateGroup = analysis.DuplicateGroup

// FileInfo holds per-file detail within a duplicate group.
type FileInfo = analysis.FileInfo

// NewDuplicateDetector creates a detector backed by the given repository.
func NewDuplicateDetector(repo *repository.ObjectRepository) *DuplicateDetector {
	return analysis.NewDuplicateDetector(repo)
}

// FormatDuplicateReport returns a human-readable report of duplicate groups.
func FormatDuplicateReport(groups []DuplicateGroup) string {
	return analysis.FormatReport(groups)
}

// Additional high-level convenience functions...
```

### 2. CLI Entry Point (`cmd/filesprawl/main.go`)

```go
package main

import (
	"os"
	"github.com/haggishunk/filesprawl/cmd/filesprawl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

### 3. Cobra Root Command (`cmd/filesprawl/root.go`)

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filesprawl",
	Short: "File indexing and duplicate detection across remote storage",
	Long: `Filesprawl indexes and analyzes files across multiple remote 
storage locations, identifying duplicates and optimizing storage usage.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add persistent flags here if needed
	// rootCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "PostgreSQL connection string")
}
```

### 4. Index Command (`cmd/filesprawl/index.go`)

```go
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index files from a remote storage location",
	Long:  `Index all files from the specified remote, computing hashes and persisting metadata.`,
	Example: `  # Index entire remote
  filesprawl index --remote media

  # Index specific path
  filesprawl index --remote media --path code/flux`,
	RunE: runIndex,
}

var (
	indexRemote string
	indexPath   string
)

func init() {
	rootCmd.AddCommand(indexCmd)
	
	indexCmd.Flags().StringVarP(&indexRemote, "remote", "r", "", "remote name (required)")
	indexCmd.Flags().StringVarP(&indexPath, "path", "p", "", "path within remote")
	indexCmd.MarkFlagRequired("remote")
}

func runIndex(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	
	// Open repository using CLI utility
	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()
	
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to determine hostname: %w", err)
	}
	
	// Normalize remote name
	remoteName := strings.TrimSuffix(strings.TrimSpace(indexRemote), ":")
	rem := remote.NewRemote(hostname, remoteName)
	
	// Create scanner
	scn := operation.NewScanner(
		operation.WithRemote(&rem),
		operation.WithRepo(repo),
	)
	
	// Configure rclone list
	lo := rclone.NewListOption(rclone.ListOptionFilesOnly())
	lc := rclone.NewListConfig(remoteName+":", indexPath, &lo)
	
	// Perform scan
	if err := operation.Scan(ctx, &scn, lc); err != nil {
		return err
	}
	
	fmt.Fprintf(cmd.OutOrStdout(), "indexed remote %s from host %s\n", rem.Name, rem.Hostname)
	return nil
}
```

### 5. CLI Utilities (`internal/cli/config.go`)

```go
package cli

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenRepository opens a database connection and returns an ObjectRepository.
// The caller must call the cleanup function when done.
func OpenRepository(ctx context.Context) (*repository.ObjectRepository, func(), error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing database url: %w", err)
	}
	
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		log.Printf("Connected to database with pid %d", conn.PgConn().PID())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}

	return repository.NewObjectRepository(database.NewPgxDatabase(pool)), cleanup, nil
}
```

### 6. Example Program (`examples/find_duplicates/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haggishunk/filesprawl"
	"github.com/haggishunk/filesprawl/internal/cli"
)

func main() {
	ctx := context.Background()

	// Open repository (using CLI utility for convenience)
	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Create duplicate detector
	detector := filesprawl.NewDuplicateDetector(repo)

	// Find duplicates across all remotes
	groups, err := detector.FindAcrossRemotes(ctx, filesprawl.FilterOptions{
		HashType: "sha256",
		MinSize:  10485760, // 10MB
		Limit:    10,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Format and display report
	report := filesprawl.FormatDuplicateReport(groups)
	fmt.Println(report)
}
```

## Build and Installation

### Building

```bash
# Build from cmd directory
go build -o build/filesprawl ./cmd/filesprawl

# Or use Makefile
make build

# Or build from root (via symlink - backward compat)
go build
```

### Installing

```bash
# Install to $GOPATH/bin or $GOBIN
go install ./cmd/filesprawl

# Now available as 'filesprawl' in PATH
filesprawl --help
```

### Using as Library

```bash
# Import in go.mod
go get github.com/haggishunk/filesprawl

# Use in code
import "github.com/haggishunk/filesprawl"
```

## Migration Strategy

### Phase 1: Structure (No Breaking Changes)
1. Create directory structure
2. Copy main.go to cmd/filesprawl/main.go
3. Create symlink at root
4. Update Makefile
5. Test everything still works

### Phase 2: Extract and Refactor
1. Create internal/cli package
2. Move CLI utilities to internal/cli
3. Split commands into separate files
4. Create public API in filesprawl.go
5. Tests still pass

### Phase 3: Enhance
1. Add Cobra commands (if combined with migrate-to-cobra-cli)
2. Create example programs
3. Document library API
4. Add godoc examples

### Phase 4: Finalize
1. Deprecate root main.go (keep symlink or remove)
2. Update all documentation
3. Publish release with new structure

## Benefits

### For CLI Users
- No breaking changes
- Better organized codebase
- Eventually: better help, completion (with Cobra)

### For Library Users
- Clean import: `import "github.com/haggishunk/filesprawl"`
- Professional API
- Reusable components
- Example programs to learn from

### For Developers
- Clear separation of concerns
- Each command in its own file
- Easy to add new commands
- Better testability
- Standard Go project structure

## References

- [Official Go Module Layout](https://go.dev/doc/modules/layout)
- [No nonsense guide to Go projects layout](https://laurentsv.com/blog/2024/10/19/no-nonsense-go-package-layout.html)
- [Go Project Layout: What Actually Works](https://alnah.io/post/go-project-layout/)
- Cobra CLI framework: https://github.com/spf13/cobra

