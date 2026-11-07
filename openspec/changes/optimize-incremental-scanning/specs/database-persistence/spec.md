## ADDED Requirements

### Requirement: Change Tracking Fields

The system SHALL persist file change tracking information in the database to support delta scanning.

#### Scenario: Store mtime in object metadata

- **WHEN** a file is scanned and hashed
- **THEN** the file's modification time (mtime) is persisted in object_meta table

#### Scenario: Store scan timestamp

- **WHEN** a file is indexed during a scan
- **THEN** a scan timestamp is recorded to indicate when the file was last scanned

#### Scenario: Retrieve change tracking data

- **WHEN** delta scanning queries the database for change information
- **THEN** mtime and scan timestamp data is returned for comparison

### Requirement: Deletion Tracking

The system SHALL track files deleted since the last scan to support accurate delta detection.

#### Scenario: Mark deleted files

- **WHEN** delta scanning detects that a previously indexed file no longer exists
- **THEN** the file is marked as deleted in the database (soft delete)

#### Scenario: Exclude deleted files from duplicates

- **WHEN** duplicate detection runs after deletion tracking
- **THEN** deleted files are excluded from duplicate results

#### Scenario: Query active files

- **WHEN** delta scanning queries for files to compare mtimes
- **THEN** only non-deleted files are returned from the database

### Requirement: Efficient Change Detection Indexes

The system SHALL maintain database indexes to efficiently query changed files during delta scans.

#### Scenario: Index on mtime column

- **WHEN** delta scanning queries files to detect modifications
- **THEN** database uses an index on object_meta.mtime for fast retrieval

#### Scenario: Index on scan timestamp

- **WHEN** identifying newly added files since last scan
- **THEN** database uses an index on scan timestamp for efficient filtering

#### Scenario: Index on deletion status

- **WHEN** querying active files during delta scanning
- **THEN** database uses an index on deletion flag for fast filtering

### Requirement: Schema Migration

The system SHALL provide database schema migrations to add delta scanning support.

#### Scenario: Add mtime column

- **WHEN** upgrading the database schema
- **THEN** migration adds mtime timestamp column to object_meta table

#### Scenario: Add deletion tracking column

- **WHEN** upgrading the database schema
- **THEN** migration adds is_deleted boolean column to object_meta with default false

#### Scenario: Create change detection indexes

- **WHEN** upgrading the database schema
- **THEN** migration creates indexes on mtime, scan_timestamp, and is_deleted columns

## MODIFIED Requirements

### Requirement: Object Metadata Persistence

The system SHALL persist complete object metadata including change tracking fields for delta scanning support.

#### Scenario: Persist metadata with mtime

- **WHEN** writing object metadata to the database
- **THEN** mtime and scan timestamp are included with name, path, and MIME type

#### Scenario: Update existing metadata

- **WHEN** a file is re-scanned and found to be modified
- **THEN** existing metadata record is updated with new mtime and current scan timestamp

#### Scenario: Retrieve metadata with change history

- **WHEN** querying metadata for delta scanning comparison
- **THEN** database returns all metadata including mtime and scan timestamp for change detection
