-- +goose Up
CREATE TABLE cage_cards(
    cc_id int PRIMARY KEY UNIQUE,
    activated_on TIMESTAMP,
    deactivated_on TIMESTAMP,
    investigator_id UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    strain UUID REFERENCES strains ON UPDATE CASCADE,
    notes TEXT,
    activated_by UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    deactivated_by UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,

    CONSTRAINT fk_strain
    FOREIGN KEY(strain)
    REFERENCES strains(id),

    CONSTRAINT fk_activated_by
    FOREIGN KEY(activated_by)
    REFERENCES investigators(id),

    CONSTRAINT fk_deactivated_by
    FOREIGN KEY(deactivated_by)
    REFERENCES investigators(id)
);


-- +goose Down
DROP TABLE cage_cards;