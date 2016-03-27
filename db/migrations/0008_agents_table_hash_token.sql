-- +migrate Up
ALTER TABLE agents ADD COLUMN token_iterations INT NOT NULL;
ALTER TABLE agents ADD COLUMN token_salt BYTEA NOT NULL;
ALTER TABLE agents ADD COLUMN token_hash BYTEA NOT NULL;
ALTER TABLE agents DROP COLUMN token;

-- +migrate Down
ALTER TABLE agents DROP COLUMN token_iterations;
ALTER TABLE agents DROP COLUMN token_salt;
ALTER TABLE agents DROP COLUMN token_hash;
ALTER TABLE agents ADD COLUMN token VARCHAR(100) NOT NULL;
