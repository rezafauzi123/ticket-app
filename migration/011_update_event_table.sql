-- +goose Up
ALTER TABLE events ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE events DROP COLUMN description;
