BEGIN;
  CREATE TABLE alerts (
    id BIGSERIAL PRIMARY KEY,
    subject TEXT NOT NULL,
    status smallint NOT NULL,
    sent_at TIMESTAMP DEFAULT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
  );

  CREATE INDEX ON alerts (user_id);
COMMIT;
