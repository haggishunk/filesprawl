# Scanning Optimization Plan: Directory-First Incremental Scanning

## Current State Analysis

**Current Flow:**

1. Walk entire directory tree (local) or recurse through all paths (rclone)
2. Hash every file found
3. Persist each file individually to database
4. No change detection - all files rescanned on each run

**Performance Issues:**

- ❌ Full tree walk even if only leaves changed
- ❌ Hash computation for unchanged files
- ❌ No way to detect deleted files
- ❌ No way to detect moved files (different path, same hash)
- ❌ Database queries for every single file

## Optimization Strategy: Three-Phase Scanning

### Phase 1: Directory Structure Snapshot (Fast)

**Goal: Detect which directories have changed**

```
GetDirectorySnapshot(location, scanRoot) -> DirectoryTree
  - For each directory: size, modtime, file count
  - Store in memory (minimal overhead)
  - Compare with previous snapshot in DB
  - Identify "changed leaves" (dirs with new/deleted files)
```

**Benefits:**

- Skip unchanged directories entirely
- Identify candidates for full scan vs delta scan

### Phase 2: Delta Scanning (Targeted)

**Goal: Scan only changed directories and their contents**

```
for each "changed leaf":
  - List current files
  - Compare against DB for this location+directory
  - Classify:
    * NEW: In filesystem, not in DB → hash & insert
    * UPDATED: In DB with different size/mtime → rehash & update
    * DELETED: In DB, not in filesystem → mark deleted
    * UNCHANGED: In DB with same size/mtime → skip hashing
```

**Benefits:**

- Hash only files that actually changed
- 50-90% reduction in hash computations
- Incremental database updates

### Phase 3: Consistency Check (Periodic)

**Goal: Catch edge cases (moved files, hash collisions)**

```
On every Nth scan:
  - Full hash check of "UNCHANGED" files
  - Detect moved files (same hash, different path)
  - Detect inode reuse edge cases
```

## Implementation: New Database Fields

Add to `object_meta` table:

```sql
ALTER TABLE object_meta ADD COLUMN IF NOT EXISTS (
  last_scanned_at TIMESTAMP,
  last_modtime BIGINT,      -- Unix timestamp
  is_deleted BOOLEAN DEFAULT false,
  deletion_detected_at TIMESTAMP
);

CREATE INDEX idx_object_meta_deleted ON object_meta(location_id, is_deleted);
CREATE INDEX idx_object_meta_scanned ON object_meta(location_id, last_scanned_at);
```

## New API: ScannerOptions

```go
type ScannerOptions struct {
  DeltaMode         bool          // Skip unchanged dirs
  HashTimeout       time.Duration // Timeout per file
  DirectorySnapshot DirectorySnapshot
  SkipHashTypes     []string      // e.g., ["dropbox"] if expensive
}

type DirectorySnapshot struct {
  DirHash string    // Tree structure hash
  FileCount int     // Count of files in tree
  ModTime time.Time // Latest modtime in tree
}
```

## Implementation Phases

**Phase A (High Priority):**

- Add last_scanned_at, last_modtime to schema
- Implement delta detection (size/mtime check)
- SKIP hash recomputation if unchanged
- Mark deleted files instead of removing

**Phase B (Medium Priority):**

- Directory snapshot comparison
- Skip hashing for unchanged directories
- Parallel scanning of independent trees
- Progress tracking/reporting

**Phase C (Low Priority):**

- Moved file detection (same hash, path changed)
- Full consistency check interval
- Compression of old scan history
- Custom hash priorities per location

## Expected Performance Gains

| Scenario              | Before           | After           | Gain       |
| --------------------- | ---------------- | --------------- | ---------- |
| Incremental 1% change | 100 files hashed | 1-3 hashed      | 95% ⬇️     |
| New files only        | All files hashed | New files only  | 80-99% ⬇️  |
| No changes            | All files hashed | 0 hashed        | 100% ⬇️    |
| Moved files           | All files hashed | All hashed once | Same       |
| Deleted files         | Add time         | ~1% overhead    | Minimal ⬆️ |

## Schema Migration

```sql
-- New migration: 0005_add_incremental_scan_fields.up.sql
ALTER TABLE object_meta ADD COLUMN last_scanned_at TIMESTAMP;
ALTER TABLE object_meta ADD COLUMN last_modtime BIGINT;
ALTER TABLE object_meta ADD COLUMN is_deleted BOOLEAN DEFAULT false;
ALTER TABLE object_meta ADD COLUMN deletion_detected_at TIMESTAMP;

CREATE INDEX idx_object_meta_deleted
  ON object_meta(location_id, is_deleted);
CREATE INDEX idx_object_meta_scanned
  ON object_meta(location_id, last_scanned_at);

-- Backfill: mark all current records as scanned now
UPDATE object_meta SET last_scanned_at = CURRENT_TIMESTAMP
WHERE last_scanned_at IS NULL;
```

## Backward Compatibility

- Non-delta scans still work (DeltaMode = false)
- Old DB works (new fields nullable)
- Gradual rollout: add fields → enable delta → optimize

## Testing Strategy

1. Unit tests: Delta detection logic
2. Integration tests: Verify deleted file detection
3. Performance tests: Compare delta vs full scan
4. Migration tests: Backward compatibility
