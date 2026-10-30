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

### Requirement: User Authentication and Session Management
The system SHALL require users to authenticate or create accounts before consolidation sessions, enabling multi-user tracking and collaborative review workflows.

#### Scenario: New user starts consolidation session
- **WHEN** user runs `filesprawl consolidate --remote media` for first time
- **THEN** system prompts for username: "Enter your name (or create new): [____]"
- **AND** system shows registered users with keyboard navigation: "← previous user | user-name | next user →"
- **AND** user can select existing user with arrow keys or type new name
- **AND** new users are created automatically with timestamp
- **AND** selected user authenticated for consolidation session

#### Scenario: Returning user selects from history
- **WHEN** consolidation starts and users have been registered
- **THEN** system displays recent users: "Last used: alice (2 days ago), bob (1 week ago)"
- **AND** user selects from list with arrow keys [↑↓] and confirms [Enter]
- **AND** or type new name to create different user account

### Requirement: Persistent Consolidation Flagging in Database
The system SHALL persist all consolidation flag decisions (keep/archive/delete) in the database with user and timestamp information, enabling pause/resume and multi-user review workflows.

#### Scenario: User flags files and exits without archiving
- **WHEN** user flags files (keep/archive/delete) during browsing
- **THEN** system stores each flag in database with:
  - File ID (object_meta.id)
  - Flag type (KEEP | ARCHIVE | DELETE)
  - User who flagged (username)
  - Timestamp when flagged
  - Flag reason (optional note: "corrupted", "best version", etc.)
- **AND** user can exit consolidation at any point: "[Q]uit and save flags"
- **AND** flags are safely persisted even if consolidation not executed
- **AND** user sees confirmation: "✓ Saved 47 flags for review later"

#### Scenario: User reviews previously flagged files
- **WHEN** user starts new consolidation session with same duplicate group
- **THEN** system shows previously flagged decisions:
```
Previous Flags (by alice, 2 days ago):
  ✓ KEEP: photo.jpg @ s3-backup
  ▶ ARCHIVE: old-photo.jpg @ archlinux
  ✗ DELETE: corrupted.jpg
```
- **AND** user can accept previous flags with one keystroke
- **AND** or modify individual flags
- **AND** all changes tracked to new user/timestamp

### Requirement: Multi-User Review and Collaboration
The system SHALL allow multiple users to review flagged decisions, confirming or modifying flags, and ensuring consensus before consolidation execution.

#### Scenario: Second user reviews flags from first user
- **WHEN** user logs in as different person: "alice flagged 47 files 2 days ago"
- **THEN** system displays: "Review pending flags by alice?"
- **AND** shows flag summary grouped by action:
```
Alice's Flagged Decisions (47 files):
  ✓ KEEP: 5 files (+ details)
  ▶ ARCHIVE: 23 files (+ destination)
  ✗ DELETE: 3 files (+ safety notes)

Review: [A]ccept All [M]odify Individual [R]eject All [C]ancel
```
- **AND** second user can [A]ccept alice's decisions unchanged
- **AND** or [M]odify specific flags with approval notation
- **AND** all changes logged with new user name and timestamp

#### Scenario: Modified flags show review trail
- **WHEN** second user modifies first user's flags
- **THEN** system tracks modification history:
```
Flag: DELETE photo.jpg
  Original: alice (FLAGGED as DELETE, 2 days ago)
  Modified: bob (CHANGED to KEEP, now)
  Reason: "Actually a good copy, bob verified quality"
```
- **AND** users can view complete review trail for each file
- **AND** final execution requires approval from last reviewer

#### Scenario: Consensus confirmation before execution
- **WHEN** multiple users have reviewed and modified flags
- **THEN** system shows final consolidated view with review trail:
```
📋 FINAL CONSOLIDATION SUMMARY (Multi-User Approved)

KEEP (5 files):
  ✓ photo.jpg — Flagged by alice, approved by bob

ARCHIVE (23 files):
  ▶ old-photo.jpg → archive-s3, Flagged by alice, bob verified

DELETE (3 files):
  ✗ corrupted.jpg — Flagged by alice, confirmed by bob
```
- **AND** execution requires confirmation from last reviewer
- **AND** system logs all reviewers and timestamps before execution
