# Change: Add Duplicate Detection

## Why
The system currently indexes files and their hashes but doesn't provide functionality to identify duplicate files. Users need to find duplicates both within a single remote storage location and across multiple remotes to optimize storage usage and understand data redundancy.

## What Changes
- Add duplicate detection capability to identify files with identical content hashes
- Support finding duplicates within a single remote storage location
- Support finding duplicates across multiple remote storage locations
- Provide query interface for retrieving duplicate file groups
- Include metadata about where duplicates are located (remote, path)

## Impact
- Affected specs: New capability `duplicate-detection`
- Affected code: 
  - New package `internal/analysis` for duplicate detection logic
  - New repository methods in `internal/repository` for duplicate queries
  - Potential CLI commands or API endpoints to expose duplicate detection

