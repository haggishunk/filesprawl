# Object Repository

## Purpose
Repository pattern implementation for CRUD operations on file objects, hashes, and their relationships.

## Requirements

### Requirement: Hash Operations
The system SHALL provide read and write operations for object hashes.

#### Scenario: Write new hash
- **WHEN** writing a hash that doesn't exist
- **THEN** it is inserted and the ID is returned

#### Scenario: Read existing hash
- **WHEN** reading a hash by value and type
- **THEN** the ID is retrieved and Persisted flag is set to true

#### Scenario: Read non-existent hash
- **WHEN** reading a hash that doesn't exist
- **THEN** no error is returned and Persisted flag remains false

### Requirement: Metadata Operations
The system SHALL provide read and write operations for object metadata.

#### Scenario: Write new metadata
- **WHEN** writing metadata that doesn't exist
- **THEN** it is inserted with name, path, and MIME type

#### Scenario: Read existing metadata
- **WHEN** reading metadata by name, path, and MIME type
- **THEN** the ID is retrieved and Persisted flag is set to true

#### Scenario: Read non-existent metadata
- **WHEN** reading metadata that doesn't exist
- **THEN** no error is returned and Persisted flag remains false

### Requirement: Junction Operations
The system SHALL manage meta-hash junction records with scan timestamps.

#### Scenario: Write new junction
- **WHEN** writing a junction that doesn't exist
- **THEN** it is inserted with meta_id, hash_id, and current timestamp

#### Scenario: Read existing junction
- **WHEN** reading a junction by meta_id and hash_id
- **THEN** the ID and scan_time are retrieved

#### Scenario: Read non-existent junction
- **WHEN** reading a junction that doesn't exist
- **THEN** no error is returned and Persisted flag remains false

### Requirement: Integrated Persistence
The system SHALL provide integrated persistence for rclone list results.

#### Scenario: Persist list response item
- **WHEN** persisting a ListResponseItem
- **THEN** metadata is read or written first

#### Scenario: Hash iteration
- **WHEN** persisting a ListResponseItem with multiple hashes
- **THEN** each hash type is processed and persisted

#### Scenario: Junction creation
- **WHEN** both metadata and hash are persisted
- **THEN** a junction record is created linking them

### Requirement: Idempotent Operations
The system SHALL support idempotent read-or-write patterns.

#### Scenario: Duplicate hash write
- **WHEN** attempting to write an existing hash
- **THEN** the existing ID is used without error

#### Scenario: Duplicate metadata write
- **WHEN** attempting to write existing metadata
- **THEN** the existing ID is used without error

### Requirement: Error Handling
The system SHALL wrap database errors with context.

#### Scenario: Query error
- **WHEN** a database query fails
- **THEN** error is wrapped with operation context using fmt.Errorf with %w

#### Scenario: No rows handling
- **WHEN** a query returns no rows
- **THEN** it is handled gracefully without error for read operations

