-- +goose Up
ALTER TABLE cage_cards
ADD COLUMN order_id UUID REFERENCES orders ON UPDATE CASCADE;
ALTER TABLE cage_cards
ADD CONSTRAINT fk_order
    FOREIGN KEY(order_id)
    REFERENCES orders(id);

-- +goose Down
ALTER TABLE cage_cards
DROP COLUMN order_id CASCADE;