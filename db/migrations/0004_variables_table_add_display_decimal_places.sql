-- +migrate Up
ALTER TABLE variables ADD COLUMN display_decimal_places INT NOT NULL DEFAULT (0);
ALTER TABLE variables ALTER COLUMN display_decimal_places DROP DEFAULT;

-- +migrate Down
ALTER TABLE variables DROP COLUMN display_decimal_places;
