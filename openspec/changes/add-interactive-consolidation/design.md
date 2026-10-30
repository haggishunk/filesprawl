# Design: Interactive Duplicate Consolidation

## Context
Users have large numbers of duplicates across multiple remotes (e.g., 57k+ groups). The consolidation workflow must be:
1. **Interactive** — preview files, make informed decisions per-group or in batch
2. **Reversible** — archive instead of delete; preserve original paths for recovery
3. **Learnable** — remember user preferences to accelerate bulk consolidation
4. **Non-destructive** — never lose data; always maintain reference to originals
5. **Integrable** — leverage existing photo metadata (DigiKam) when available

## Goals
- Enable visual verification of duplicates before consolidation (critical for photo/media)
- Reduce manual per-file decisions via smart preference memory
- Support archive-first workflows with path preservation
- Optional photo metadata integration (DigiKam) for intelligent selection

## Non-Goals
- Permanent deletion (archiving is the final step, deletions are manual)
- Real-time sync across remotes
- Full media player integration (previews only; actual viewing external)

## Decisions
1. **Archive Strategy**: Copy to archive location preserving directory structure; delete source only after user confirmation
2. **Preferences**: Store consolidation decisions (per-remote, per-duplicate-pattern) in database; enable batch re-application
3. **Preview**: Lightweight image thumbnails (100x100px) embedded in ANSI terminal or external viewer option
4. **DigiKam**: Optional integration; detect photo libraries via EXIF/sidecar files; provide ratings/tags for smart selection
5. **Concurrency**: Single-threaded by default; optional parallel operations with progress tracking

## Risks
- **Large batches**: 118k+ groups could overwhelm terminal. Mitigation: paginate, batch mode, progress bars
- **Path collisions**: Archives may have path conflicts across remotes. Mitigation: prefix with remote name / hash
- **Rclone performance**: Moving files via RC API could be slow. Mitigation: batch operations, async support
