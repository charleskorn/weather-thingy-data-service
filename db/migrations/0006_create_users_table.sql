-- +migrate Up
CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(254) NOT NULL UNIQUE,
  password_iterations INT NOT NULL,
  password_salt BYTEA NOT NULL,
  password_hash BYTEA NOT NULL,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE agents ADD COLUMN owner_user_id INT NOT NULL REFERENCES users (user_id);

-- +migrate Down
ALTER TABLE agents DROP COLUMN owner_id;

DROP TABLE users;
