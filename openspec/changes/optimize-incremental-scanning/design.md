# Design: Incremental Scanning with Delta Detection

## Context

Current: Walk entire filesystem, hash all files, persist all, every time.  
Target: Skip unchanged files, track deletions, enable weekly/daily scans for 1M+ libraries.  
Constraints: Backward compatible, conservative (fallback to full scan), minimal schema changes.

## Goals

✅ Reduce scan time 95% for typical incremental updates (1% file changes)  
✅ Enable weekly scans for 1M photo libraries (27h → 18min)  
✅ Track file lifecycle (creation, modification, deletion times)  
✅ Detect moved files (same hash, different path)  
✅ Backward compatible (optional flag, old DBs work)

**Non-Goals:**

- Remove deleted files (keep history)
- Change hash algorithm
- Parallel scanning in Phase 1 (Phase 4)
- Cloud-synced scan state

## Decisions

### 1. Classification Approach

**Decision:** File state classified by (size, mtime) match, not full hash

- **Why:** Size+mtime check is ~10000x faster than hashing
- **Alternative:** Trust filesystem timestamps (risky with clock skew)
- **Mitigation:** Periodic full verification every 10th scan

### 2. Deletion Tracking

**Decision:** Mark deleted, don't purge (is_deleted=TRUE, deletion_detected_at=NOW)

- **Why:** Enables archive cleanup decisions, lifecycle analysis
- **Alternative:** Delete immediately (loses historical data)
- **Trade-off:** DB grows, but controlled by retention policy

### 3. Change Detection Fields

**Decision:** Store (last_scanned_at, last_modtime, is_deleted, deletion_detected_at)

- **Why:** Minimal schema (4 fields), sufficient for delta detection
- **Alternative:** Store full directory tree snapshot (complex, large)
- **Trade-off:** Can't skip entire subtrees, but simpler implementation

### 4. API Design

**Decision:** --delta flag (opt-in), not automatic behavior change

- **Why:** Conservative rollout, backward compatible, explicit user intent
- **Alternative:** Auto-enable delta (riskier, harder to debug)
- **Process:** Start disabled, enable after 2 weeks validation

### 5. Hash Caching

**Decision:** Trust size+mtime, skip hash if unchanged

- **Why:** 99% correctness for typical incremental scans
- **Alternative:** Always verify hash (safe but slow)
- **Fallback:** --force flag to skip delta, rescan all

## Risks & Mitigations

| Risk                          | Probability | Impact            | Mitigation                           |
| ----------------------------- | ----------- | ----------------- | ------------------------------------ |
| Clock skew (mtime unreliable) | Medium      | Wrong "unchanged" | Periodic full verify, --force option |
| Rclone missing mtime          | High        | Can't use delta   | Detect mtime, fallback to hash       |
| Rapid file update (2x mtime)  | Low         | File missed       | Periodic full verification           |
| Deleted files pile up         | Medium      | DB bloat          | Retention policy, is_deleted filter  |

## Migration Plan

**Phase 1:** Deploy schema, repo methods

- New fields nullable: existing queries work
- Old DB works unchanged
- No data migration needed

**Phase 2:** Deploy delta scanning (--delta disabled)

- Code change only
- No runtime behavior change
- Backward compatible

**Phase 3:** Enable --delta flag

- Manual opt-in for testing
- Monitor metrics (unchanged/new/changed/deleted counts)
- Collect performance data

**Phase 4:** Default to enabled after validation

- If issues: disable with --no-delta
- Keep --force option for full rescan

## Open Questions

1. **Periodic verification interval?** Current: every 10th scan. Configurable?
2. **Deletion retention window?** Current: keep forever. Should old deletions be purged?
3. **Rclone mtime support?** Need to validate availability across remote types.
4. **Directory snapshot optimization?** Phase 4 feature, defer decision until Phase 2 complete.
