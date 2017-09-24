DROP SEQUENCE IF EXISTS app.users_id_seq CASCADE;
CREATE SEQUENCE app.users_id_seq;

DROP TABLE IF EXISTS app.users;
CREATE TABLE app.users
(
  id INT PRIMARY KEY NOT NULL DEFAULT NEXTVAL('app.users_id_seq'),
  username VARCHAR(30) NOT NULL,
  password BYTEA NOT NULL,
  verified SMALLINT NOT NULL DEFAULT 0,
  disabled SMALLINT NOT NULL DEFAULT 0,
  email VARCHAR(255) NOT NULL
);

DROP INDEX IF EXISTS udx_user_username;
CREATE UNIQUE INDEX udx_user_username
  ON app.users (username);

DROP INDEX IF EXISTS udx_user_email;
CREATE UNIQUE INDEX udx_user_email
  ON app.users (email);

DROP TABLE IF EXISTS app.user_verifications;
CREATE TABLE app.user_verifications
(
  user_id INT PRIMARY KEY NOT NULL,
  hash BYTEA NOT NULL,
  expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 day'
);

DROP INDEX IF EXISTS udx_user_verification_hash;
CREATE UNIQUE INDEX udx_user_verification_hash
  ON app.user_verifications (hash);

DROP TABLE IF EXISTS app.user_sessions;
CREATE TABLE app.user_sessions
(
  user_id INT PRIMARY KEY NOT NULL,
  character_id INT NOT NULL DEFAULT 0,
  key BYTEA NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DROP INDEX IF EXISTS idx_user_sessions_key;
CREATE INDEX idx_user_sessions_key
  ON app.user_sessions (key);

DROP TABLE IF EXISTS app.user_authorizations;
CREATE TABLE app.user_authorizations
(
  user_id INT NOT NULL,
  character_id INT NOT NULL,
  role INT NOT NULL,
  token BYTEA NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(user_id, role)
);