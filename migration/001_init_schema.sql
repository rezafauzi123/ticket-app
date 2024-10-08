-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied.
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back.
DROP TABLE IF EXISTS users;
