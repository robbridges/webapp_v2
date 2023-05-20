
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    message TEXT  NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
