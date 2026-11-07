## ADDED Requirements

### Requirement: Incremental Remote File Scanning

The system SHALL support scanning only new/missing directories on remote storage (via rclone), enabling efficient incremental updates via the `--only-missing-dirs` flag on the `index` command.

#### Scenario: Incremental scan skips indexed remote directories

- **WHEN** user runs `filesprawl index --remote media --only-missing-dirs`
- **THEN** the system queries database for previously indexed directories for that remote
- **AND** the system skips any directories already indexed
- **AND** only new/missing directories are listed and processed

#### Scenario: Full remote scan remains default behavior

- **WHEN** user runs `filesprawl index --remote media` (without flag)
- **THEN** the system performs a full scan as before
- **AND** all remote directories are indexed regardless of prior scans

#### Scenario: Scanner accepts scanned directory set

- **WHEN** Scanner is created with `WithScannerScannedDirs(dirs)` option
- **THEN** the scanner skips directories before descending into them
- **AND** directories in the set are completely bypassed during recursion

### Requirement: Remote Directory Filtering

The system SHALL efficiently filter remote directories by comparing rclone-discovered directories against a pre-computed set of indexed directories.

#### Scenario: Remote directory filtering during DirsOnly phase

- **WHEN** Scanner processes directories in DirsOnly mode
- **THEN** each directory is checked against the scanned set
- **AND** matching directories are skipped before recursing
- **AND** non-matching directories are descended into for file processing

#### Scenario: Location ID correctly set before scanning

- **WHEN** incremental scan for a remote begins
- **THEN** the system ensures the location exists in the database
- **AND** the Remote.LocationId is synchronized with Location.Id
- **AND** all file persistence operations use the correct location ID
