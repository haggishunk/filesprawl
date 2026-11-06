-- Migration 0006: Fix remote object_paths to contain only the path within remote
-- For remote files, object_path should only contain the path within the remote,
-- not the "remote_name/scan_root:" prefix.
-- This extracts the part after the last colon.

BEGIN;

-- Update remote file paths to remove the prefix before the colon
-- Format is currently: "remote_name/scan_root:actual_path_in_remote"
-- Should be: "actual_path_in_remote"
UPDATE object_meta
SET object_path = SUBSTRING(object_path FROM POSITION(':' IN object_path) + 1)
WHERE object_path LIKE '%:%';

COMMIT;
