# Implementation Tasks: Incremental Directory-Based Scanning

## 1. LocalScanner Enhancement

- [x] 1.1 Add `OnlyMissingDirs` and `ScannedDirs` fields to LocalScanner struct
- [x] 1.2 Implement `WithOnlyMissingDirs()` option function
- [x] 1.3 Implement `WithScannedDirs()` option function
- [x] 1.4 Modify `ScanLocal()` to filter directories using `filepath.SkipDir`

## 2. Scanner (Remote) Enhancement

- [x] 2.1 Add `OnlyMissingDirs` and `ScannedDirs` fields to Scanner struct
- [x] 2.2 Implement `WithScannerOnlyMissingDirs()` option function
- [x] 2.3 Implement `WithScannerScannedDirs()` option function
- [x] 2.4 Modify `Scan()` to skip directories before recursion

## 3. CLI scan-local Command

- [x] 3.1 Add `--only-missing-dirs` flag
- [x] 3.2 Implement hostname resolution for local location
- [x] 3.3 Call `EnsureLocation()` to get location ID
- [x] 3.4 Call `GetScannedDirectories()` when flag is set
- [x] 3.5 Pass scanned dirs to LocalScanner options

## 4. CLI index Command

- [x] 4.1 Add `--only-missing-dirs` flag
- [x] 4.2 Call `EnsureLocation()` for the remote
- [x] 4.3 Sync `Remote.LocationId` after ensuring location
- [x] 4.4 Call `GetScannedDirectories()` when flag is set
- [x] 4.5 Pass scanned dirs to Scanner options

## 5. Testing & Validation

- [x] 5.1 Build without errors
- [x] 5.2 All repository integration tests pass
- [x] 5.3 All analysis tests pass
- [x] 5.4 Verify help text shows new flags
- [x] 5.5 Test basic command execution

## 6. Documentation

- [x] 6.1 Create INCREMENTAL_SCANNING_DESIGN.md
- [x] 6.2 Create INCREMENTAL_SCANNING_CLI.md

## 7. Git & PR

- [ ] 7.1 Create feature branch for proposal
- [ ] 7.2 Commit proposal, tasks, design.md, and spec deltas
- [ ] 7.3 Push and open PR for review
- [ ] 7.4 Get approval before merging to main
