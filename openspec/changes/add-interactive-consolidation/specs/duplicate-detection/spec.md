## MODIFIED Requirements

### Requirement: Duplicate Group Retrieval with Consolidation Status
The system SHALL return duplicate groups with consolidation status and history when requested.

#### Scenario: Consolidation status included
- **WHEN** user queries duplicates with `--with-consolidation-status`
- **THEN** each group includes: prior consolidation decisions, archive destinations, preference matches
- **AND** user knows which duplicates were already handled in previous runs

#### Scenario: Filter consolidated groups
- **WHEN** user runs `filesprawl duplicates --exclude-consolidated`
- **THEN** system returns only new/unprocessed duplicate groups
- **AND** previously consolidated groups are hidden from view

## ADDED Requirements

### Requirement: Consolidation-Aware Duplicate Filtering
The system SHALL filter duplicates based on consolidation history to avoid re-processing.

#### Scenario: Skip already consolidated
- **WHEN** duplicate group matches prior consolidation preference
- **THEN** system can skip that group in batch mode if preference is still valid
