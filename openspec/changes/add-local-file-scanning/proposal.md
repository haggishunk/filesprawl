# Change: Add Local File Scanning

## Why
Filesprawl can index remote-backed files and compare their hashes, but users also need to scan local filesystem content so duplicates can be found across remotes and local storage. Local files need to participate in the same hash index and duplicate workflows without requiring an rclone remote definition.

## What Changes
- Add a local file scanning capability that walks local filesystem roots and hashes files directly
- Persist local scan results into the existing metadata/hash index so duplicate lookups can compare local and remote files together
- Reuse the persisted source identity model by storing local scan origins as host-qualified source records with a local source type
- Add a CLI entrypoint for indexing a local path
- Keep the locality mapping proposal separate; this change is about indexing local content, not resolving remote/local path mappings

## Impact
- Affected specs:
  - New capability `local-file-scanning`
  - Modified capability `object-repository`
  - Modified capability `remote-configuration`
- Affected code:
  - `main.go` CLI command dispatch
  - `internal/operation` for local scan traversal
  - `internal/repository` for persistence of local scan origins and hashes
  - new or extended hashing utilities for local files
  - tests and documentation for mixed local/remote duplicate workflows