-- name: GetCageCardsByInvestigator :many
SELECT * FROM cage_cards
WHERE $1 = investigator_id
AND activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC;

-- name: GetAllActiveCageCards :many
SELECT * FROM cage_cards
WHERE activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC;

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

-- name: DeactivateCageCard :exec
UPDATE cage_cards
SET deactivated_on = $2,
    deactivated_by = $3
WHERE cc_id = $1;

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

    