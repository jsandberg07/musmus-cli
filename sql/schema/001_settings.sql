-- settings is saving the settings. this should ONLY ever have one row

-- +goose Up
CREATE TABLE settings(
    id INT PRIMARY KEY UNIQUE,
    settings_complete BOOLEAN NOT NULL DEFAULT false,
    only_activate_self BOOLEAN NOT NULL DEFAULT false,
    test_data_loaded BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE settings;