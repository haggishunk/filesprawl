# Change: Add Duplicates CLI Command

## Why
The duplicate detection backend (`internal/analysis` package) is fully implemented and tested, but there is no CLI command to expose this functionality to users. Users need a command-line interface to discover and analyze duplicate files within and across their indexed remotes, making it easy to identify storage waste and redundant files.

## What Changes
- Add a `duplicates` CLI command that exposes the duplicate detection functionality
- Support finding duplicates within a specific remote storage location
- Support finding duplicates across all indexed remotes
- Provide filtering options for hash type, minimum file size, and result limits
- Use the existing `FormatReport()` function for human-readable output
- Return proper exit codes and error messages

## Impact
- Affected specs:
  - New capability `duplicates-cli`
- Affected code:
  - `main.go` CLI command dispatch and argument parsing
  - New command handler functions for duplicate detection
  - Integration with existing `internal/analysis` package
  - Integration with existing `internal/repository` package
- Users can now:
  - Quickly identify duplicate files to reclaim storage space
  - Understand file redundancy across their remote storage locations
  - Make informed decisions about which duplicate copies to remove

