-- +goose Up
CREATE TABLE positions(
    id UUID PRIMARY KEY UNIQUE,
    title TEXT NOT NULL UNIQUE,
    can_activate BOOLEAN NOT NULL,
    can_deactivate BOOLEAN NOT NULL,
    can_add_orders BOOLEAN NOT NULL,
    can_query BOOLEAN NOT NULL,
    can_change_protocol BOOLEAN NOT NULL,
    can_add_staff BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE positions;