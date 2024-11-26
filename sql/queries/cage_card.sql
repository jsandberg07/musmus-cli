-- name: GetCageCardsByInvestigator :many
SELECT * FROM cage_cards
WHERE $1 = investigator_id;

