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

## Impact
- Affected specs: `duplicate-detection`, `object-repository` (new write operations), `consolidation-workflow` (new capability)
- Affected code: 
  - `cmd/filesprawl/consolidate.go` (new command)
  - `internal/consolidation/` (new package for workflow logic)
  - `internal/archive/` (new package for archiving with path preservation)
  - `internal/previews/` (new package for image thumbnail generation)
  - `internal/digikam/` (new optional integration package)
  - Database schema additions (consolidation history, preferences)
- **BREAKING**: No breaking changes; purely additive
