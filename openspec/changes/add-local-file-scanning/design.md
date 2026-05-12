## Context
Duplicate detection becomes much more useful when Filesprawl can compare remote-backed files with files that exist on local disks. The project already has a persistence model for file metadata, hashes, and host-qualified source identities. The local scanning change should extend that model without depending on rclone or locality mapping.

## Goals / Non-Goals
- Goals:
  - Index local filesystem files directly and compute their hashes
  - Persist local scan results into the same hash index used for remotes
  - Make local and remote files comparable in duplicate detection workflows
- Non-Goals:
  - Replace locality mapping or infer local paths for remote objects
  - Require local paths to be configured as rclone remotes
  - Add file modification watching or background sync

## Decisions
- Decision: Represent local scan origins using the existing source identity table with `remote_type = local`.
  - Rationale: this keeps duplicate queries and provenance handling unified across local and remote scans.
- Decision: Use the host-qualified local scan root as the persisted local source identity.
  - Rationale: a normalized local root path is deterministic and does not require a separate alias-management system.
- Decision: Local file hashing will be performed by Filesprawl directly rather than delegated to rclone RC.
  - Rationale: local scanning should not depend on an rclone remote definition to participate in duplicate indexing.

## Risks / Trade-offs
- Local scans may be slower than remote metadata scans because hashes must be computed locally.
  - Mitigation: start with a simple synchronous hashing pipeline and optimize later if needed.
- A local root path may be less user-friendly than a named alias.
  - Mitigation: use the normalized path now and leave alias support for a future change.

## Migration Plan
1. Add the local scan CLI entrypoint.
2. Add local traversal and hashing.
3. Persist local source identities and object links using the existing tables.
4. Add tests for mixed local/remote duplicate lookup behavior.

## Open Questions
- Should a future change allow choosing which hash algorithms to compute for local files?
- Should a future change support exclusion patterns for local scans?