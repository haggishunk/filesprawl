-- Migration 0005: Drop absolute_path column and consolidate into object_path
-- Converges on object_path as the single source of truth for file paths:
-- - For local files: object_path contains the absolute filesystem path
-- - For remote files: object_path contains the path within the remote (as before)
-- This eliminates redundancy and makes the schema simpler.

BEGIN;

-- Drop the index that used absolute_path
DROP INDEX IF EXISTS idx_object_meta_location_path;

-- Update object_path to contain absolute paths for records that had absolute_path set
-- For rows where absolute_path differs from object_path, use absolute_path
UPDATE object_meta
SET object_path = absolute_path
WHERE absolute_path IS NOT NULL AND absolute_path != object_path;

-- Drop the absolute_path column
ALTER TABLE object_meta DROP COLUMN absolute_path;

-- Recreate the index on location_id and object_path
CREATE INDEX idx_object_meta_location_path 
  ON object_meta(location_id, object_path);

COMMIT;
