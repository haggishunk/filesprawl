## ADDED Requirements

### Requirement: Remote Discovery Command
The system SHALL provide a `list-remotes` CLI command that lists the local rclone remote names available through the active `rcd` session.

#### Scenario: List locally configured remotes
- **WHEN** a user runs `filesprawl list-remotes`
- **THEN** the system calls the rclone RC endpoint for listing configured remotes
- **AND** returns the remote names reported by the active local `rcd` session

#### Scenario: Surface remote discovery failures
- **WHEN** the rclone RC remote-list operation fails
- **THEN** the command returns an error that explains remote discovery could not be completed

### Requirement: Index Command Invocation
The system SHALL provide an `index` CLI command for starting a scan against a selected rclone remote.

#### Scenario: Index the root of a selected remote
- **WHEN** a user runs `filesprawl index --remote media`
- **THEN** the system starts a scan using `media:` as the rclone `fs` value
- **AND** the scan targets the root path for that remote

#### Scenario: Index a selected remote subpath
- **WHEN** a user runs `filesprawl index --remote media --path code/flux`
- **THEN** the system starts a scan using `media:` as the rclone `fs` value
- **AND** the scan targets `code/flux` as the remote path

### Requirement: Local RCD Remote Resolution
The system SHALL interpret the CLI remote argument as a remote name configured on the same host as the active rclone `rcd` session.

#### Scenario: Forward local remote name to rclone RC
- **WHEN** the application invokes the scan pipeline for `--remote media`
- **THEN** the local remote name is normalized for persistence
- **AND** `media:` is passed to the rclone RC list operation used by Filesprawl

#### Scenario: Missing remote argument
- **WHEN** a user invokes `filesprawl index` without specifying a remote
- **THEN** the command fails with a validation error explaining that a remote name is required