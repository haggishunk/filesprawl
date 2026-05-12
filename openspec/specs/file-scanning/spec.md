# File Scanning

## Purpose
Recursive directory traversal and file indexing across remote storage locations using rclone integration.

## Requirements

### Requirement: Scanner Configuration
The system SHALL provide a Scanner type configured with remote and repository.

#### Scenario: Scanner creation
- **WHEN** creating a Scanner
- **THEN** it is configured using functional options (WithRemote, WithRepo)

#### Scenario: Remote association
- **WHEN** Scanner is configured with a remote
- **THEN** it targets that specific remote storage location

#### Scenario: Repository association
- **WHEN** Scanner is configured with a repository
- **THEN** it uses that repository for persistence operations

### Requirement: Recursive Directory Traversal
The system SHALL recursively traverse directories on remote storage.

#### Scenario: Directory descent
- **WHEN** scanning with DirsOnly option
- **THEN** each directory is recursively scanned for files

#### Scenario: File processing
- **WHEN** scanning with FilesOnly option
- **THEN** files are processed and persisted before descending into subdirectories

#### Scenario: Depth-first traversal
- **WHEN** scanning a directory tree
- **THEN** files in current directory are processed before descending into subdirectories

### Requirement: File Persistence
The system SHALL persist file information during scanning.

#### Scenario: File metadata persistence
- **WHEN** a file is encountered during scanning
- **THEN** its metadata (name, path, MIME type) is persisted

#### Scenario: Hash persistence
- **WHEN** a file has hash values
- **THEN** each hash type and value is persisted

#### Scenario: Junction creation
- **WHEN** file metadata and hashes are persisted
- **THEN** a junction record links them with a scan timestamp

### Requirement: Scan Context
The system SHALL accept context for scan operations.

#### Scenario: Context propagation
- **WHEN** scanning is initiated
- **THEN** context is propagated to all rclone and database operations

#### Scenario: Cancellation support
- **WHEN** context is cancelled during scanning
- **THEN** scan operations are terminated

### Requirement: Error Handling
The system SHALL handle and propagate errors during scanning.

#### Scenario: List operation failure
- **WHEN** rclone list operation fails
- **THEN** error is wrapped with context and returned

#### Scenario: Persistence failure
- **WHEN** database persistence fails
- **THEN** error is wrapped with context and returned

#### Scenario: Recursive scan failure
- **WHEN** scanning a subdirectory fails
- **THEN** error is propagated up the call stack

