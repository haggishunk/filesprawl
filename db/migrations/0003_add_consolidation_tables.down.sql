-- Rollback migration 0003: Remove consolidation tables

BEGIN;

-- Drop tables in reverse order (reviews before flags before users)
-- to avoid foreign key constraint violations

DROP TABLE IF EXISTS consolidation_reviews CASCADE;

DROP TABLE IF EXISTS consolidation_flags CASCADE;

DROP TABLE IF EXISTS consolidation_users CASCADE;

COMMIT;
