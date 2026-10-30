-- Migration 0002: Add location table and refactor remote schema
-- Separates location (source of data) from remote configuration
-- Prevents false duplicates from multi-root scans and fixes hostname semantics
-- Reversible via 0002_add_location_table.down.sql

BEGIN;

-- Create location table: identifies the source of indexed data
-- For local scans: name=hostname, type='local'
-- For remotes: name=remote identifier, type='remote'
CREATE TABLE IF NOT EXISTS location (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('local', 'remote')),
    UNIQUE(name, type)
);

-- Migrate existing remote records to location table
-- This creates location records from unique (hostname, remote_type) pairs
INSERT INTO location (name, type)
SELECT DISTINCT hostname, CASE 
    WHEN remote_type = 'local' THEN 'local'
    ELSE 'remote'
END as type
FROM remote
ON CONFLICT (name, type) DO NOTHING;

-- Add location_id column to remote table
ALTER TABLE remote ADD COLUMN location_id INT REFERENCES location(id) ON DELETE RESTRICT;

-- Populate location_id by joining with the newly created location records
UPDATE remote r
SET location_id = (
    SELECT l.id FROM location l 
    WHERE l.name = r.hostname 
    AND l.type = CASE WHEN r.remote_type = 'local' THEN 'local' ELSE 'remote' END
);

-- Ensure all location_id values are set (fail if not)
ALTER TABLE remote ALTER COLUMN location_id SET NOT NULL;

-- Remove hostname column (it's now in location table)
ALTER TABLE remote DROP COLUMN hostname;

-- Drop old remote_type column (type is now in location table)
ALTER TABLE remote DROP COLUMN remote_type;

COMMIT;
