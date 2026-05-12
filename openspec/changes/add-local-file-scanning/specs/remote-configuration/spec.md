## ADDED Requirements

### Requirement: Local Source Identity
The system SHALL represent locally scanned file origins as persisted host-qualified source identities.

#### Scenario: Persist local source identity
- **WHEN** a local filesystem root is indexed on a host
- **THEN** the persisted source identity records the hostname, the normalized local root identifier, and source type `local`

#### Scenario: Distinguish same local root across hosts
- **WHEN** two hosts each index the same local path string
- **THEN** the persisted local source identities remain distinct by hostname