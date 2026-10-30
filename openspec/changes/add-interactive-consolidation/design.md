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
6. **Data Loss Prevention**: Integrate three-layer safety detection:
   - **Name Collision Detection**: Identify files with same name but different hashes; verify all variants preserved before consolidation
   - **Completeness Verification**: Detect 0-byte vs non-zero copies; ensure complete copies exist elsewhere before consolidation
   - **Directory Safety Check**: Verify all files in consolidation directory exist elsewhere; block if risky files detected
7. **Safety Display**: Present consolidated safety panel during interactive mode showing name collision status, completeness status, and directory deletion safety

## Risks
- **Large batches**: 118k+ groups could overwhelm terminal. Mitigation: paginate, batch mode, progress bars
- **Path collisions**: Archives may have path conflicts across remotes. Mitigation: prefix with remote name / hash
- **Rclone performance**: Moving files via RC API could be slow. Mitigation: batch operations, async support
- **Data Loss**: Multiple risk vectors during consolidation (filename collisions, incomplete transfers, missing variants). Mitigation: Three-layer safety detection with explicit warnings before consolidation

## Data Loss Prevention Architecture

### Layer 1: Name Collision Detection
- Query: Find files with identical object_name but different MD5 hashes
- Display: Show count of unique variants, status (all preserved / some missing / critical loss)
- Action: Flag for user review, recommend hash-based (not filename-based) consolidation
- Source: `db/scripts/detect_name_collisions.sql`

### Layer 2: Completeness Verification
- Query: For each file hash in duplicate group, count zero-byte vs non-zero copies
- Display: Show status by location (0-byte safe, non-zero verified, missing copies flagged)
- Action: Block consolidation if non-zero copy not found elsewhere
- Source: Hash-based queries in repository layer during consolidation analysis

### Layer 3: Directory Safety Check
- Query: For target consolidation directory, verify all files have copies elsewhere (by hash)
- Display: Categorize as safe (copies elsewhere), zero-byte (corrupted), or risky (only copy)
- Action: Block consolidation if risky files detected; require explicit override
- Source: `db/scripts/check_directory_deletion_safety.sql`

### Safety Panel Display
During interactive consolidation, present one-page safety summary:
```
📋 Group Summary: N copies across M remotes
🔗 Filename Status: ✅/⚠️/🚨 (X unique versions)
💾 Completeness: ✅/⚠️/🚨 (X zero-byte, Y complete)
📂 Directory Safety: ✅/⚠️/🚨 (X safe, Y risky)
```
User can expand each section for details before confirming consolidation.
