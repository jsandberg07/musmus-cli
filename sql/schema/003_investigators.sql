-- +goose Up
CREATE TABLE investigators(
    id UUID PRIMARY KEY UNIQUE,
    i_name TEXT NOT NULL,
    nickname TEXT,
    email TEXT UNIQUE,
    position UUID NOT NULL REFERENCES positions ON UPDATE CASCADE,
    active BOOLEAN NOT NULL,
    hashed_password TEXT,

    CONSTRAINT fk_investigators
    FOREIGN KEY(position)
    REFERENCES positions(id)
);

-- +goose Down
DROP TABLE investigators;