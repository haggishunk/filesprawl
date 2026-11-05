# Change: Optimize Incremental Scanning with Delta Detection

## Why

Current scanning implementation rescans **all** files and recomputes hashes for **all** files on every scan operation. For a 1M photo library:

- Scan time: 27+ hours (prohibitive for weekly/daily scans)
- CPU waste: 99% of files unchanged, but all get hashed
- No change visibility: Can't track which files were deleted
- Storage growth: Deleted files remain in DB forever

The goal is to make scanning practical for large libraries by enabling incremental updates that skip unchanged files.

## What Changes

**BREAKING CHANGES:**

- None. Feature is backward compatible via optional `--delta` flag

**New Capabilities:**

- Delta scanning mode: Skip hashing for unchanged files (size+mtime match)
- Deletion tracking: Mark deleted files in database with timestamp
- Change classification: Files categorized as UNCHANGED/NEW/CHANGED/DELETED
- Incremental performance: 95% faster for typical incremental scans (1% changes)

**Affected Specs:**

- `file-scanning`: Add delta detection algorithm and change classification
- `database-persistence`: Add fields for change tracking (last_scanned_at, last_modtime, is_deleted, deletion_detected_at)

## Impact

### Behavioral Changes

- New `--delta` flag enables incremental scanning (default: disabled for compatibility)
- Files with unchanged size+mtime skip hash computation
- Deleted files marked in DB instead of purged
- Progress reporting shows what was skipped

### Performance Gains

| Scenario   | Before | After | Reduction |
| ---------- | ------ | ----- | --------- |
| 1% change  | 27h    | 18min | 95%       |
| 10% change | 27h    | 3h    | 90%       |
| No changes | 27h    | 2min  | 100%      |

### Database Schema

- Add: last_scanned_at, last_modtime, is_deleted, deletion_detected_at
- Add indexes: idx_object_meta_deleted, idx_object_meta_scanned

### Code Changes

- New: `internal/operation/delta_scan.go` - DeltaScan type
- New: `internal/operation/delta_detection.go` - Classification logic
- Modify: `internal/repository/repository.go` - New query methods
- Modify: CLI scan commands - Add --delta flag
- Modify: `internal/operation/local_scan.go`, `scan.go` - Integrate delta logic

### Testing

- Unit tests for delta detection (unchanged/new/changed/deleted classification)
- Integration tests for --delta flag with real filesystems
- Migration tests for schema additions
- Backward compatibility tests

## Implementation Phases

- **Phase 1 (Critical):** Schema migration + repository methods (1-2 weeks)
- **Phase 2 (Critical):** Delta scanning engine (1-2 weeks)
- **Phase 3 (Important):** CLI integration (1 week)
- **Phase 4-5 (Optional):** Optimizations and advanced features (3+ weeks)

## Risk Mitigation

- **Timestamp issues:** Use UTC, Unix timestamps, fallback to full rescan
- **Mtime unavailable:** Rclone compatibility layer with graceful degradation
- **Edge cases:** Periodic full verification every 10th scan
- **Rollback:** --delta flag disabled by default, easy to disable

See `design.md` for technical decisions and `tasks.md` for implementation checklist.
