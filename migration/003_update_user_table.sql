-- +goose Up
ALTER TABLE users
    ADD COLUMN email VARCHAR(50),
    ADD COLUMN password VARCHAR(100);                 

-- +goose Down
DROP TABLE IF EXISTS users;
