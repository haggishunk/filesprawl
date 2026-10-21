-- Migration 0001: symlink marking and scan provenance on object_meta.
-- Part of the fix-symlink-handling change. Reversible via
-- 0001_add_symlink_and_scan_root.down.sql.
BEGIN;

ALTER TABLE object_meta
    ADD COLUMN IF NOT EXISTS is_symlink BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE object_meta
    ADD COLUMN IF NOT EXISTS scan_root TEXT;

-- Best-effort backfill of scan provenance from the linked remote record.
-- For local scans the remote name is the scanned root path; for rclone
-- remotes it is the remote name (the indexed sub-path was not recorded).
UPDATE object_meta om
SET scan_root = r.remote_name
FROM object_remote_junction orj
JOIN remote r ON r.id = orj.remote_id
WHERE orj.object_meta_id = om.id
  AND om.scan_root IS NULL;

CREATE INDEX IF NOT EXISTS idx_object_meta_scan_root
    ON object_meta (scan_root);

COMMIT;
