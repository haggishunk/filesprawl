-- Rollback Migration 0002: Restore remote schema to original state

BEGIN;

-- Restore hostname and remote_type columns from location table
ALTER TABLE remote ADD COLUMN hostname TEXT;
ALTER TABLE remote ADD COLUMN remote_type TEXT;

-- Populate hostname and remote_type from location table
UPDATE remote r
SET 
    hostname = l.name,
    remote_type = CASE WHEN l.type = 'local' THEN 'local' ELSE 'rclone' END
FROM location l
WHERE r.location_id = l.id;

-- Make hostname NOT NULL again
ALTER TABLE remote ALTER COLUMN hostname SET NOT NULL;
ALTER TABLE remote ALTER COLUMN remote_type SET NOT NULL;

-- Remove location_id column and its constraints
ALTER TABLE remote DROP COLUMN location_id;

-- Drop location table
DROP TABLE IF EXISTS location CASCADE;

COMMIT;
