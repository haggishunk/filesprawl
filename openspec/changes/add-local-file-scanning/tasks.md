# Implementation Tasks

## 1. Local Scan Command Surface
- [ ] 1.1 Define a CLI command for indexing a local filesystem root
- [ ] 1.2 Validate local path inputs and return clear usage errors
- [ ] 1.3 Document example local indexing commands

## 2. Local Traversal and Hashing
- [ ] 2.1 Add local filesystem traversal that walks files under a selected root
- [ ] 2.2 Compute supported content hashes for local files during scanning
- [ ] 2.3 Skip or report unreadable files with clear error context
- [ ] 2.4 Add unit tests for local traversal and hashing behavior

## 3. Persistence and Source Identity
- [ ] 3.1 Persist local scan results into the existing object metadata and hash tables
- [ ] 3.2 Reuse the source identity model with a `local` source type and host-qualified local origin identity
- [ ] 3.3 Persist object-to-source links so duplicate queries can compare local and remote results together
- [ ] 3.4 Add repository and integration tests for mixed local/remote duplicate lookups

## 4. Validation
- [ ] 4.1 Run `go test ./...`
- [ ] 4.2 Run `openspec validate add-local-file-scanning --strict`