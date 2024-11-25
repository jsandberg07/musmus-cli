-- +goose Up
CREATE TABLE cage_cards(
    cc_id int PRIMARY KEY UNIQUE,
    activated TIMESTAMP,
    deactivated TIMESTAMP,
    investigator TEXT NOT NULL
);


-- +goose Down
DROP TABLE cage_cards;