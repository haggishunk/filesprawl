-- Migration 0004: Consolidate location and remote tables
-- Merges remote table into location table, stores absolute paths in object_meta
-- Simplifies schema: location now has name + type + absolute_path
-- object_meta gains location_id, absolute_path; loses scan_root column

BEGIN;

-- Step 1: Add new columns to object_meta
-- location_id: FK to location (identifies source of file)
-- absolute_path: full path to file (replaces relative path model)
ALTER TABLE object_meta
  ADD COLUMN location_id INT,
  ADD COLUMN absolute_path TEXT;

-- Step 2: Backfill object_meta with data from remote table
-- For each object, find its remote, get location, build absolute path
UPDATE object_meta om
SET 
  location_id = r.location_id,
  absolute_path = CASE 
    WHEN om.scan_root IS NOT NULL THEN 
      -- Local file: absolute_path = scan_root + '/' + object_path
      om.scan_root || '/' || om.object_path
    ELSE
      -- Remote file: absolute_path = remote_name + ':' + object_path
      r.remote_name || ':' || om.object_path
  END
FROM object_remote_junction orj
JOIN remote r ON orj.remote_id = r.id
WHERE om.id = orj.object_meta_id;

-- Step 3: Add NOT NULL constraints and foreign key
ALTER TABLE object_meta
  ALTER COLUMN location_id SET NOT NULL,
  ALTER COLUMN absolute_path SET NOT NULL,
  ADD CONSTRAINT fk_object_meta_location 
    FOREIGN KEY (location_id) REFERENCES location(id) ON DELETE RESTRICT;

-- Step 4: Create index on (location_id, absolute_path) for efficient lookups
CREATE INDEX idx_object_meta_location_path 
  ON object_meta(location_id, absolute_path);

-- Step 5: Drop old junction and remote table
DROP TABLE object_remote_junction;
DROP TABLE remote;

-- Step 6: Drop scan_root column (absolute path now in absolute_path)
ALTER TABLE object_meta DROP COLUMN scan_root;

-- Step 7: Drop old scan_root index
DROP INDEX IF EXISTS idx_object_meta_scan_root;

COMMIT;
