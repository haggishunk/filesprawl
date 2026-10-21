-- Rollback for migration 0001: remove symlink marking and scan provenance
-- from object_meta.
BEGIN;

DROP INDEX IF EXISTS idx_object_meta_scan_root;

ALTER TABLE object_meta
    DROP COLUMN IF EXISTS scan_root;

ALTER TABLE object_meta
    DROP COLUMN IF EXISTS is_symlink;

COMMIT;
