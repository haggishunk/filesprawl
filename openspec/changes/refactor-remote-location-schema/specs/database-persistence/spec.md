## ADDED Requirements

### Requirement: Location Table for Source Identification
The system SHALL maintain a location table identifying the source of indexed data (local hostname or remote name).

#### Scenario: Local scan location
- **WHEN** user scans local filesystem from host `archlinux`
- **THEN** location record stores (name: `archlinux`, type: `local`)

#### Scenario: Remote scan location
- **WHEN** user indexes rclone remote `dropbox`
- **THEN** location record stores (name: `dropbox`, type: `remote`)

### Requirement: Unique Remote Identity
The system SHALL enforce unique identity of remotes using (location_id, remote_name, scan_root).

#### Scenario: Multi-root prevention
- **WHEN** same filesystem is scanned from `/media` and `/media/2024`
- **THEN** separate remote records are created (different scan_roots)
- **AND** no false duplicates result from path variance

## MODIFIED Requirements

### Requirement: Database Schema
The system SHALL maintain remote metadata linked to location sources instead of hostname columns.

#### Scenario: Remote-location relationship
- **WHEN** persisting a remote
- **THEN** remote references location by ID, not by hostname string
- **AND** type (local vs remote) is explicit in location record
