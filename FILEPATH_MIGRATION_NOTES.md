# Filepath Storage & Display Migration

## Current Database State

The database stores file paths split into:

- `remote.remote_name` = scan root (e.g., `/goblet/light`)
- `object_meta.object_path` = relative path (e.g., `photos/travis-transfer/2019/03/10/IMG_0587.MOV`)
- `object_meta.scan_root` = redundant copy of scan root

## New Design Requirements

The UI now expects:

- **Local files**: `RemoteName` = full canonical path (scan_root + relative path)
- **Rclone files**: `RemoteName` = remote identifier (e.g., `goblet`)

## Migration Strategy

### Option 1: Re-Scan (Recommended)

Run scan again for all indexed locations:

```bash
filesprawl index --remote <name> --path <path>
```

This automatically stores correct data:

- `remote.remote_name` = full canonical path for local files
- `object_meta.object_path` = relative path (unchanged)
- `object_meta.scan_root` = scan root (for type detection)

### Option 2: Database Migration (Complex)

Compute full paths from existing data:

```sql
UPDATE remote
SET remote_name = object_meta.scan_root || '/' || object_meta.object_path
WHERE location.type = 'local'
```

**Issue**: This requires joining tables and won't work if scan_root isn't set on all records.

## Recommendation

**Re-scan all sources** to ensure data integrity and consistency. The UI has been updated to work with the correct data structure, but existing scans need to be re-indexed.

## Current UI Behavior

The UI code now:

1. Detects file type using `scan_root` (local vs rclone)
2. Displays `RemoteName` as the remote/hostname
3. Constructs full paths for preview: `RemoteName/Path` (local) or `RemoteName:Path` (rclone)

This assumes `RemoteName` is already the full canonical path for local files.
