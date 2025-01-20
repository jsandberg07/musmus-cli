-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY UNIQUE,
    order_number TEXT NOT NULL UNIQUE,
    expected_date TIMESTAMP NOT NULL,
    investigator_id UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    strain_id UUID NOT NULL REFERENCES strains ON UPDATE CASCADE,
    note TEXT,
    received BOOLEAN NOT NULL,

    CONSTRAINT fk_investigators
    FOREIGN KEY(investigator_id)
    REFERENCES investigators(id),

    CONSTRAINT fk_strains
    FOREIGN KEY(strain_id)
    REFERENCES strains(id)
);

-- +goose Down
DROP TABLE orders;