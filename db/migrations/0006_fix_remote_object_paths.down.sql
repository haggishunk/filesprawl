-- Rollback: Restore remote object_paths with prefix
-- This migration cannot be perfectly reversed as we no longer have the
-- remote_name/scan_root information. This is a data-destructive operation.
-- In practice, you should not roll back this migration in production.

BEGIN;

-- Note: This down migration is provided for schema completeness only.
-- The actual restoration would require the original remote metadata.
-- For now, we'll just leave paths as-is (without prefix) if rolled back.

ROLLBACK;
