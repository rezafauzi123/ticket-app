-- +goose Up
ALTER TABLE users
ADD CONSTRAINT unique_email UNIQUE (email);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS unique_email;

