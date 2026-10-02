## Context
Local file scanning needs to handle symlinks correctly while tracking scan provenance. Current implementation crashes on directory symlinks, falsely matches symlink targets, and creates collisions when same-named files exist in different scan roots.

## Goals / Non-Goals
- Goals:
  - Skip directory symlinks gracefully (no crash)
  - Mark file symlinks in metadata
  - Track scan root for each file
  - Correctly detect duplicates across multiple scan roots
  - Maintain backward compatibility for non-symlink files
- Non-Goals:
  - Follow symlinks to create linked copies
  - Resolve symlink chains transitively
  - Store symlink target paths

## Decisions

**Decision: Use os.Lstat() for symlink detection**
- Rationale: Doesn't follow symlinks, safe for detection
- Alternative: Walk with follow logic (more complex, harder to control)

**Decision: Store is_symlink in object_meta**
- Rationale: Preserves metadata as single source of truth
- Alternative: Separate table (adds complexity)

**Decision: Skip hashing symlinks entirely**
- Rationale: Symlinks shouldn't match their targets
- Alternative: Hash target (confuses duplicate detection)

**Decision: Track scan_root in object_meta**
- Rationale: Fixes path collision issue, enables accurate scan provenance
- Alternative: Store as separate metadata field (would complicate queries)

## Risks / Trade-offs
- Database migration required (can be scripted for existing data)
- Queries become slightly more complex with scan_root filtering
- Backward compatibility: Old records without scan_root need migration

## Migration Plan
1. Add schema columns (is_symlink, scan_root)
2. Update local_scan.go for symlink handling
3. Update repository methods
4. Update duplicate detection logic
5. Migrate existing data (set scan_root to remote name + path for remote files)
6. Test with symlink-heavy directories

## Open Questions
- Should remote scans also track scan_root? (Recommendation: Yes, for consistency)
- How to handle broken symlinks? (Skip with warning)
