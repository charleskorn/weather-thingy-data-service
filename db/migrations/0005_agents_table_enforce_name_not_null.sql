-- +migrate Up
ALTER TABLE agents ALTER COLUMN name SET NOT NULL;

-- +migrate Down
ALTER TABLE agents ALTER COLUMN name DROP NOT NULL;
