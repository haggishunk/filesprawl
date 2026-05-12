## Context
Filesprawl now persists host-qualified remote identities and can query duplicate hashes. However, the current cross-remote query counts duplicate files without requiring more than one remote identity. This causes same-remote duplicates to appear in a report that users interpret as a cross-remote comparison.

## Goals / Non-Goals
- Goals:
  - Make cross-remote duplicate results require at least two distinct persisted remote identities
  - Preserve current within-remote duplicate behavior
  - Keep the semantics aligned with host-qualified remote identity persistence
- Non-Goals:
  - Add a new "across all indexed files" query in the same change
  - Change the persistence model for remotes or hashes

## Decisions
- Decision: Cross-remote duplicate detection will key distinctness off persisted `remote.id` values.
  - Rationale: persisted remote IDs already capture hostname + remote name identity, so they are the correct comparison unit.
- Decision: Files without remote associations will not qualify a hash group for cross-remote duplicate results.
  - Rationale: a result cannot be proven cross-remote without remote identity information.

## Risks / Trade-offs
- Existing users may notice fewer "across remotes" results than before.
  - Mitigation: the updated semantics more accurately match the report name and user expectation.

## Migration Plan
1. Update the repository query.
2. Add tests showing same-remote duplicates are excluded.
3. Document the stricter semantics in the proposal/spec trail.