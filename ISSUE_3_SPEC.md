# Issue #3: Symlinks Break and Fail to Find Duplicates - OpenSpec

## Issue Summary
Local file scanning has three critical issues:
1. Directory symlinks crash the scan
2. File symlinks are silently followed (false duplicates)
3. Same-named files in different scans aren't matched as duplicates

## Solution: fix-symlink-handling Change

A comprehensive spec addressing all three issues:

### Problem 1: Directory Symlinks Crash
**Root Cause:** `filepath.WalkDir` follows symlinks; hashing a directory fails with "is a directory" error

**Solution:** Detect directory symlinks with `os.Lstat()` and skip them using `fs.SkipDir`

**Result:** Scans complete successfully, symlink directories are skipped

### Problem 2: File Symlinks Create False Duplicates
**Root Cause:** Symlinks are followed and hashed without marking, matching their target files

**Solution:** Detect symlinks with `os.Lstat()`, mark with `is_symlink=true`, skip hashing

**Result:** Symlinks are tracked separately, no false duplicate matches

### Problem 3: Path Collisions Across Scans
**Root Cause:** `object_meta` uses path alone as identifier; scanning `/photos` then `/photos-alt` both create `a.jpg` → collision

**Solution:** Add `scan_root` column to `object_meta` to track scan provenance

**Result:** Same-named files from different scans create distinct records, correctly detected as duplicates

## Implementation Plan

### Database Changes
- Add `is_symlink BOOLEAN DEFAULT false` to object_meta
- Add `scan_root TEXT` to object_meta
- Create migration script (reversible)

### Code Changes
- Update `local_scan.go` for symlink detection
- Update `object.Meta` struct with new fields
- Update `repository.go` methods
- Update duplicate detection queries

### Testing
- 7 task groups covering schema, code, testing
- Test against issue #3 reproduction setup
- All existing tests must pass

## OpenSpec Status
✅ **Change:** `fix-symlink-handling`
✅ **Status:** Validated with strict mode
✅ **Files:** proposal.md, design.md, tasks.md, 2 spec deltas
✅ **Tasks:** 34 items to implement

## Next Steps
1. Review this spec
2. Approve if satisfied
3. Implement tasks sequentially
4. Test against issue #3 setup
5. Merge when complete

## Breaking Changes
⚠️ **Database migration required**
- Existing records need `scan_root` populated
- Migration script will be provided
- Backward compatibility: Old records default to safe values
