# Change: Add Interactive Duplicate Consolidation Workflow

## Why
Users currently can identify duplicates across remotes but lack an interactive workflow to consolidate them. This blocks the primary use case: discovering duplicates, previewing them (especially images), deciding which to keep, and consolidating to a chosen location—all while preserving options for later cleanup. Users need memory of consolidation patterns to speed up batch decisions.

## What Changes
- New `consolidate` command enabling interactive discovery, preview, and consolidation of duplicates
- Image preview support (thumbnails, embedded in terminal output) for visual verification
- Consolidation memory: persist "consolidation preferences" (e.g., "keep newest", "archive to S3") to speed up batch workflows
- Archive-first approach: consolidations move files to an archive location (preserving original directory structure) rather than permanently deleting
- Optional DigiKam database integration: detect photo collections and use existing tags/ratings for smart consolidation
- Per-remote consolidation targets: define default consolidation destinations per remote for one-click decisions
- **Multi-user flagging and review workflow** (NEW):
  - Persistent consolidation flagging: users flag files (keep/archive/delete) and decisions are stored in database
  - Session independence: users can flag files, exit, and resume later without losing progress
  - User authentication: simple username-based system to track who flags files
  - Multi-user collaboration: second user can review and modify first user's flags with full audit trail
  - Consensus workflow: multiple reviewers approve decisions with complete modification history before execution
- **Data loss prevention integration** (NEW):
  - Name collision detection: identify files with same name but different content; verify all variants preserved
  - Completeness verification: detect 0-byte vs non-zero copies; ensure complete copies exist elsewhere
  - Directory deletion safety: verify all files in deletion target exist elsewhere before allowing delete
- **Dual-view browsing** (NEW):
  - Hash-centric view: group identical files by content across all locations
  - Filename-centric view: group different versions of same filename to choose best version
  - Location presence indicators: 🟢/🟡/🔴 color-coded backup replication health

## Impact
- Affected specs: `duplicate-detection`, `object-repository` (new write operations), `consolidation-workflow` (new capability), `consolidation-browsing` (new capability)
- Affected code:
  - `cmd/filesprawl/consolidate.go` (new command)
  - `internal/consolidation/` (new package for workflow logic)
  - `internal/consolidation/users/` (new subpackage for user management)
  - `internal/consolidation/flags/` (new subpackage for persistent flagging)
  - `internal/consolidation/browsing/` (new subpackage for dual-view interface)
  - `internal/safety/` (new package for data loss prevention)
  - `internal/archive/` (new package for archiving with path preservation)
  - `internal/previews/` (new package for image thumbnail generation)
  - `internal/digikam/` (new optional integration package)
  - Database schema additions:
    - `consolidation_users` (id, username, created_at, last_active)
    - `consolidation_flags` (id, file_id, user_id, flag_type, flagged_at, reason, archive_dest)
    - `consolidation_reviews` (id, file_id, original_flag, original_user, original_at, reviewed_by, reviewed_at, new_flag, review_reason)
    - Database migrations for all three tables
- **BREAKING**: No breaking changes; purely additive
