# Change: Add Incremental Directory-Based Scanning

## Why

Large storage locations contain thousands of files across hundreds of directories. Re-scanning the entire location to find new files is inefficient—most directories haven't changed and have already been indexed. This proposal enables scanning only new/missing directories, reducing re-scan time by 10-20x for typical incremental updates.

## What Changes

- **New CLI flag** `--only-missing-dirs` for `scan-local` and `index` commands
- **Database tracking** of scanned directories via `object_directory` column
- **Incremental filtering** at LocalScanner and Scanner levels
- **Skip unchanged** directories by comparing with pre-indexed set

Includes:

- LocalScanner supports filtering by `OnlyMissingDirs` and `ScannedDirs`
- Scanner (remote) supports filtering by `OnlyMissingDirs` and `ScannedDirs`
- `GetScannedDirectories()` repository method to query indexed directories
- CLI integration to fetch scanned dirs before scanning

## Impact

- **Affected specs**: `file-scanning`, `rclone-integration`
- **Affected code**:
  - `internal/operation/local_scan.go` - LocalScanner struct and options
  - `internal/operation/scan.go` - Scanner struct and options
  - `cmd/filesprawl/scan_local.go` - CLI flag and workflow
  - `cmd/filesprawl/index.go` - CLI flag and workflow
- **Performance**: 10-20x faster re-scans for incremental updates
- **Database**: Uses existing `object_directory` column (already added in prior work)
- **Backward compatible**: Full scan is still default; `--only-missing-dirs` is opt-in

## Deployment Notes

- No schema changes required (uses existing `object_directory` column)
- No migrations needed
- Requires `GetScannedDirectories()` repository method (already implemented)
- CLI changes are backward compatible
