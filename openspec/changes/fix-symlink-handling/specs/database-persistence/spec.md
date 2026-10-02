## MODIFIED Requirements

### Requirement: Database Schema
The system SHALL maintain an extended PostgreSQL schema with symlink and scan origin tracking.

#### Scenario: is_symlink column
- **WHEN** storing object metadata
- **THEN** a boolean `is_symlink` column records whether the file is a symlink

#### Scenario: scan_root column
- **WHEN** storing object metadata from local file scanning
- **THEN** a text `scan_root` column records the directory root of the scan

#### Scenario: Schema migration
- **WHEN** upgrading from previous schema
- **THEN** migration adds is_symlink (default false) and scan_root columns to object_meta

### Requirement: Indexed Queries
The system SHALL efficiently query files accounting for scan origin.

#### Scenario: Find by hash and scan_root
- **WHEN** querying duplicates across multiple scan roots
- **THEN** queries return files grouped by hash while respecting scan_root distinctness

#### Scenario: Same name, different roots
- **WHEN** two files have same path but different scan_root values
- **THEN** they are treated as distinct records for duplicate matching
