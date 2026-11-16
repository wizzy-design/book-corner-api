-- +goose Up
ALTER TABLE tags DROP COLUMN IF EXISTS kind;

-- +goose Down
ALTER TABLE tags ADD COLUMN kind VARCHAR(50) NOT NULL DEFAULT 'genre';