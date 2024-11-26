-- +goose Up
CREATE TABLE config(
    first_time_run BOOLEAN NOT NULL DEFAULT true,
    only_activate_self BOOLEAN NOT NULL DEFAULT true
);

-- +goose Down
DROP TABLE config;