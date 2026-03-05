# rclone Integration

## Purpose
Integration with rclone RC (remote control) server for unified access to multiple remote storage backends (Dropbox, S3, Google Drive, etc.).

## Requirements

### Requirement: RC Server Communication
The system SHALL communicate with rclone RC server via HTTP POST with JSON payloads.

#### Scenario: RC server endpoint
- **WHEN** making rclone operations
- **THEN** requests are sent to http://localhost:5572

#### Scenario: JSON request encoding
- **WHEN** calling rclone operations
- **THEN** configuration is encoded as JSON in the request body

#### Scenario: JSON response decoding
- **WHEN** receiving rclone responses
- **THEN** responses are decoded from JSON into structured types

### Requirement: Operations List Support
The system SHALL support the operations/list endpoint for directory and file listing.

#### Scenario: List files
- **WHEN** listing with FilesOnly option
- **THEN** only files are returned in the response

#### Scenario: List directories
- **WHEN** listing with DirsOnly option
- **THEN** only directories are returned in the response

#### Scenario: Hash retrieval
- **WHEN** listing with ShowHash enabled
- **THEN** file hashes are included in the response

### Requirement: List Configuration
The system SHALL provide structured configuration for list operations.

#### Scenario: Remote and path specification
- **WHEN** creating a ListConfig
- **THEN** it includes fs (remote name) and remote (path) fields

#### Scenario: List options
- **WHEN** configuring list operations
- **THEN** options include DirsOnly, FilesOnly, HashTypes, ShowHash, Recurse, etc.

### Requirement: Response Parsing
The system SHALL parse rclone responses into structured types.

#### Scenario: List response items
- **WHEN** receiving list responses
- **THEN** items include ID, Name, Path, Size, IsDir, MimeType, ModTime, and Hashes

#### Scenario: Hash map parsing
- **WHEN** files have multiple hash types
- **THEN** hashes are parsed as a map[string]string

### Requirement: Error Handling
The system SHALL handle rclone RC errors with context.

#### Scenario: Connection failure
- **WHEN** rclone RC server is unavailable
- **THEN** return error with "connection failed" context

#### Scenario: HTTP error status
- **WHEN** rclone returns non-200 status
- **THEN** return error with operation name and error details

### Requirement: Authentication Support
The system SHALL support optional basic authentication for RC server.

#### Scenario: No authentication
- **WHEN** authUser and authPass are empty
- **THEN** requests are sent without authentication headers

#### Scenario: Basic auth
- **WHEN** authUser and authPass are configured
- **THEN** requests include basic authentication headers

