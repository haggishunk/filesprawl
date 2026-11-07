## ADDED Requirements

### Requirement: Delta Scanning Mode

The system SHALL support a delta scanning mode that tracks file changes between scans instead of re-scanning all files.

#### Scenario: Enable delta scanning with flag

- **WHEN** user invokes scan command with `--delta` flag
- **THEN** the system activates incremental scanning mode

#### Scenario: Skip unchanged files

- **WHEN** delta scanning is active and a file has not changed since last scan
- **THEN** the file is skipped and not re-hashed

#### Scenario: Detect new files

- **WHEN** delta scanning encounters a file not present in previous scan database
- **THEN** the file is hashed and indexed as a new object

### Requirement: File Classification

The system SHALL classify files as new, modified, or unchanged based on delta scanning logic.

#### Scenario: New file classification

- **WHEN** a file is encountered that has no prior database entry
- **THEN** it is classified as new and marked for hashing

#### Scenario: Modified file classification

- **WHEN** a file's mtime differs from its last scan database record
- **THEN** it is classified as modified and re-hashed

#### Scenario: Unchanged file classification

- **WHEN** a file's mtime matches its last scan database record
- **THEN** it is classified as unchanged and skipped

### Requirement: Mtime Change Detection

The system SHALL use file modification time to detect changes during delta scanning.

#### Scenario: Mtime comparison

- **WHEN** delta scanning processes a file
- **THEN** the current file mtime is compared against the stored mtime from the last scan

#### Scenario: Mtime update on modification

- **WHEN** a file is re-scanned and found modified
- **THEN** its stored mtime is updated to the current value in the database

### Requirement: Delta Scanning Statistics

The system SHALL collect and report statistics on delta scanning performance.

#### Scenario: Count file classifications

- **WHEN** delta scanning completes
- **THEN** statistics include count of new, modified, and unchanged files

#### Scenario: Report scanning efficiency

- **WHEN** delta scanning statistics are generated
- **THEN** efficiency metrics show percentage of files skipped vs. re-hashed

## MODIFIED Requirements

### Requirement: File Scanning

The system SHALL support both full and incremental scanning modes via command-line options.

#### Scenario: Full scan (default behavior)

- **WHEN** scanning without `--delta` flag
- **THEN** all files are scanned and hashed regardless of prior scan history

#### Scenario: Delta scan with --delta flag

- **WHEN** scanning with `--delta` flag
- **THEN** only new and modified files are hashed, unchanged files are skipped

#### Scenario: Mixed scan results

- **WHEN** delta scanning completes
- **THEN** final database contains combined results from unchanged files (from prior scans) and newly scanned files

### Requirement: Progress Reporting

The system SHALL report progress during file scanning with delta mode information.

#### Scenario: Report files processed in delta mode

- **WHEN** delta scanning is in progress
- **THEN** progress output indicates files skipped, new files found, and modified files being hashed

#### Scenario: Show scanning mode in output

- **WHEN** scan completes
- **THEN** output explicitly indicates whether full or delta scanning was performed
