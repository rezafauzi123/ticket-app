-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

UPDATE users SET new_id = uuid_generate_v4();

ALTER TABLE users ALTER COLUMN id DROP DEFAULT;

ALTER TABLE users ALTER COLUMN id TYPE UUID USING (uuid_generate_v4());

UPDATE users SET id = new_id;

ALTER TABLE users DROP COLUMN new_id;

-- +goose Down
ALTER TABLE users ADD COLUMN new_id SERIAL;

UPDATE users SET new_id = nextval('users_id_seq');

ALTER TABLE users ALTER COLUMN id TYPE SERIAL;

UPDATE users SET id = new_id;

ALTER TABLE users DROP COLUMN new_id;
