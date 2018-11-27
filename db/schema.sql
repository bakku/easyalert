CREATE EXTENSION citext;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    password_digest TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON users (admin);

CREATE TABLE schema_migrations (
	migration CHAR(14)
) ;

INSERT INTO schema_migrations VALUES ("20180611170754") ;
