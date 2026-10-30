-- DETECT NAME COLLISIONS (different files, same filename)
-- Identifies files with identical names but different content hashes.
-- This is a data loss risk if you delete by filename!
--
-- Usage: psql $DATABASE_URL -f db/scripts/detect_name_collisions.sql
--
-- Results include:
-- - CRITICAL: Some variants missing (data loss risk!)
-- - WARNING: All variants present but only 1 unique version exists
-- - OK: All different variants are preserved

\echo ''
\echo '╔════════════════════════════════════════════╗'
\echo '║     FILENAME COLLISION DETECTION           ║'
\echo '║  (Same name, Different content)            ║'
\echo '╚════════════════════════════════════════════╝'
\echo ''

WITH filename_hashes AS (
  SELECT 
    om.object_name,
    oh.hash_value,
    oh.hash_type,
    COUNT(DISTINCT om.id) as total_copies,
    COUNT(DISTINCT CASE WHEN om.object_size > 0 THEN om.id END) as nonzero_copies
  FROM object_meta om
  LEFT JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
  LEFT JOIN object_hash oh ON ohj.object_hash_id = oh.id
  WHERE oh.hash_type = 'md5' AND om.object_size > 0
  GROUP BY om.object_name, oh.hash_value, oh.hash_type
),
collision_summary AS (
  SELECT 
    object_name,
    COUNT(DISTINCT hash_value) as distinct_versions,
    SUM(total_copies) as total_instances,
    SUM(CASE WHEN nonzero_copies > 0 THEN 1 ELSE 0 END) as nonzero_versions
  FROM filename_hashes
  GROUP BY object_name
  HAVING COUNT(DISTINCT hash_value) > 1
)
SELECT 
  object_name,
  distinct_versions as unique_variants,
  nonzero_versions as preserved_variants,
  total_instances as total_copies,
  CASE 
    WHEN nonzero_versions < distinct_versions 
      THEN '🚨 CRITICAL: ' || (distinct_versions - nonzero_versions) || ' variant(s) MISSING!'
    WHEN nonzero_versions = 1 
      THEN '⚠️  WARNING: Only 1 version preserved, ' || (distinct_versions-1) || ' lost'
    ELSE '✅ OK: All ' || distinct_versions || ' variants preserved'
  END as status
FROM collision_summary
ORDER BY 
  CASE 
    WHEN nonzero_versions < distinct_versions THEN 0
    WHEN nonzero_versions = 1 THEN 1
    ELSE 2
  END,
  object_name
LIMIT 100;

\echo ''
\echo '=== SUMMARY ==='
WITH filename_hashes AS (
  SELECT 
    om.object_name,
    oh.hash_value,
    COUNT(DISTINCT CASE WHEN om.object_size > 0 THEN om.id END) as nonzero_copies
  FROM object_meta om
  LEFT JOIN object_hash_junction ohj ON om.id = ohj.object_meta_id
  LEFT JOIN object_hash oh ON ohj.object_hash_id = oh.id
  WHERE oh.hash_type = 'md5' AND om.object_size > 0
  GROUP BY om.object_name, oh.hash_value
),
collision_summary AS (
  SELECT 
    COUNT(DISTINCT CASE WHEN COUNT(DISTINCT hash_value) > 1 THEN object_name END) as files_with_collisions,
    COUNT(DISTINCT object_name) as unique_filenames
  FROM filename_hashes
  GROUP BY object_name
)
SELECT 
  'Total filenames with collisions: ' || COALESCE((SELECT COUNT(DISTINCT object_name) FROM filename_hashes WHERE (SELECT COUNT(DISTINCT hash_value) FROM filename_hashes fh2 WHERE fh2.object_name = filename_hashes.object_name) > 1), 0) as info;
