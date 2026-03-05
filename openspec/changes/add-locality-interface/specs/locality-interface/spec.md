## ADDED Requirements

### Requirement: Locality Configuration Storage
The system SHALL store locality mapping configuration in remote storage and support local overrides.

#### Scenario: Remote configuration storage
- **WHEN** locality configuration is stored
- **THEN** it is saved to `.filesprawl/locality.json` on the remote

#### Scenario: Local configuration override
- **WHEN** local override configuration exists
- **THEN** it takes precedence over remote configuration

#### Scenario: Configuration merging
- **WHEN** both remote and local configurations exist
- **THEN** local rules override remote rules with the same pattern

### Requirement: Path Mapping Rules
The system SHALL support configurable rules to map remote paths to local paths.

#### Scenario: Simple path mapping
- **WHEN** a mapping rule is defined (e.g., `/documents` -> `/home/user/Documents`)
- **THEN** remote paths under `/documents` are mapped to local paths under `/home/user/Documents`

#### Scenario: Pattern-based mapping
- **WHEN** a mapping rule uses patterns (e.g., `/{category}/*` -> `/home/user/{category}/*`)
- **THEN** variable substitution is applied during path resolution

#### Scenario: Multiple mapping rules
- **WHEN** multiple mapping rules are defined
- **THEN** the most specific matching rule is applied

### Requirement: Remote to Local Resolution
The system SHALL resolve remote storage paths to local filesystem paths.

#### Scenario: Resolve remote path
- **WHEN** given a remote path and locality configuration
- **THEN** return the corresponding local filesystem path

#### Scenario: No mapping found
- **WHEN** no mapping rule matches a remote path
- **THEN** return error indicating no mapping exists

#### Scenario: Hostname-specific mapping
- **WHEN** different hosts have different local paths
- **THEN** apply hostname-specific mapping rules

### Requirement: Local to Remote Resolution
The system SHALL resolve local filesystem paths to remote storage paths.

#### Scenario: Resolve local path
- **WHEN** given a local path and locality configuration
- **THEN** return the corresponding remote storage path

#### Scenario: Multiple remote matches
- **WHEN** a local path could map to multiple remotes
- **THEN** return all matching remote paths

#### Scenario: Reverse pattern matching
- **WHEN** resolving local to remote with patterns
- **THEN** extract variables from local path and apply to remote pattern

### Requirement: Configuration Caching
The system SHALL cache locality configuration to minimize remote reads.

#### Scenario: Configuration cache
- **WHEN** locality configuration is read from remote
- **THEN** it is cached for subsequent lookups

#### Scenario: Cache invalidation
- **WHEN** configuration is updated
- **THEN** cache is invalidated and refreshed on next access

#### Scenario: Cache TTL
- **WHEN** cached configuration exceeds time-to-live
- **THEN** configuration is re-read from remote

