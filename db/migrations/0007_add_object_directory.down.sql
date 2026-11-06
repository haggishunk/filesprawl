-- Rollback: Remove object_directory column and its index

BEGIN;

-- Drop the index
DROP INDEX IF EXISTS idx_object_meta_directory;

-- Drop the column
ALTER TABLE object_meta DROP COLUMN object_directory;

COMMIT;
