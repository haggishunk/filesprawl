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

### Requirement: Name Collision Risk Detection
The system SHALL detect when files with identical names contain different content (different hashes) and alert the user to prevent data loss from filename-based consolidation.

#### Scenario: Multiple versions of same filename exist
- **WHEN** user views a duplicate group containing `01 Intro.mp3` with 3 different hashes
- **THEN** system displays warning: "⚠️ NAME COLLISION: 3 different versions of '01 Intro.mp3' exist"
- **AND** system shows count and status of each variant: "✅ All 3 variants preserved"
- **AND** system recommends hash-based consolidation logic instead of filename-based

#### Scenario: Some variants of colliding filename are missing
- **WHEN** user views group with `cover.jpg` but only 1 of 5 versions is preserved
- **THEN** system displays critical warning: "🚨 CRITICAL: 4 variants of 'cover.jpg' are MISSING (data loss detected!)"
- **AND** system blocks consolidation until user manually verifies
- **AND** system recommends backup/recovery procedures

### Requirement: Zero-Byte and Incomplete File Detection
The system SHALL categorize duplicate files by completeness status (zero-byte vs non-zero) and verify copies exist elsewhere before allowing consolidation.

#### Scenario: Mixed zero-byte and complete copies
- **WHEN** user views duplicate group with `photo.jpg` (2 zero-byte, 5 complete copies)
- **THEN** system displays status: "✅ 2 zero-byte (corrupted/incomplete) — safe to delete"
- **AND** system shows: "✅ 5 complete copies exist elsewhere — safe to consolidate"
- **AND** system verifies at least one non-zero copy exists in another location before consolidation

#### Scenario: Only zero-byte copies in current location
- **WHEN** user views group where all copies in target location are 0 bytes
- **THEN** system checks if non-zero copies exist elsewhere
- **AND** if yes: displays "✅ Safe: Zero-byte corrupted files, complete copies exist elsewhere"
- **AND** if no: displays "🚨 CRITICAL: No complete copies exist anywhere!" and blocks consolidation

### Requirement: Directory Deletion Safety Information
The system SHALL verify that all files in a consolidation directory exist elsewhere before presenting consolidation as safe.

#### Scenario: User consolidates entire directory
- **WHEN** user runs `filesprawl consolidate --directory ~/photos/old-backup`
- **THEN** system analyzes all files in directory and reports:
  - Count of files with copies elsewhere (safe to delete)
  - Count of zero-byte corrupted files (safe to delete)
  - Count of risky files (only copy exists here) — blocks consolidation
- **AND** system displays summary: "✅ SAFE TO CONSOLIDATE: All 1,204 files represented elsewhere"

#### Scenario: Risky files detected in consolidation directory
- **WHEN** directory contains files that only exist in that location
- **THEN** system displays: "⚠️ RISKY: 3 file(s) only exist in this directory. These will be LOST if consolidation proceeds!"
- **AND** system requires explicit user confirmation with file listing
- **AND** system suggests archiving to backup location as alternative

### Requirement: Safety Information Display in Interactive Consolidation
The system SHALL present all data loss prevention information (name collisions, completeness status, directory safety) to user during interactive consolidation workflow for informed decision-making.

#### Scenario: User reviews safety information before consolidation
- **WHEN** user reviews a duplicate group in interactive mode
- **THEN** system displays consolidated safety panel:
  - "📋 Group Summary: 47 copies across 3 remotes"
  - "🔗 Filename Status: ✅ No collisions (all versions preserved)"
  - "💾 Completeness: ✅ All 47 copies non-zero, 1+ in each location"
  - "📂 Directory Safety: ✅ Target directory contents all exist elsewhere"
- **AND** user can request detailed breakdown for each category
- **AND** safety information influences consolidation recommendation (green/yellow/red status)
