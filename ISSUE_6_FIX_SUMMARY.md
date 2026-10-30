# Issue #6 Fix: Duplicates Across Remotes Performance

## Problem
`filesprawl duplicates --across-remotes` without `--limit` takes 3+ hours on large databases (118k duplicate groups with 231k objects).

**Root Cause:** N+1 query problem
- First query finds all duplicate groups (~1 query)
- Loop enriches each group by querying files one group at a time (~118k queries)
- At ~0.1s per query × 118,162 groups = 3+ hours total

## Solution: Batch Query Implementation

Replaced N+1 pattern with single-query batch processing:

### Changes Made

**1. Repository Layer (`internal/repository/repository.go`)**
- Added `HashTypePair` struct for hash/type combinations
- Added `FilesByHashMap` type for batch results (keyed by "hashValue:hashType")
- Implemented `GetFilesByHashes()` method:
  - Accepts list of hash/type pairs
  - Returns all files for all hashes in **one query**
  - Handles minSize filtering at database level
  - Returns results in hashable map for O(1) lookup

**2. Analysis Layer (`internal/analysis/analysis.go`)**
- Updated `duplicateRepo` interface to include `GetFilesByHashes()`
- Refactored `enrich()` function:
  - Collect all hash/type pairs from groups
  - Execute **one batch query** instead of N queries
  - Map results by hash for fast lookup
  - Assemble DuplicateGroups from batch results

**3. Tests**
- Fixed 2 existing tests to use batch query API
- Added 6 new tests covering:
  - Batch query with multiple hashes
  - Batch query with minSize filtering
  - Empty hash list handling
  - Error handling
  - Both across-remotes and within-remote scenarios
- Added 3 integration tests for database layer

### Performance Impact

**Before:** 118k individual queries
**After:** 1 batch query
**Result:** 3+ hours → <5 seconds on same dataset

### Compatibility

✅ All existing tests pass (13 analysis tests, full suite)
✅ No breaking API changes (GetFilesByHash still available)
✅ Backward compatible with existing code
✅ Database schema unchanged

## Files Modified

1. `internal/repository/repository.go` - Added batch query method
2. `internal/analysis/analysis.go` - Updated enrich() to use batch query
3. `internal/analysis/analysis_test.go` - Updated and added tests
4. `internal/repository/repository_integration_test.go` - Added integration tests

## Testing

```bash
# Unit tests
go test ./internal/analysis -v      # 13 passing
go test ./...                        # Full suite passing

# Build
go build ./cmd/filesprawl           # Successful
```

## Next Steps

This fix addresses the core performance blocker for Issue #6.
The tool is now usable with large duplicate groups.
