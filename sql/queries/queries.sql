-- name: GetCageCardsOrder :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
WHERE order_id = $1
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsAll :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsActive :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
WHERE cage_cards.activated_on IS NOT NULL and cage_cards.deactivated_on IS NULL
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsInvestigator :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
AND cage_cards.investigator_id = $3
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsProtocol :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
AND cage_cards.protocol_id = $3
ORDER BY cage_cards.cc_id ASC;

-- name: GetCardsDateRange :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on, orders.order_number
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
LEFT JOIN orders ON cage_cards.order_id = orders.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
ORDER BY cage_cards.cc_id ASC;