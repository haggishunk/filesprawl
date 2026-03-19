## ADDED Requirements

### Requirement: Indexed Remote Identity
The system SHALL persist the identity of an indexed remote as the host computer together with the local rclone remote name used by the active `rcd` session.

#### Scenario: Persist host and local remote name for a scan
- **WHEN** Filesprawl indexes remote `dbox:` from host `guru`
- **THEN** the persisted remote identity records hostname `guru` and remote name `dbox:` for that scan context

#### Scenario: Distinguish same remote name across hosts
- **WHEN** two different host computers each index a locally configured remote named `dbox:`
- **THEN** the persisted remote identities remain distinct by hostname

#### Scenario: Use local configuration context
- **WHEN** a scan is launched through a local rclone `rcd` session
- **THEN** the persisted remote identity reflects the local rclone configuration context rather than assuming a backend-global remote identifier