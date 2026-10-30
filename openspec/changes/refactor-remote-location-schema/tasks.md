# Implementation Tasks: Remote Location Schema Refactor

## Phase 1: Database Migration
- [ ] 1.1 Create migration adding `location` table with (id, name, type)
- [ ] 1.2 Add `location_id` column to `remote` table
- [ ] 1.3 Migrate existing data: insert locations from (hostname, remote_type) pairs
- [ ] 1.4 Add FK constraint from remote.location_id to location.id
- [ ] 1.5 Remove `hostname` column from `remote` table (verify no orphaned usage)
- [ ] 1.6 Create reverse migration for rollback

## Phase 2: Repository Layer Updates
- [ ] 2.1 Create `Location` type in internal/location package
- [ ] 2.2 Update `Remote` type to use location_id instead of hostname
- [ ] 2.3 Refactor ReadRemote/WriteRemote to work with Location
- [ ] 2.4 Update all queries (FindDuplicatesWithinRemote, etc.) to join location table
- [ ] 2.5 Add unique constraint check in ensureRemote for (location, remote_name, scan_root)

## Phase 3: CLI & Display Layer
- [ ] 3.1 Update report formatting to use location.name + remote_name (no hostname display)
- [ ] 3.2 Update duplicates command output formatting
- [ ] 3.3 Update --remote parsing to accept format: `location/remote_name` or just `remote_name`

## Phase 4: Tests
- [ ] 4.1 Update integration tests for location-based remote identification
- [ ] 4.2 Test false duplicate prevention (multi-root scans)
- [ ] 4.3 Test rclone remote without hostname
- [ ] 4.4 Verify migration forward and backward
