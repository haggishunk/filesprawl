## ADDED Requirements

### Requirement: Dual-View Duplicate Browsing
The system SHALL provide two complementary browsing views for duplicate consolidation: hash-centric and filename-centric, allowing users to understand relationships from both data and organizational perspectives.

#### Scenario: User browses by hash to see all files with same content
- **WHEN** user selects "View by Hash" in consolidation interface
- **THEN** system groups all identical files by content hash
- **AND** displays each hash with count of copies: "Hash abc123 (47 copies)"
- **AND** user can expand hash to see all files and their locations
- **AND** selecting a hash shows all files grouped by location (remote/local name)

#### Scenario: User browses by filename to see different versions
- **WHEN** user selects "View by Filename" in consolidation interface
- **THEN** system groups files by object_name (filename only)
- **AND** displays each filename with count of unique hashes: "photo.jpg (5 different versions)"
- **AND** expanding filename shows each unique hash variant
- **AND** user can see count and locations for each variant

### Requirement: Location Presence Indicators with Color Coding
The system SHALL display file locations (local/remote names) with visual presence indicators to show backup replication status and identify single points of failure.

#### Scenario: User sees complete replication (green status)
- **WHEN** user views a file in hash-by-view showing locations
- **THEN** system displays each location with presence indicator:
  - "🟢 archlinux//goblet/light (local) — 3 copies"
  - "🟢 archlinux/do (dropbox) — 2 copies"
  - "🟢 s3-backup/archive (remote) — 1 copy"
- **AND** green indicator means "backup represented in this location"
- **AND** user can verify coverage before consolidation

#### Scenario: User identifies risky single locations (red status)
- **WHEN** a file exists in only one location with one copy
- **THEN** system displays warning indicator:
  - "🔴 archlinux//goblet (local) — ONLY copy!"
  - "⚠️ This file has no backup elsewhere"
- **AND** red/yellow triggers data loss warning
- **AND** system recommends archiving to additional location first

#### Scenario: User sees mixed status (yellow warning)
- **WHEN** file has limited replication (e.g., only 2 locations instead of 3+)
- **THEN** system displays intermediate indicator:
  - "🟡 archlinux/do (dropbox) — 1 copy (limited)"
- **AND** yellow indicates "present but limited backup"

### Requirement: File Flagging for Consolidation Decisions
The system SHALL allow users to flag individual files for keep/delete/archive actions and persist these decisions across browsing sessions.

#### Scenario: User flags files for archival
- **WHEN** user browses hash or filename view
- **THEN** user can select individual file with keyboard shortcut (e.g., 'a' for archive)
- **AND** file displays visual flag marker: "▶ [ARCHIVE] photo.jpg in archlinux//goblet/light"
- **AND** flag persists while user continues browsing other files
- **AND** user can toggle flag with same shortcut

#### Scenario: User flags file to keep in place
- **WHEN** user views duplicate group with multiple variants
- **THEN** user can flag preferred version with shortcut (e.g., 'k' for keep)
- **AND** flagged file shows marker: "✓ [KEEP] photo.jpg in s3-backup (best quality)"
- **AND** other versions automatically suggest archive/delete

#### Scenario: User flags file for deletion
- **WHEN** user identifies corrupted or 0-byte file
- **THEN** user can flag for deletion with shortcut (e.g., 'd' for delete)
- **AND** file shows marker: "✗ [DELETE] photo.jpg (0 bytes, corrupted)"
- **AND** deletion requires explicit confirmation before execution
- **AND** system verifies backup exists before allowing deletion

#### Scenario: User reviews and modifies flags before execution
- **WHEN** user completes browsing and flags decisions
- **THEN** system shows summary: "Ready to consolidate: Archive 23, Keep 5, Delete 3"
- **AND** user can review each flagged file before executing
- **AND** user can modify individual flags or clear all
- **AND** consolidation only executes on explicitly flagged files

### Requirement: Consolidated Decision Summary Before Execution
The system SHALL present a final summary of all user decisions grouped by action type before executing any consolidation operations.

#### Scenario: User reviews final decision summary
- **WHEN** user requests summary after flagging decisions
- **THEN** system displays consolidated view:
```
📋 CONSOLIDATION SUMMARY

KEEP (5 files):
  ✓ photo.jpg (s3-backup) — 2.5 MB, complete
  ✓ photo.jpg (archlinux/do) — 2.5 MB, complete
  ...

ARCHIVE (23 files):
  ▶ old-photo.jpg (archlinux//goblet) → archive-s3
  ▶ incomplete-video.mp4 (archive-old) → s3-backup
  ...

DELETE (3 files):
  ✗ corrupted.jpg (0 bytes, local copy exists)
  ✗ duplicate.jpg (broken, 3+ good copies exist)
  ...

[Confirm] [Cancel] [Modify]
```
- **AND** user can review each action with full file details
- **AND** system performs final safety check before execution
