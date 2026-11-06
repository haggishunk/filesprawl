-- Rollback Migration 0005: Restore absolute_path column
-- Splits object_path back into object_path (relative) and absolute_path (full path)

BEGIN;

-- Drop the index that uses object_path
DROP INDEX IF EXISTS idx_object_meta_location_path;

-- Re-add the absolute_path column
ALTER TABLE object_meta ADD COLUMN absolute_path TEXT;

-- Backfill absolute_path: use object_path (they were unified during migration)
UPDATE object_meta
SET absolute_path = object_path;

-- Make absolute_path NOT NULL
ALTER TABLE object_meta ALTER COLUMN absolute_path SET NOT NULL;

-- Recreate the index on location_id and absolute_path
CREATE INDEX idx_object_meta_location_path 
  ON object_meta(location_id, absolute_path);

COMMIT;
