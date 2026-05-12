## ADDED Requirements

### Requirement: Local Scan Result Persistence
The system SHALL persist local scan results into the same object metadata and hash index used for remote-backed scans.

#### Scenario: Persist local file metadata and hash
- **WHEN** a local file is scanned and hashed
- **THEN** its metadata, hashes, and metadata-to-hash associations are persisted through the repository layer

#### Scenario: Mixed local and remote hash lookup
- **WHEN** duplicate lookup is performed for a hash shared by local and remote-backed files
- **THEN** repository file lookup can return both local and remote-backed file records