# Change: Refactor Remote & Location Schema to Fix Duplicate False Positives

## Why
Current schema conflates location (hostname/rclone remote) with scan scope. Two problems:
1. **Local scan false duplicates**: Scanning same filesystem from different roots (e.g., `/media` vs `/media/2024`) creates duplicate metadata + false duplicate groups
2. **Hostname on remote-only sources**: `archlinux/do` shouldn't have hostname; hostname is irrelevant for rclone remotes

## What Changes
- **NEW: `location` table** — Uniquely identifies scan source (hostname for local, remote name for rclone)
- **MODIFIED: `remote` table** — Replace `hostname` column with `location_id` FK to location table
- **CLARIFIED: `scan_root` semantics** — Only applies to local scans; null/empty for remote scans
- **NEW: Unique constraint** — `(location_id, remote_name, scan_root)` to prevent false duplicates

## Impact
- Affected specs: `remote-configuration`, `database-persistence`
- Affected code: repository layer (ReadRemote, WriteRemote, queries)
- Database migration required (reversible)
- **BREAKING**: API unchanged; internal representation changes

## Benefits
- Eliminates false duplicates from multi-root local scans
- Correctly distinguishes remotes from local scans
- Cleaner semantics: location = source, remote/scan_root = path within source
