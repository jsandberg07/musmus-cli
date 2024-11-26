-- name: ResetDatabase :exec
DELETE FROM cage_cards *;
DELETE FROM strains *;
DELETE FROM added_to_protocol *;
DELETE FROM protocols *;
DELETE FROM investigators *;
DELETE FROM positions *;
DELETE FROM config *;