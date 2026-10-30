CREATE TYPE hash_type_enum AS ENUM ('md5', 'dropbox', 'sha1', 'sha256');

CREATE TABLE IF NOT EXISTS object_hash (
    id SERIAL PRIMARY KEY,
    hash_value TEXT NOT NULL,
    hash_type hash_type_enum NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_object_hash_value ON object_hash (hash_value);

CREATE TABLE IF NOT EXISTS location (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('local', 'remote')),
  UNIQUE(name, type)
);

CREATE TABLE IF NOT EXISTS remote (
  id SERIAL PRIMARY KEY,
  location_id INT NOT NULL REFERENCES location(id) ON DELETE RESTRICT,
  remote_name TEXT NOT NULL,
  scan_root TEXT
);

CREATE TABLE IF NOT EXISTS object_meta (
  id SERIAL PRIMARY KEY,
  object_name TEXT NOT NULL,
  object_path TEXT NOT NULL,
  object_mime_type TEXT NOT NULL,
  object_size BIGINT NOT NULL DEFAULT 0,
  is_symlink BOOLEAN NOT NULL DEFAULT false,
  scan_root TEXT
);

CREATE INDEX IF NOT EXISTS idx_object_meta_scan_root ON object_meta (scan_root);

CREATE TABLE IF NOT EXISTS object_hash_junction (
  id SERIAL PRIMARY KEY,
  scan_time TIMESTAMP,
  object_meta_id INT,
  object_hash_id INT,
  FOREIGN KEY (object_meta_id) REFERENCES object_meta(id),
  FOREIGN KEY (object_hash_id) REFERENCES object_hash(id)
);

CREATE TABLE IF NOT EXISTS object_remote_junction (
  id SERIAL PRIMARY KEY,
  object_meta_id INT,
  remote_id INT,
  FOREIGN KEY (object_meta_id) REFERENCES object_meta(id),
  FOREIGN KEY (remote_id) REFERENCES remote(id)
);

-- Consolidation tables for multi-user workflow
CREATE TABLE IF NOT EXISTS consolidation_users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
