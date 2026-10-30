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
8. **Dual-View Browsing**: Hash-centric and filename-centric views for understanding duplicate relationships
   - **Hash View**: Group by content (hash), show all variants across locations
   - **Filename View**: Group by name, show all different versions (hashes)
   - **Decision Flagging**: Individual file marking for keep/archive/delete with persistent state
9. **Location Presence Indicators**: Color-coded status (🟢 green/🟡 yellow/🔴 red) showing backup replication health and single-point-of-failure warnings

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

## Dual-View Browsing Interface Architecture

### Hash-Centric View (Browse by Content)
```
View: By Hash (47 total copies)
═══════════════════════════════════════════

Hash abc123... (47 copies)  [Expand ▼]
├─ 🟢 archlinux//goblet/light (local)      — 3 copies
│  ├─ photo.jpg @ /2020/03/10/
│  ├─ photo.jpg @ /2020/04/15/
│  └─ photo.jpg @ /2019/misc/
├─ 🟢 archlinux/do (dropbox)               — 2 copies
│  ├─ photo.jpg @ /photos/2020/
│  └─ photo.jpg @ /old-photos/
└─ 🟢 s3-backup/archive (S3)               — 1 copy
   └─ photo.jpg @ /archive/recovered/

[k]eep  [a]rchive  [d]elete  [TAB]switch view
```

**Advantage**: See complete duplication picture for one piece of content
**Use Case**: "Is this file backed up well enough to delete one copy?"

### Filename-Centric View (Browse by Name)
```
View: By Filename (118k files with variants)
═══════════════════════════════════════════

photo.jpg (5 different versions)  [Expand ▼]
├─ Version 1 (hash abc123... — 47 copies)
│  ├─ 🟢 archlinux//goblet/light — 3
│  ├─ 🟢 archlinux/do — 2
│  └─ 🟢 s3-backup — 1
├─ Version 2 (hash xyz789... — 12 copies)
│  ├─ 🟢 archlinux//goblet/light — 2
│  └─ 🟡 old-archive — 1
├─ Version 3 (hash 000fff... — 0 bytes — 8 copies)
│  ├─ 🔴 archlinux//goblet — ONLY copy
│  └─ 🔴 google-photos — ONLY copy
...

[k]eep  [a]rchive  [d]elete  [TAB]switch view
```

**Advantage**: Understand all versions of "the same" file at once
**Use Case**: "Which version of photo.jpg should I keep?"

### Location Presence Indicators
- 🟢 **Green**: Multiple copies across locations (safe to delete one)
- 🟡 **Yellow**: Limited replication (only 1-2 locations, recommend backup)
- 🔴 **Red**: Only copy in entire system (CRITICAL - no backup exists!)

### Decision State Machine per File
```
┌─────────────────────────────────────┐
│          Initial State               │
│  (no flag/decision made)             │
└────────────┬────────────────────────┘
             │
    ┌────────┼────────┬────────┐
    │        │        │        │
    ▼        ▼        ▼        ▼
 [KEEP]   [ARCH]   [DEL]   [SKIP]
    │        │        │        │
    └────────┼────────┼────────┘
             │
    ┌────────▼────────┐
    │ Summary Screen  │
    │ for Execution   │
    └─────────────────┘
```

User can:
- Flag individual files with keyboard shortcuts
- Toggle flags while browsing
- View summary grouped by action (KEEP/ARCHIVE/DELETE)
- Modify decisions in summary view
- Execute consolidation only on flagged files

### Data Model for Browsing State
```
type ConsolidationSession {
  DuplicateGroup duplicates
  Map[FileID] decision (KEEP | ARCHIVE | DELETE | SKIP)
  CurrentView (BY_HASH | BY_FILENAME)
  ExpandedItems Set[FileID | HashID]
  PageCursor Position
}
```
