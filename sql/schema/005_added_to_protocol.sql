-- +goose Up
CREATE TABLE added_to_protocol(
    id UUID PRIMARY KEY,
    investigator_id UUID NOT NULL REFERENCES investigators ON UPDATE CASCADE,
    protocol_id UUID NOT NULL REFERENCES protocols ON UPDATE CASCADE,
    UNIQUE(investigator_id, protocol_id),

    CONSTRAINT fk_atp_investigator
    FOREIGN KEY(investigator_id)
    REFERENCES investigators(id),

    CONSTRAINT fk_apt_protocol
    FOREIGN KEY(protocol_id)
    REFERENCES protocols(id)

);

-- +goose Down
DROP TABLE added_to_protocol;