## ADDED Requirements

### Requirement: Remote Association Persistence
<<<<<<< HEAD
The system SHALL persist and reuse remote records and object-to-remote associations during indexed scans.

#### Scenario: Ensure remote exists before linking objects
=======

The system SHALL persist and reuse remote records and object-to-remote associations during indexed scans.

#### Scenario: Ensure remote exists before linking objects

>>>>>>> 562cedc (Add index command and remote discovery)
- **WHEN** the repository persists scan results for a selected host/remote identity
- **THEN** it reads or writes the corresponding remote record before creating object-to-remote associations

#### Scenario: Link indexed objects to the selected remote
<<<<<<< HEAD
=======

>>>>>>> 562cedc (Add index command and remote discovery)
- **WHEN** file metadata is persisted from an indexed scan
- **THEN** the repository creates an `object_remote_junction` record linking the object metadata to the selected remote identity

#### Scenario: Reuse existing persisted remote identity
<<<<<<< HEAD
- **WHEN** additional files are indexed from the same host and local rclone remote name
- **THEN** the repository reuses the existing remote record instead of creating a duplicate remote identity
=======

- **WHEN** additional files are indexed from the same host and local rclone remote name
- **THEN** the repository reuses the existing remote record instead of creating a duplicate remote identity
>>>>>>> 562cedc (Add index command and remote discovery)
