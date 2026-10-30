# Design: Remote Location Schema Refactor

## Problem Statement
1. **Multi-root scan duplicates**: Same file in `/photos` and `/photos/2024` → two metadata entries → false duplicate group
2. **Hostname semantic mismatch**: Rclone remote `do` shouldn't have hostname "archlinux"
3. **Confusion in reporting**: Display shows `archlinux//goblet/light root: /goblet/light` which mixes location with path

## Solution
**Separate `location` (where data comes from) from `remote_name` + `scan_root` (path within source)**

### New Schema
```sql
location
  id: int PK
  name: string (hostname for local, remote identifier for rclone)
  type: enum ('local', 'remote')
  UNIQUE(name, type)

remote
  id: int PK
  location_id: int FK → location
  remote_name: string (path or rclone name)
  remote_type: string ('local', 'rclone')
  UNIQUE(location_id, remote_name, scan_root)

object_meta
  scan_root: string (local directory root; NULL for rclone)
```

## Key Decisions
1. **Location identifies source** — hostname for local, remote name for rclone
2. **Type prevents confusion** — 'local' vs 'remote' makes semantics explicit
3. **Scan root only for local** — NULL for remote indexing
4. **Unique constraint** — (location, remote_name, scan_root) prevents false duplicates

## Migration Strategy
- Create location table, migrate data, add FK, remove hostname column
- Backward-compatible: API layer converts back to `hostname/remote_name` if needed
