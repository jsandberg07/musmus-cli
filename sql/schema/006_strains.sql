-- +goose Up
CREATE TABLE strains(
    id UUID PRIMARY KEY,
    s_name TEXT NOT NULL,
    vendor TEXT NOT NULL,
    vendor_code TEXT NOT NULL
);

-- +goose Down
DROP TABLE strains;