# Implementation Tasks

## 1. Repository Query Semantics
- [ ] 1.1 Update the cross-remote duplicate query to require at least two distinct remote identities per hash group
- [ ] 1.2 Keep within-remote duplicate queries unchanged
- [ ] 1.3 Add repository coverage for same-remote-only duplicates versus distinct-remote duplicates

## 2. Analysis / Reporting Semantics
- [ ] 2.1 Ensure duplicate reporting continues to surface hostname and remote name for cross-remote groups
- [ ] 2.2 Add or update tests proving same-remote duplicates are excluded from cross-remote results

## 3. Validation
- [ ] 3.1 Run `go test ./...`
- [ ] 3.2 Run `openspec validate update-cross-remote-duplicate-semantics --strict`