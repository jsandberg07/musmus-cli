-- +goose Up
CREATE TABLE protocols(
    id UUID PRIMARY KEY UNIQUE,
    p_number TEXT UNIQUE,
    primary_investigator UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    title TEXT NOT NULL,
    allocated INTEGER NOT NULL,
    balance INTEGER NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL,
    previous_protocol UUID,

    CONSTRAINT fk_investigators
    FOREIGN KEY(primary_investigator)
    REFERENCES investigators(id)
);

-- +goose Down
DROP TABLE protocols;