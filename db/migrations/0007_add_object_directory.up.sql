-- Migration 0007: Add object_directory column for incremental scanning
-- Stores the immediate parent directory of each file.
-- Enables --only-missing-dirs feature to skip re-scanning indexed directories.

BEGIN;

-- Add object_directory column
ALTER TABLE object_meta ADD COLUMN object_directory TEXT;

-- Populate object_directory from existing object_path data
-- For all files: extract parent directory (everything before the last /)
-- Examples:
--   /path/to/file.txt → /path/to
--   photos/subdir/file.jpg → photos/subdir
--   file.txt → (empty string for files in root)
UPDATE object_meta
SET object_directory = CASE
  WHEN object_path LIKE '%/%' THEN
    -- Has at least one slash: extract everything before the last slash
    SUBSTR(object_path, 1, LENGTH(object_path) - POSITION('/' IN REVERSE(object_path)))
  ELSE
    -- No slash: file is in root directory
    ''
END;

-- Make column NOT NULL now that it's populated
ALTER TABLE object_meta ALTER COLUMN object_directory SET NOT NULL;

-- Create index for efficient directory queries
CREATE INDEX IF NOT EXISTS idx_object_meta_directory 
  ON object_meta(location_id, object_directory);

COMMIT;
