# Change: Fix Symlink Handling in Local File Scanning

## Why
Local file scanning encounters three critical issues with symlinks and path deduplication:

1. **Directory symlinks crash the scan** - `filepath.WalkDir` follows directory symlinks, leading to "is a directory" errors when attempting to hash them. The scan fails entirely instead of gracefully skipping.

2. **File symlinks are silently followed** - File symlinks are read and hashed without any indication they are symlinks. The resulting file is indistinguishable from its target, creating false duplicates.

3. **Path deduplication breaks cross-scan matching** - Files with identical names in separately scanned folders (e.g., `photos/a.jpg` and `photos-alt/a.jpg`) reuse the same metadata record instead of creating distinct entries. This prevents duplicate detection across separate scan roots.

## What Changes

### Symlink Handling
- Use `filepath.WalkDir` with `fs.SkipDir` to avoid descending into directory symlinks
- Detect file symlinks using `os.Lstat()` (does not follow the link)
- Record symlink status in `object_meta` via new `is_symlink` boolean column
- Skip hashing symlinks entirely (or hash link target with explicit marking)

### Path Deduplication Fix
- Include the scan root in the metadata record (new `scan_root` column)
- Normalize the path to be relative to the scan root
- Modify duplicate queries to account for scan origin when matching files
- Same-named files from different scans now create distinct records

### Database Schema
- Add `is_symlink` BOOLEAN to `object_meta` table (default: false)
- Add `scan_root` TEXT to `object_meta` table (identifies scan origin)
- Update indices to optimize queries including scan_root

## Impact

**Affected Specs:**
- Modified: `local-file-scanning` - Update scanning behavior for symlinks
- Modified: `object-repository` - Extend metadata schema and queries
- Modified: `database-persistence` - Schema changes required

**Affected Code:**
- `internal/operation/local_scan.go` - Add symlink detection and skipping
- `internal/repository/repository.go` - Update PersistLocalResult to include scan_root
- `db/init/schema.sql` - Add is_symlink and scan_root columns
- `internal/analysis/analysis.go` - Update duplicate matching logic

**Breaking Changes:**
- ⚠️ **Database migration required** - New columns in object_meta table
- ⚠️ Existing duplicate detection may show different results post-migration

**User Benefits:**
- Scans with symlinks complete successfully
- False duplicates from symlink targets eliminated
- Correct duplicate detection across multiple scan roots
- Proper tracking of scan provenance

## Scenarios

1. User scans `/photos` then `/photos-alt` - files with same name detected as duplicates
2. Directory contains symlinks - scan completes without error, symlinks skipped
3. File symlink exists - marked as symlink, not falsely matched to target
4. Large directory tree with mixed symlinks - scans efficiently, reporting complete
