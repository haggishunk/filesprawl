## ADDED Requirements

### Requirement: Interactive Consolidation Command
The system SHALL provide a `consolidate` command allowing users to discover, preview, and consolidate duplicate files interactively.

#### Scenario: User discovers and consolidates duplicates
- **WHEN** user runs `filesprawl consolidate --remote media`
- **THEN** system presents first duplicate group with image previews (if available)
- **AND** user can decide: keep here / consolidate to destination / skip / batch mode
- **AND** consolidation moves files to archive location preserving directory structure

#### Scenario: Batch consolidation with preferences
- **WHEN** user runs `filesprawl consolidate --remote media --batch`
- **THEN** system applies stored consolidation preferences for known patterns
- **AND** only ambiguous groups are presented for interactive decision

### Requirement: Consolidation Memory
The system SHALL persist user consolidation decisions and apply them to similar duplicates in future runs.

#### Scenario: Preference reuse
- **WHEN** user consolidates 10 groups to archive-s3 location
- **THEN** system learns preference: "duplicates in media → archive to s3"
- **AND** future consolidations of media duplicates auto-apply this preference

### Requirement: Archive Operations with Path Preservation
The system SHALL move consolidated files to archive location while preserving original directory structure for recovery.

#### Scenario: Successful archival
- **WHEN** user consolidates `photos/2020/pic.jpg` to archive
- **THEN** file is copied to `archive/media//photos/2020/pic.jpg`
- **AND** original file location is recorded in consolidation history
- **AND** user can recover archived files by path if needed
