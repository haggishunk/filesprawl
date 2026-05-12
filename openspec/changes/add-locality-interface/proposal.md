# Change: Add Locality Interface

## Why
rclone provides excellent remote storage configuration, but users need a way to map remote storage paths to local filesystem paths. This mapping enables users to understand where remote files would be located locally and vice versa. The configuration should be stored in a known location on the remote with optional local overrides.

## What Changes
- Add locality mapping configuration to map remote paths to local paths
- Support storing configuration in remote storage (e.g., `.filesprawl/locality.json`)
- Support local configuration overrides
- Provide API to resolve remote paths to local paths and vice versa
- Support per-remote and per-path mapping rules

## Impact
- Affected specs: New capability `locality-interface`
- Affected code:
  - New package `internal/locality` for mapping logic
  - New configuration types for locality mappings
  - Repository methods for storing/retrieving locality configurations
  - Integration with rclone to read/write configuration files on remotes

