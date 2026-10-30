-- CHECK DIRECTORY DELETION SAFETY
-- Usage: psql $DATABASE_URL -v dir_pattern="'%path/to/check/%'" -f check_directory_deletion_safety.sql
-- 
-- This script verifies that all files in a target directory exist elsewhere before deletion.
-- Returns three sections:
--   1. Summary: OK, RISKY, or SAFE_ONLY_ZERO_BYTES
--   2. Safe files (exist elsewhere or are 0-byte)
--   3. RISKY files (only copy exists in target directory)

\echo ''
\echo '╔════════════════════════════════════════════╗'
\echo '║   DIRECTORY DELETION SAFETY CHECKER       ║'
\echo '╚════════════════════════════════════════════╝'
\echo ''

-- Verify the pattern is set
\if :{?dir_pattern}
  \echo 'Checking directory pattern: :dir_pattern'
  \echo ''
\else
  \echo 'ERROR: Must set dir_pattern variable with:'
  \echo '  psql $DATABASE_URL -v "dir_pattern='"'"'%path/to/check/%'"'"'" -f check_directory_deletion_safety.sql'
  \q
\endif

WITH target_files AS (
  -- Find all files in the target directory
  SELECT
    om.id,
    om.object_name,
    om.object_path,
    om.object_size,
    oh.hash_value,
    oh.hash_type
  FROM object_meta om
  LEFT JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
  LEFT JOIN object_hash oh ON ohj.object_hash_id = oh.id
  WHERE om.object_path LIKE :dir_pattern
),
copies_elsewhere AS (
  -- Count non-target copies for each hash
  SELECT
    tf.id,
    tf.object_name,
    tf.object_path,
    tf.object_size,
    tf.hash_value,
    tf.hash_type,
    COUNT(DISTINCT om.id) as copies_elsewhere
  FROM target_files tf
  LEFT JOIN object_hash oh ON tf.hash_value = oh.hash_value AND tf.hash_type = oh.hash_type
  LEFT JOIN object_hash_junction ohj ON oh.id = ohj.object_hash_id
  LEFT JOIN object_meta om ON ohj.object_meta_id = om.id
    AND om.object_path NOT LIKE :dir_pattern
  GROUP BY tf.id, tf.object_name, tf.object_path, tf.object_size, tf.hash_value, tf.hash_type
),
safety_analysis AS (
  SELECT
    object_name,
    object_path,
    object_size,
    copies_elsewhere,
    CASE
      WHEN object_size = 0 THEN 'ZERO-BYTE (safe)'
      WHEN copies_elsewhere > 0 THEN 'Safe (exists elsewhere)'
      ELSE 'RISKY (only copy!)'
    END as status
  FROM copies_elsewhere
)
SELECT
  status,
  COUNT(*) as file_count,
  SUM(object_size) as total_bytes
FROM safety_analysis
GROUP BY status
ORDER BY
  CASE status
    WHEN 'RISKY (only copy!)' THEN 1
    WHEN 'Safe (exists elsewhere)' THEN 2
    ELSE 3
  END;
