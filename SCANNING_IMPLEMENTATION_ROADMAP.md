# Scanning Optimization - Implementation Roadmap

## Phase 1: Foundation (1-2 weeks) - CRITICAL PATH

### 1.1: Schema Migration

**File:** `db/migrations/0005_add_incremental_scan_fields.up.sql`

```sql
ALTER TABLE object_meta ADD COLUMN last_scanned_at TIMESTAMP;
ALTER TABLE object_meta ADD COLUMN last_modtime BIGINT;
ALTER TABLE object_meta ADD COLUMN is_deleted BOOLEAN DEFAULT false;
ALTER TABLE object_meta ADD COLUMN deletion_detected_at TIMESTAMP;

CREATE INDEX idx_object_meta_deleted ON object_meta(location_id, is_deleted);
CREATE INDEX idx_object_meta_scanned ON object_meta(location_id, last_scanned_at);
```

**Rollback:** `0005_add_incremental_scan_fields.down.sql`

### 1.2: Repository Methods

**File:** `internal/repository/repository.go`

Add:

- `GetFilesForLocation(ctx, locID) ([]FileRecord, error)` - with mtime
- `UpdateLastScanned(ctx, fileID) error` - touch timestamp only
- `MarkDeleted(ctx, fileID) error` - set is_deleted + timestamp
- `UpdateFileMtime(ctx, fileID, mtime) error` - update modtime

### 1.3: Unit Tests

**File:** `internal/repository/repository_integration_test.go`

- Test MarkDeleted sets is_deleted + timestamp
- Test UpdateLastScanned updates timestamp
- Test GetFilesForLocation returns mtime
- Test deleted files are excluded from most queries

## Phase 2: Delta Scanning Core (1-2 weeks)

### 2.1: DeltaScan Type

**File:** `internal/operation/delta_scan.go` (NEW)

```go
type DeltaScan struct {
  Location location.Location
  ScanRoot string
  Force    bool // Skip incremental, full rescan

  // Stats
  StatsUnchanged, StatsNew, StatsChanged, StatsDeleted int
}

func (ds *DeltaScan) ScanLocal(ctx context.Context) error
func (ds *DeltaScan) ScanRclone(ctx context.Context) error
func (ds *DeltaScan) DetectDeletions(ctx context.Context) error
```

### 2.2: Delta Detection Logic

**File:** `internal/operation/delta_detection.go` (NEW)

```go
type FileState struct {
  AbsPath string
  Size    int64
  Mtime   int64
  Hash    string
}

// ClassifyFile: UNCHANGED, NEW, CHANGED
func ClassifyFile(current FileState, existing *FileRecord) FileClass

// Return early if unchanged
func ShouldSkipHash(current, existing FileState) bool {
  return current.Size == existing.Size &&
         current.Mtime == existing.ModTime
}
```

### 2.3: Tests

**File:** `internal/operation/delta_scan_test.go`

- Test unchanged files skip hashing
- Test new files get hashed
- Test changed files get rehashed
- Test deletion detection
- Test stats tracking

## Phase 3: Integration (1 week)

### 3.1: CLI Updates

**File:** `internal/cli/commands.go`

Add `--delta` flag to scan commands:

```
filesprawl scan-local --delta /path
filesprawl scan-rclone --delta myremote:path
```

Default: false (backward compatible)

### 3.2: Progress Reporting

Print stats:

```
Scanning /path (delta mode)...
✓ 10,000 unchanged (skipped hash)
✓ 500 new (hashed)
✓ 50 changed (rehashed)
✓ 20 deleted (marked)
Duration: 2m 34s
```

### 3.3: Tests

- Integration test: scan-local with --delta flag
- Integration test: reclone with --delta flag
- Verify backward compat: --delta not specified works

## Phase 4: Optimization (2 weeks) - OPTIONAL

### 4.1: Directory Snapshot

Skip entire directory trees if unchanged:

- `GetDirectorySnapshot(path) -> DirHash`
- Compare snapshots between scans
- Skip descending into unchanged dirs

### 4.2: Parallel Scanning

Scan independent directory branches in parallel:

- For each top-level directory
- Spawn goroutine to scan subtree
- Merge results

### 4.3: Hash Caching

Cache hashes from previous scans:

- If size+mtime unchanged from last successful hash
- Trust previous hash (optional setting)

## Phase 5: Polish (1 week) - OPTIONAL

### 5.1: Moved File Detection

```
SELECT h.id FROM object_meta m
JOIN object_hash_junction j ON m.id = j.object_meta_id
WHERE h.hash_id = ?
  AND m.location_id = ?
  AND m.absolute_path != ?
```

Track files that moved but have same hash.

### 5.2: Consistency Check Interval

Every 10th scan: Force full hash verification of "unchanged" files.

### 5.3: CLI Reporting

```
filesprawl scan-local /path --report
Shows:
- Files unchanged since last scan
- Average time between changes
- Deleted file patterns
- Moved files detected
```

## Success Metrics

| Metric                | Before    | Target       | Status           |
| --------------------- | --------- | ------------ | ---------------- |
| 1% change scan time   | 27 hours  | <30 min      | ✅ 95% reduction |
| Hash computations     | All files | 1% of files  | ✅ 99% reduction |
| Deleted file tracking | None      | Full history | ✅ Complete      |
| DB queries            | Per file  | Per location | ✅ 1000x fewer   |
| Memory usage          | Low       | Low          | ✅ Same          |

## Risk Mitigation

- **Risk:** Timestamp issues across timezones
  - **Mitigation:** Use UTC everywhere, store as Unix timestamp

- **Risk:** File deduplication by mtime fails for fast updates
  - **Mitigation:** Add option to force full rescan, default to safety

- **Risk:** Rclone doesn't always provide mtime
  - **Mitigation:** Fall back to hash for Rclone, full scan for safety

## Rollout Plan

1. Deploy Phase 1 (schema + repo methods)
2. Deploy Phase 2 (delta scanning) with --delta flag OFF by default
3. Manual testing with --delta on test data
4. Gradual rollout: enable for users on opt-in basis
5. Monitor: track delta stats, watch for edge cases
6. Default to on after 2 weeks of clean production data
7. Deploy Phase 3+ (optimizations) after foundation is stable
