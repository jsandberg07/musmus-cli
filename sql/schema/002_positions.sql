-- +goose Up
CREATE TABLE positions(
    id UUID PRIMARY KEY UNIQUE,
    title TEXT NOT NULL UNIQUE,
    can_activate BOOLEAN NOT NULL DEFAULT false,
    can_deactivate BOOLEAN NOT NULL DEFAULT false,
    can_add_orders BOOLEAN NOT NULL DEFAULT false,
    can_receive_orders BOOLEAN NOT NULL DEFAULT false,
    can_query BOOLEAN NOT NULL DEFAULT false,
    can_change_protocol BOOLEAN NOT NULL DEFAULT false,
    can_add_staff BOOLEAN NOT NULL DEFAULT false,
    can_add_reminders BOOLEAN NOT NULL DEFAULT false,
    is_admin BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE positions;