## ADDED Requirements

### Requirement: Configured Remote Enumeration
The system SHALL support the rclone RC `config/listremotes` endpoint for discovering local remote names.

#### Scenario: Enumerate configured remote names
- **WHEN** Filesprawl requests configured remotes from the active local `rcd` session
- **THEN** it calls the `config/listremotes` endpoint
- **AND** parses the returned remote names into a structured result for CLI use

#### Scenario: Propagate configured remote enumeration errors
- **WHEN** the `config/listremotes` endpoint returns an error or non-200 response
- **THEN** Filesprawl wraps and returns the remote enumeration failure with rclone RC context