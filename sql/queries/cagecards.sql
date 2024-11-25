-- name: ActivateCageCard :one
INSERT INTO cage_cards(cc_id, activated, deactivated, investigator)
VALUES($1, $2, NULL, $3)
RETURNING *;

-- name: GetCageCards :many
SELECT * FROM cage_cards
ORDER BY cc_id ASC;

-- name: ResetCageCards :exec
DELETE FROM cage_cards *;