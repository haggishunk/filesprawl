# Database Persistence

## Purpose
PostgreSQL-based persistence layer for file indexing data including object hashes, metadata, and their relationships.

## Requirements

### Requirement: Database Schema
The system SHALL maintain a PostgreSQL database schema with tables for object hashes, metadata, remotes, and junction tables.

#### Scenario: Schema initialization
- **WHEN** the database is initialized
- **THEN** tables for object_hash, object_meta, remote, object_hash_junction, and object_remote_junction are created

#### Scenario: Hash type enumeration
- **WHEN** storing hash values
- **THEN** hash_type is constrained to enum values: md5, dropbox, sha1, sha256

### Requirement: Connection Pooling
The system SHALL use pgxpool for efficient database connection management.

#### Scenario: Pool lifecycle
- **WHEN** the application starts
- **THEN** a connection pool is created and maintained for the application lifetime

#### Scenario: Connection logging
- **WHEN** a database connection is established
- **THEN** the connection PID is logged

### Requirement: Database Interface Abstraction
The system SHALL provide a Database interface for testability and flexibility.

#### Scenario: Interface implementation
- **WHEN** database operations are performed
- **THEN** they use the Database interface (QueryRow, Exec methods)

#### Scenario: PostgreSQL implementation
- **WHEN** using PostgreSQL
- **THEN** PgxDatabase struct implements the Database interface

### Requirement: Context Propagation
The system SHALL accept context.Context for all database operations.

#### Scenario: Context-aware queries
- **WHEN** executing database queries
- **THEN** context is passed to enable cancellation and timeout handling

### Requirement: Foreign Key Relationships
The system SHALL enforce referential integrity through foreign key constraints.

#### Scenario: Junction table constraints
- **WHEN** inserting into object_hash_junction
- **THEN** object_meta_id and object_hash_id MUST reference valid records

#### Scenario: Remote junction constraints
- **WHEN** inserting into object_remote_junction
- **THEN** object_meta_id and remote_id MUST reference valid records

