CREATE TYPE hash_type_enum AS ENUM ('md5', 'dropbox', 'sha1', 'sha256');

CREATE TABLE IF NOT EXISTS object_hash (
    id SERIAL PRIMARY KEY,
    hash_value TEXT NOT NULL,
    hash_type hash_type_enum NOT NULL
);

CREATE TABLE IF NOT EXISTS remote (
  id SERIAL PRIMARY KEY,
  remote_name TEXT NOT NULL,
  remote_type TEXT NOT NULL,
  hostname TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS object_meta (
  id SERIAL PRIMARY KEY,
  object_id TEXT NOT NULL,
  object_name TEXT NOT NULL,
  object_path TEXT NOT NULL,
  object_mime_type TEXT NOT NULL,
  object_size BIGINT,
  remote_id INT,
  FOREIGN KEY (remote_id) REFERENCES remote(id)
);

CREATE TABLE IF NOT EXISTS object_hash_junction (
  id SERIAL PRIMARY KEY,
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
