CREATE EXTENSION citext;

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

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    password_digest TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE schema_migrations (
	migration CHAR(14)
) ;

INSERT INTO schema_migrations VALUES ("20180611170754") ;
INSERT INTO schema_migrations VALUES ("20181127180911") ;
INSERT INTO schema_migrations VALUES ("20181204181116") ;
