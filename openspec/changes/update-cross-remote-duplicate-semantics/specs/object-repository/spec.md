## ADDED Requirements

### Requirement: Cross-Remote Duplicate Distinctness
The system SHALL only report a hash as a cross-remote duplicate when matching files are linked to at least two distinct persisted remote identities.

#### Scenario: Same hash on two files in one remote
- **WHEN** two files in the same persisted remote share the same hash
- **THEN** that hash is not returned by the cross-remote duplicate query

#### Scenario: Same hash on two distinct remotes
- **WHEN** files linked to two distinct persisted remote identities share the same hash
- **THEN** that hash is returned by the cross-remote duplicate query

#### Scenario: Same remote name on different hosts
- **WHEN** two files share the same hash and are linked to remote records with the same remote name but different hostnames
- **THEN** the hash qualifies as a cross-remote duplicate because the persisted remote identities are distinct