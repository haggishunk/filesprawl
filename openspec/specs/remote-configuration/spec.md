# Remote Configuration

## Purpose
Configuration and management of remote storage locations with hostname and name associations.

## Requirements

### Requirement: Remote Type
The system SHALL provide a Remote type with hostname and name fields.

#### Scenario: Remote creation
- **WHEN** creating a Remote
- **THEN** it includes Hostname and Name fields

#### Scenario: Hostname identification
- **WHEN** a Remote is configured
- **THEN** Hostname identifies the machine/host accessing the remote

#### Scenario: Remote name
- **WHEN** a Remote is configured
- **THEN** Name identifies the rclone remote configuration (e.g., "dbox:")

### Requirement: Multi-Host Support
The system SHALL support different hosts having the same remote with different names.

#### Scenario: Host-specific naming
- **WHEN** multiple hosts access the same storage backend
- **THEN** each can use different rclone remote names

#### Scenario: Remote identification
- **WHEN** identifying a remote
- **THEN** both hostname and name are used for uniqueness

### Requirement: Database Schema
The system SHALL maintain a remote table in the database.

#### Scenario: Remote table structure
- **WHEN** the database is initialized
- **THEN** remote table includes id, remote_name, remote_type, and hostname fields

#### Scenario: Remote persistence
- **WHEN** storing remote configurations
- **THEN** they are persisted with name, type, and hostname

