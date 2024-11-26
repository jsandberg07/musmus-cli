-- this is bunk and it's here just to not break sqlc

-- name: ActivateCageCard :one
INSERT INTO cage_cards(cc_id, activated_on, deactivated_on, investigator_id)
VALUES($1, $2, NULL, $3)
RETURNING *;

-- name: GetCageCards :many
SELECT * FROM cage_cards
ORDER BY cc_id ASC;

-- name: ResetCageCards :exec
DELETE FROM cage_cards *;