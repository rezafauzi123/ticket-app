-- +goose Up
ALTER TABLE users ADD COLUMN role_id VARCHAR(50);

-- +goose Down
ALTER TABLE users DROP COLUMN role_id;
