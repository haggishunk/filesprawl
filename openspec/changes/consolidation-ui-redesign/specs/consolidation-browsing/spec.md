## MODIFIED Requirements

### Requirement: Dual-View Duplicate Browsing

The system SHALL provide two complementary browsing views for duplicate consolidation accessed via side-by-side panels: hash-centric list (left) and file details table (right), allowing users to understand relationships and make informed decisions.

#### Scenario: User browses hashes in left panel

- **WHEN** user launches consolidation and sees duplicate groups
- **THEN** left panel shows list of unique hashes with:
  - Safety indicator (🔴/🟡/🟢)
  - Hash prefix (first 8 chars)
  - Copy count in parentheses "(5 copies)"
  - Location count indicator "📍 3 loc"
- **AND** user can scroll with j/k to select different hash groups
- **AND** selection highlights current hash in left panel

#### Scenario: User inspects files in right panel

- **WHEN** user selects a hash group in left panel
- **THEN** right panel shows table of all files in that group with:
  - Columns: remote | filename | size | status
  - Status showing safety (✓ safe in N locations, 🟡 caution, 🔴 only 1 location)
  - Flag indicators if KEEP/ARCHIVE/DELETE already marked
- **AND** user can scroll with j/k to select individual files
- **AND** pressing 'v' toggles "wide" view to add filepath column

#### Scenario: User flags files for action

- **WHEN** user selects a file in right panel
- **THEN** user can press K/a/d/s to flag with:
  - K = KEEP this copy
  - a = ARCHIVE to backup
  - d = DELETE (safety checked)
  - s = SKIP/clear flags
- **AND** flag status appears in table immediately

### Requirement: Deletion Safety Constraint

The system SHALL prevent deletion of all copies of a (remote, filename) pair and inform users of safety violations.

#### Scenario: User attempts unsafe deletion

- **WHEN** user tries to flag multiple files with same (remote, filename) as DELETE
- **THEN** if only 1 copy exists for that (remote, filename) pair
- **AND** user flags it DELETE
- **THEN** system shows error: "Cannot delete only copy of [remote] / [filename]"
- **AND** flag is rejected and table not updated

#### Scenario: User safely deletes redundant copies

- **WHEN** same filename exists on multiple remotes
- **THEN** user can DELETE copies on secondary remotes
- **AND** as long as at least 1 copy per (remote, filename) remains unflagged
- **AND** deletion is allowed and table shows DELETE status
