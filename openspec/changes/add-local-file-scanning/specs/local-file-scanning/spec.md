## ADDED Requirements

### Requirement: Local Scan Command
The system SHALL provide a CLI command for indexing a local filesystem path.

#### Scenario: Index a local root
- **WHEN** a user runs a local indexing command with a valid filesystem root
- **THEN** the system starts a local filesystem scan for that root

#### Scenario: Missing local path
- **WHEN** a user invokes the local indexing command without a required path
- **THEN** the command fails with a validation error explaining that a local path is required

### Requirement: Local File Hashing
The system SHALL compute supported content hashes for local files discovered during local scanning.

#### Scenario: Hash a readable file
- **WHEN** a readable local file is encountered during local scanning
- **THEN** the system computes supported content hashes from the file content

#### Scenario: Unreadable file during scan
- **WHEN** a local file cannot be opened or read
- **THEN** the scan returns or records an error with the local file path in context

### Requirement: Local Files in Duplicate Workflows
The system SHALL make locally scanned files available to the same duplicate lookup workflows used for remote-backed files.

#### Scenario: Local and remote files share a hash
- **WHEN** a local file and a remote-backed file are indexed with the same hash
- **THEN** duplicate lookup can return both files in the same duplicate group