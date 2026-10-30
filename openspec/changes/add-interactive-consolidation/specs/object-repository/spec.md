## ADDED Requirements

### Requirement: Consolidation History Tracking
The system SHALL record all consolidation operations with source/destination mapping for recovery and audit.

#### Scenario: Consolidation recorded
- **WHEN** file at `s3://bucket/media.jpg` is consolidated to archive
- **THEN** system records: source path, destination path, timestamp, user decision, remote info
- **AND** query system can retrieve all consolidations for a given file or group

### Requirement: Consolidation Preference Storage
The system SHALL persist and retrieve user consolidation preferences for pattern-based reuse.

#### Scenario: Preference lookup
- **WHEN** user consolidates duplicate group with pattern (remote type, file type, etc.)
- **THEN** system stores preference with match criteria
- **AND** future similar groups can be auto-consolidated using stored preference
- **AND** user can list, modify, or delete stored preferences

### Requirement: Archive Path Mapping
The system SHALL track mapping of original file locations to archive locations for recovery.

#### Scenario: Archive lookup
- **WHEN** user queries for archived files from remote `media`
- **THEN** system returns list: original path → archive path with retrieval instructions
