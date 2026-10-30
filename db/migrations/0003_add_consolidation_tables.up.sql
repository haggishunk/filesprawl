-- Migration 0003: Add consolidation tables for multi-user workflow
-- Creates tables for persistent flagging and audit trail
-- Enables pause/resume consolidation and multi-user review

BEGIN;

-- User authentication table
-- Tracks who flags files and who reviews decisions
CREATE TABLE IF NOT EXISTS consolidation_users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Persistent flag storage
-- Each flag decision stored immediately with user and timestamp
-- Enables session independence (flag, exit, resume later)
CREATE TABLE IF NOT EXISTS consolidation_flags (
  id SERIAL PRIMARY KEY,
  file_id INT NOT NULL REFERENCES object_meta(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES consolidation_users(id) ON DELETE RESTRICT,
  flag_type TEXT NOT NULL CHECK (flag_type IN ('KEEP', 'ARCHIVE', 'DELETE')),
  flagged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  reason TEXT,
  archive_dest TEXT
);

CREATE INDEX IF NOT EXISTS idx_consolidation_flags_file_id ON consolidation_flags (file_id);
CREATE INDEX IF NOT EXISTS idx_consolidation_flags_user_id ON consolidation_flags (user_id);

-- Audit trail for multi-user review
-- Tracks modifications and approvals of flags
-- Records who reviewed, what they changed, and why
CREATE TABLE IF NOT EXISTS consolidation_reviews (
  id SERIAL PRIMARY KEY,
  file_id INT NOT NULL REFERENCES object_meta(id) ON DELETE CASCADE,
  original_flag TEXT CHECK (original_flag IN ('KEEP', 'ARCHIVE', 'DELETE')),
  original_user INT REFERENCES consolidation_users(id) ON DELETE SET NULL,
  original_at TIMESTAMP,
  reviewed_by INT NOT NULL REFERENCES consolidation_users(id) ON DELETE RESTRICT,
  reviewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  new_flag TEXT CHECK (new_flag IN ('KEEP', 'ARCHIVE', 'DELETE')),
  review_reason TEXT
);

CREATE INDEX IF NOT EXISTS idx_consolidation_reviews_file_id ON consolidation_reviews (file_id);
CREATE INDEX IF NOT EXISTS idx_consolidation_reviews_reviewed_by ON consolidation_reviews (reviewed_by);

COMMIT;
