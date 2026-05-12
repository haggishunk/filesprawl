## ADDED Requirements

### Requirement: Duplicate Detection Within Remote
The system SHALL identify files with identical content hashes within a single remote storage location.

#### Scenario: Find duplicates in single remote
- **WHEN** querying for duplicates within a specific remote
- **THEN** return groups of files that share the same hash value

#### Scenario: Multiple hash types
- **WHEN** files have multiple hash types (md5, sha1, sha256)
- **THEN** duplicates are identified by any matching hash type

#### Scenario: No duplicates found
- **WHEN** no files share hash values within a remote
- **THEN** return empty result set

### Requirement: Duplicate Detection Across Remotes
The system SHALL identify files with identical content hashes across multiple remote storage locations.

#### Scenario: Find duplicates across all remotes
- **WHEN** querying for duplicates across all remotes
- **THEN** return groups of files from different remotes that share the same hash value

#### Scenario: Remote identification
- **WHEN** duplicates are found across remotes
- **THEN** each file includes remote identification (hostname, remote name)

#### Scenario: Cross-remote deduplication candidates
- **WHEN** a file exists on multiple remotes
- **THEN** identify which copies could be removed to save storage

### Requirement: Duplicate Grouping
The system SHALL group duplicate files by their content hash.

#### Scenario: Hash-based grouping
- **WHEN** retrieving duplicates
- **THEN** files are grouped by hash value with all matching files in each group

#### Scenario: Group metadata
- **WHEN** displaying duplicate groups
- **THEN** include hash value, hash type, and count of duplicates

#### Scenario: File details in groups
- **WHEN** displaying files within a duplicate group
- **THEN** include file name, path, remote, and size for each file

### Requirement: Duplicate Filtering
The system SHALL support filtering duplicate detection by criteria.

#### Scenario: Filter by remote
- **WHEN** filtering duplicates by remote
- **THEN** only include files from specified remotes

#### Scenario: Filter by hash type
- **WHEN** filtering duplicates by hash type
- **THEN** only use specified hash types for duplicate detection

#### Scenario: Filter by minimum file size
- **WHEN** filtering duplicates by minimum file size
- **THEN** only include files above the specified size threshold

### Requirement: Performance Optimization
The system SHALL optimize duplicate detection queries for large datasets.

#### Scenario: Hash value indexing
- **WHEN** querying for duplicates
- **THEN** use database indexes on hash_value for efficient lookups

#### Scenario: Pagination support
- **WHEN** duplicate result sets are large
- **THEN** support pagination to retrieve results in chunks

