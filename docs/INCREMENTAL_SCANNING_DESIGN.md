# Incremental Directory-Based Scanning Design

## Overview

The `--only-missing-dirs` feature enables efficient incremental scanning by only indexing directories that haven't been previously scanned. This significantly speeds up re-scans of large storage locations by skipping already-indexed directories.

## Schema

### object_directory Column

```sql
ALTER TABLE object_meta ADD COLUMN object_directory TEXT NOT NULL;
```

**Purpose**: Stores the immediate parent directory of each file.

**Format**:

- **Local files**: Absolute filesystem directory path
  - Example: `/goblet/light/photos/2019`
- **Remote files**: Directory path within remote
  - Example: `xforants/light/wedding photos professional`
- **Root files**: Empty string (though rare in practice)

**Index**:

```sql
CREATE INDEX idx_object_meta_directory
  ON object_meta(location_id, object_directory);
```

## Implementation Details

### Derivation Logic

The `object_directory` is derived from `object_path` using `filepath.Dir()`:

```
object_path = /path/to/file.txt
object_directory = /path/to

object_path = photos/subdir/image.jpg
object_directory = photos/subdir
```

This is automatically computed in `object.Meta.BasePath` during `NewMeta()`.

### Database Operations

**Query indexed directories**:

```go
GetScannedDirectories(ctx, locationID) -> map[string]bool
```

Returns all unique directories that have been indexed for a location.

**Insertion**:
WriteMeta automatically persists `object_directory` from `Meta.BasePath`.

## Usage: --only-missing-dirs Flag

### Local Scans

```bash
# Discover directories on local machine
find /path/to/scan -type d > /tmp/local_dirs.txt

# Query database for indexed directories
filesprawl scan-local /path/to/scan --only-missing-dirs

# CLI logic:
# 1. Get indexed directories from DB for this location
# 2. Query filesystem for all directories
# 3. Find diff = directories in filesystem but not in DB
# 4. Only walk those directories
```

### Remote Scans

```bash
# Discover directories on remote via rclone
rclone rc operations/list fs=remote: recursive=true dirs-only=true

# Query database and scan only missing directories
filesprawl index --remote media --only-missing-dirs

# CLI logic:
# 1. Get indexed directories from DB for this location
# 2. Call rclone lsd with recursive=true
# 3. Find diff = directories in remote but not in DB
# 4. Only scan those directories with targeted path filters
```

## Performance Benefits

**Before**:

- 133,353 files across 447 directories
- Every re-scan walks entire tree
- Average scan time: minutes (depending on storage I/O)

**After**:

- First scan: unchanged (full tree walk)
- Incremental scan: only walk new directories
- Example: 5 new directories = ~0.5 minute scan vs 5+ minute full scan

## Data Statistics

Current database state:

- **Total files**: 133,353
- **Unique directories**: 447
- **Largest directory**: 6,671 files (xforants/light/cara-gphoto)
- **Average files/directory**: ~299

## Future Enhancements

1. **Directory modification tracking**: Track scan timestamps per directory
2. **Change detection**: Skip unchanged directories via mtime/etag checks
3. **Parallel directory scanning**: Scan independent directories concurrently
4. **Sync mode**: Detect deleted files/directories and clean up database
