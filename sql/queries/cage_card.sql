-- name: GetCageCardsByInvestigator :many
SELECT * FROM cage_cards
WHERE $1 = investigator_id
AND activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC;

-- name: GetAllActiveCageCards :many
SELECT * FROM cage_cards
WHERE activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC;

-- name: GetCageCardsRange :many
SELECT * FROM cage_cards
WHERE cc_id >= $1 AND cc_id <= $2;

-- name: GetAllCageCards :many
SELECT * FROM cage_cards
ORDER BY cc_id ASC;

-- name: AddCageCard :one
INSERT INTO cage_cards(cc_id, protocol_id, investigator_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: NewActivateCageCard :exec
UPDATE cage_cards
SET activated_on = $2,
    activated_by = $3
WHERE cc_id = $1;

-- name: DeactivateCageCard :one
UPDATE cage_cards
SET deactivated_on = $2,
    deactivated_by = $3
WHERE cc_id = $1
RETURNING *;

-- name: AddNote :exec
UPDATE cage_cards
SET notes = $2
WHERE cc_id = $1;

-- name: TrueActivateCageCard :one
UPDATE cage_cards
SET activated_on = $2,
    activated_by = $3,
    strain = $4,
    notes = $5
WHERE cc_id = $1
RETURNING *;

-- name: GetCageCardByID :one
SELECT * FROM cage_cards
WHERE $1 = cc_id;

-- name: GetActivationDate :one
SELECT activated_on FROM cage_cards
WHERE $1 = cc_id;

-- name: GetDeactivationDate :one
SELECT deactivated_on FROM cage_cards
WHERE $1 = cc_id;

-- name: ReactivateCageCard :exec
UPDATE cage_cards
SET deactivated_on = NULL
WHERE $1 = cc_id;

-- name: InactivateCageCard :exec
UPDATE cage_cards
SET activated_on = NULL
WHERE $1 = cc_id;

-- name: GetActiveTestCards :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE cage_cards.activated_on IS NOT NULL and cage_cards.deactivated_on IS NULL
ORDER BY cage_cards.cc_id ASC;

-- name: GetCardsDateRange :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsProtocol :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
AND cage_cards.protocol_id = $3
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsInvestigator :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on >= $1) AND (deactivated_on <= $2 OR deactivated_on IS NULL)
AND cage_cards.investigator_id = $3
ORDER BY cage_cards.cc_id ASC;


-- name: GetCageCardsActive :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE cage_cards.activated_on IS NOT NULL and cage_cards.deactivated_on IS NULL
ORDER BY cage_cards.cc_id ASC;

-- name: GetCageCardsAll :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
ORDER BY cage_cards.cc_id ASC;

-- name: ReceiveCageCard :one
INSERT INTO cage_cards(cc_id, protocol_id, activated_on, investigator_id, strain, notes, activated_by)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING *;