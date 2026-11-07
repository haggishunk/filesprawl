## ADDED Requirements

### Requirement: Incremental Local File Scanning

The system SHALL support scanning only new/missing directories on the local filesystem, enabling efficient incremental updates via the `--only-missing-dirs` flag on the `scan-local` command.

#### Scenario: Incremental scan skips indexed directories

- **WHEN** user runs `filesprawl scan-local /path --only-missing-dirs`
- **THEN** the system queries database for previously indexed directories
- **AND** the system skips any directories already indexed
- **AND** only new/missing directories are walked and hashed

#### Scenario: Full scan remains default behavior

- **WHEN** user runs `filesprawl scan-local /path` (without flag)
- **THEN** the system performs a full scan as before
- **AND** all directories are indexed regardless of prior scans

#### Scenario: LocalScanner accepts scanned directory set

- **WHEN** LocalScanner is created with `WithScannedDirs(dirs)` option
- **THEN** the scanner filters directories using `filepath.SkipDir`
- **AND** directories in the set are completely skipped with their subtrees

### Requirement: Directory-Based Filtering

The system SHALL efficiently determine which directories to scan by comparing discovered directories against a pre-computed set of indexed directories.

#### Scenario: Directory set lookup is O(1)

- **WHEN** filtering directories during scan
- **THEN** the system uses a map-based lookup (not linear search)
- **AND** performance scales to thousands of directories

#### Scenario: Parent directory skip includes children

- **WHEN** a parent directory is in the scanned set
- **THEN** the system skips all child directories automatically
- **AND** uses `filepath.SkipDir` to prevent subtree traversal
