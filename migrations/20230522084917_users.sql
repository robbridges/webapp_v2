-- +goose Up
-- +goose StatementBegin
CREATE table users (
   id SERIAL PRIMARY KEY,
   email TEXT UNIQUE NOT NULL,
   password_hash TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd