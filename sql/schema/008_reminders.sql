-- +goose Up
CREATE TABLE reminders(
    id UUID PRIMARY KEY UNIQUE,
    r_date TIMESTAMP NOT NULL,
    r_cc_id INT NOT NULL REFERENCES cage_cards ON UPDATE CASCADE ON DELETE CASCADE,
    investigator_id UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    note TEXT NOT NULL,

    CONSTRAINT fk_cage_cards
    FOREIGN KEY(r_cc_id)
    REFERENCES cage_cards(cc_id),

    CONSTRAINT fk_investigators
    FOREIGN KEY(investigator_id)
    REFERENCES investigators(id)
);

-- +goose Down
DROP TABLE reminders;