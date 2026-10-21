# Implementation Tasks

## 1. Database Schema Migration
- [x] 1.1 Add `is_symlink BOOLEAN DEFAULT false` column to object_meta table
- [x] 1.2 Add `scan_root TEXT` column to object_meta table
- [x] 1.3 Create migration script with rollback capability
- [x] 1.4 Update schema.sql for new installations
- [x] 1.5 Test migration against existing test database

## 2. Local File Scanning - Symlink Handling
- [x] 2.1 Update local_scan.go to use os.Lstat() for symlink detection
- [x] 2.2 Skip directory symlinks (WalkDir never descends into them; skip without fs.SkipDir since the entry is not a directory)
- [x] 2.3 Detect and skip file symlinks (set is_symlink=true, don't hash)
- [x] 2.4 Add error handling for broken symlinks with clear messages
- [x] 2.5 Add unit tests for symlink handling

## 3. Local File Scanning - Scan Root Tracking
- [x] 3.1 Update ScanLocal signature to accept and track scan root
- [x] 3.2 Normalize scan root path (absolute, resolved)
- [x] 3.3 Pass scan_root through to PersistLocalResult
- [x] 3.4 Update local_scan_test.go with scan_root tests

## 4. Repository - Metadata Persistence
- [x] 4.1 Update object.Meta struct to include IsSymlink and ScanRoot
- [x] 4.2 Update PersistLocalResult to persist is_symlink and scan_root
- [x] 4.3 Update WriteMeta query to handle new columns
- [x] 4.4 Update ReadMeta query to retrieve new columns

## 5. Duplicate Detection - Path Deduplication
- [x] 5.1 Update FindDuplicatesWithinRemote to account for scan_root (distinct records per scan root; verified by integration tests)
- [x] 5.2 Update FindDuplicatesAcrossRemotes to account for scan_root (distinct records per scan root; verified by integration tests)
- [x] 5.3 Update GetFilesByHash to include scan_root in results
- [x] 5.4 Add integration tests for cross-scan duplicate detection

## 6. Testing
- [x] 6.1 Create test setup with symlinked directories
- [x] 6.2 Test directory symlink skipping (no crash)
- [x] 6.3 Test file symlink detection and marking
- [x] 6.4 Test same-named files from different scans
- [x] 6.5 Test duplicates command with symlink-heavy structure
- [x] 6.6 Run go test ./... (all tests passing)

## 7. Validation & Documentation
- [x] 7.1 Run openspec validate fix-symlink-handling --strict
- [x] 7.2 Update README.md with symlink handling note
- [x] 7.3 Document migration steps
- [x] 7.4 Test against issue #3 reproduction case
