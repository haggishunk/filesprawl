# Scanning Optimization - Detailed Technical Specification

## Problem Statement

Current scanning approach:

- **Walks entire filesystem tree** every scan
- **Computes hashes for all files**, including unchanged ones
- **No change detection** between scans
- **No deleted file tracking** - deleted files remain in DB forever
- **No moved file detection** - treats moves as delete + insert

For a 1M file photo library:

- Scanning takes **hours** (dominated by hash computation)
- 99% of files unchanged between scans
- Storage grows unbounded with deleted files
- Can't distinguish "file deleted" from "file not yet indexed"

## Solution Architecture

### 1. Metadata Enhancements

**Add scan tracking fields to `object_meta`:**

```
last_scanned_at: TIMESTAMP
  - When this file was last verified
  - Enables incremental scanning

last_modtime: BIGINT (Unix timestamp)
  - From filesystem/rclone (mtime)
  - Changes = file likely changed

is_deleted: BOOLEAN
  - TRUE = file gone, but record kept
  - Enables tracking of deletion patterns

deletion_detected_at: TIMESTAMP
  - When we first noticed deletion
  - Track lifecycle of files
```

### 2. Delta Scanning Algorithm

**Phase 1: Collect Current State**

```
for each file in filesystem:
  - Get: name, size, mtime (NOT hash)
  - O(n) file system operations
  - Very fast (milliseconds per file)
```

**Phase 2: Classify Files**

```
for each current file:
  if file in DB with same size+mtime:
    → UNCHANGED: skip hash, update last_scanned_at
  else:
    → CHANGED/NEW: compute hash, update/insert
```

**Phase 3: Detect Deletions**

```
for each file in DB (for this location):
  if file NOT in filesystem:
    → DELETED: set is_deleted=TRUE, deletion_detected_at=NOW
```

### 3. Implementation: DeltaScan Type

```go
// DeltaScan performs incremental scanning
type DeltaScan struct {
  Location location.Location
  ScanRoot string

  // Options
  Force bool  // Force full rescan ignoring timestamps

  // Tracking
  StatsUnchanged int
  StatsChanged   int
  StatsNew       int
  StatsDeleted   int
  StatsDuration  time.Duration

  // Handlers
  OnFile     func(file ScanFile) error
  OnDeletion func(path string) error
}

type ScanFile struct {
  AbsPath  string
  Mtime    int64
  Size     int64
  Hash     string // Empty if unchanged
  IsNew    bool
  Changed  bool
  Unchanged bool
}

// Scan returns early if file unchanged
func (ds *DeltaScan) Scan(ctx context.Context) error {
  existingFiles := ds.repo.GetFilesForLocation(ctx, ds.Location.Id)

  return filepath.WalkDir(ds.ScanRoot, func(path string, d os.DirEntry, err error) error {
    if d.IsDir() { return nil }

    // Get current file metadata
    info, _ := d.Info()
    absPath := filepath.Join(ds.ScanRoot, relPath)

    // Check if unchanged
    existing := existingFiles[absPath]
    if existing != nil &&
       existing.Size == info.Size() &&
       existing.ModTime.Unix() == info.ModTime().Unix() {
      // Skip hashing
      ds.StatsUnchanged++
      ds.repo.UpdateLastScanned(ctx, existing.Id)
      return nil
    }

    // File changed or new: compute hash
    hash := computeHash(path)
    ds.OnFile(ScanFile{
      AbsPath: absPath,
      Size: info.Size(),
      Mtime: info.ModTime().Unix(),
      Hash: hash,
      IsNew: existing == nil,
    })

    if existing == nil {
      ds.StatsNew++
    } else {
      ds.StatsChanged++
    }
    return nil
  })
}
```

### 4. Deletion Tracking

```go
// After scanning all current files
func (ds *DeltaScan) DetectDeletions(ctx context.Context) error {
  for _, file := range ds.repo.GetFilesForLocation(ctx, ds.Location.Id) {
    if !file.IsDeleted && !currentFiles[file.AbsolutePath] {
      ds.repo.MarkDeleted(ctx, file.Id)
      ds.OnDeletion(file.AbsolutePath)
      ds.StatsDeleted++
    }
  }
}
```

### 5. Repository Methods

```go
// Get all files for a location (with mtime)
GetFilesForLocation(ctx context.Context, locID int) ([]FileRecord, error)

// Update last_scanned_at timestamp only
UpdateLastScanned(ctx context.Context, fileID int) error

// Mark file as deleted (set is_deleted=TRUE, deletion_detected_at=NOW)
MarkDeleted(ctx context.Context, fileID int) error

// Update file with new hash (for changed files)
UpdateFileHash(ctx context.Context, fileID int, hash string) error
```

### 6. Benefits

✅ **95% reduction in hash computations** for typical incremental scans  
✅ **Deleted file tracking** - know what was removed  
✅ **Moved file detection** - same hash, different path  
✅ **Archive lifecycle tracking** - when files appeared/disappeared  
✅ **Bandwidth savings** - no reclone of unchanged remote files  
✅ **Backward compatible** - old DB just has NULL values

## Example: 1M Photo Library

**Scenario: Weekly incremental scan with 10K new photos**

**Before Optimization:**

```
1,000,000 files × 100ms per hash = 100,000 seconds = 27+ hours
Result: Prohibitive for weekly scans
```

**After Optimization:**

```
Phase 1: Scan directory structure = 5 seconds
Phase 2: Check size/mtime = 30 seconds (1M files, fast ops)
Phase 3: Hash only new files = 10,000 × 100ms = 1,000 seconds
Phase 4: Mark deletions = 10 seconds
Total: ~1,055 seconds = 18 minutes
Result: **95% faster** ✓
```
