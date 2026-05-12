# Change: Update Cross-Remote Duplicate Semantics

## Why
The current duplicate-query behavior groups files by shared hashes across the entire indexed corpus. That is useful, but it does not match the stronger meaning of "across remotes" that users expect. A hash should only be reported as a cross-remote duplicate when the matching files are linked to at least two distinct persisted remote identities.

## What Changes
- Tighten the semantics of "duplicates across remotes" so results require at least two distinct remote identities
- Preserve existing within-remote duplicate behavior
- Ensure remote identity distinctness is based on the persisted remote records rather than just file count
- Document the difference between "across all indexed files" and "across remotes"

## Impact
- Affected specs:
  - Modified capability `object-repository`
- Affected code:
  - `internal/repository` duplicate query logic
  - `internal/analysis` reporting behavior that depends on cross-remote duplicate groups