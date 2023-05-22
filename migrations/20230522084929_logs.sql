-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
      id SERIAL PRIMARY KEY,
      message TEXT  NOT NULL,
      timestamp TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
