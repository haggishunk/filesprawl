## MODIFIED Requirements

### Requirement: Local File Hashing
The system SHALL compute supported content hashes for local files discovered during local scanning, excluding symlinks.

#### Scenario: Skip directory symlinks
- **WHEN** local scanning encounters a symbolic link to a directory
- **THEN** the symlink is skipped and scanning continues with remaining files

#### Scenario: Mark file symlinks
- **WHEN** a file symlink is encountered during local scanning
- **THEN** it is marked as `is_symlink=true` in metadata and not hashed

#### Scenario: Hash readable regular files
- **WHEN** a readable local regular file is encountered
- **THEN** the system computes supported content hashes from file content

#### Scenario: Unreadable file during scan
- **WHEN** a local file cannot be opened or read
- **THEN** scan returns or records an error with the local file path in context

### Requirement: Scan Origin Tracking
The system SHALL track the scan root directory for each indexed local file.

#### Scenario: Record scan root
- **WHEN** indexing a local path `/home/user/photos`
- **THEN** each file is persisted with `scan_root=/home/user/photos`

#### Scenario: Distinguish same-named files
- **WHEN** scanning `/photos/a.jpg` and then `/photos-alt/a.jpg`
- **THEN** both files are stored as distinct records (not reused)

#### Scenario: Duplicate detection across scans
- **WHEN** finding duplicates across multiple scanned roots
- **THEN** same-named files from different scan roots are correctly identified as duplicates if hashes match
