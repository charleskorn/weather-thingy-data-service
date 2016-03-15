-- +migrate Up
ALTER TABLE agents ADD COLUMN token VARCHAR(100) NOT NULL;

-- +migrate Down
ALTER TABLE agents DROP COLUMN token;
