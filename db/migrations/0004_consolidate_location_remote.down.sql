-- Rollback Migration 0004: Restore location and remote table structure
-- Splits location back to location + remote
-- Converts absolute_path back to relative paths

BEGIN;

-- Step 1: Recreate remote table
CREATE TABLE remote (
  id SERIAL PRIMARY KEY,
  location_id INT NOT NULL REFERENCES location(id) ON DELETE RESTRICT,
  remote_name TEXT NOT NULL,
  scan_root TEXT
);

-- Step 2: Recreate object_remote_junction
CREATE TABLE object_remote_junction (
  id SERIAL PRIMARY KEY,
  object_meta_id INT,
  remote_id INT,
  FOREIGN KEY (object_meta_id) REFERENCES object_meta(id),
  FOREIGN KEY (remote_id) REFERENCES remote(id)
);

-- Step 3: Re-add scan_root column to object_meta
ALTER TABLE object_meta ADD COLUMN scan_root TEXT;

-- Step 4: Backfill remote table from location + object_meta data
-- Extract remote_name and scan_root from absolute_path
INSERT INTO remote (location_id, remote_name, scan_root)
SELECT DISTINCT 
  om.location_id,
  CASE 
    WHEN om.scan_root IS NOT NULL THEN om.scan_root
    ELSE SUBSTRING(om.absolute_path FROM 1 FOR POSITION(':' IN om.absolute_path) - 1)
  END AS remote_name,
  om.scan_root
FROM object_meta om;

-- Step 5: Backfill object_remote_junction
INSERT INTO object_remote_junction (object_meta_id, remote_id)
SELECT 
  om.id,
  r.id
FROM object_meta om
JOIN remote r ON om.location_id = r.location_id
WHERE (om.scan_root IS NOT NULL AND r.scan_root = om.scan_root)
   OR (om.scan_root IS NULL AND 
       r.remote_name = SUBSTRING(om.absolute_path FROM 1 FOR POSITION(':' IN om.absolute_path) - 1));

-- Step 6: Drop constraints and columns
ALTER TABLE object_meta 
  DROP CONSTRAINT fk_object_meta_location,
  DROP COLUMN location_id,
  DROP COLUMN absolute_path;

-- Step 7: Drop new index
DROP INDEX IF EXISTS idx_object_meta_location_path;

-- Step 8: Re-add old index
CREATE INDEX idx_object_meta_scan_root ON object_meta (scan_root);

COMMIT;
