BEGIN;
  ALTER TABLE users
  ADD COLUMN admin BOOLEAN NOT NULL DEFAULT FALSE;

  CREATE INDEX ON users (admin);
COMMIT;