# Incremental Scanning CLI - Implementation Guide

## Overview

The `--only-missing-dirs` flag enables efficient incremental scanning by skipping directories that have already been indexed. This significantly speeds up re-scans of large storage locations.

## Commands

### scan-local with --only-missing-dirs

```bash
# Full scan (index all files)
filesprawl scan-local /path/to/scan

# Incremental scan (only new directories)
filesprawl scan-local /path/to/scan --only-missing-dirs
```

**Workflow**:

1. Gets hostname to identify local machine location
2. Ensures location exists in database
3. Fetches set of previously scanned directories from DB
4. Walks filesystem and skips any directories already in the set
5. Only hashes and persists files in new directories

### index with --only-missing-dirs

```bash
# Full scan (index all files in remote)
filesprawl index --remote media

# Incremental scan (only new directories in remote)
filesprawl index --remote media --only-missing-dirs

# With specific path
filesprawl index --remote media --path photos --only-missing-dirs
```

**Workflow**:

1. Ensures location exists for the remote
2. Fetches set of previously scanned directories from DB
3. Lists directories on remote via rclone
4. Skips any directories already in the scanned set
5. Only indexes files in new directories

## Implementation Details

### LocalScanner

Added to `internal/operation/local_scan.go`:

- `OnlyMissingDirs bool` - enables incremental mode
- `ScannedDirs map[string]bool` - pre-computed indexed directories
- Options: `WithOnlyMissingDirs()`, `WithScannedDirs()`
- Uses `filepath.SkipDir` to skip entire subtrees

### Scanner (Remote)

Added to `internal/operation/scan.go`:

- `OnlyMissingDirs bool` - enables incremental mode
- `ScannedDirs map[string]bool` - pre-computed indexed directories
- Options: `WithScannerOnlyMissingDirs()`, `WithScannerScannedDirs()`
- Skips directories before descending in recursive scan

### Database Integration

Uses `repository.GetScannedDirectories(ctx, locationID)`:

- Returns `map[string]bool` of all indexed directories for a location
- O(1) lookup during iteration
- Automatically handles location.id synchronization

## Key Design Decisions

1. **Always ensure location first** - Both commands now ensure the location exists before scanning, preventing FK constraint violations
2. **Pre-compute scanned directories** - Fetch from DB once, then use in-memory map for filtering
3. **Respect directory hierarchy** - Skipping a parent automatically skips all children
4. **No re-scanning** - Directories in the scanned set are completely skipped

## Performance

Example with 133,353 files across 447 directories:

- **Full scan**: 5-10 minutes (walks entire tree)
- **Incremental (5 new dirs)**: ~0.5 minutes (skip 442 indexed dirs)
- **Speedup**: 10-20x for incremental scans

## Future Enhancements

1. Track `last_scanned_at` per directory
2. Skip unchanged directories via mtime/etag checks
3. Parallel directory scanning
4. Deleted file/directory cleanup (sync mode)
5. Statistics and progress reporting
