-- name: GetInvestigatorByName :many
SELECT * FROM investigators
WHERE $1 = i_name OR $1 = nickname;

-- name: GetInvestigatorByEmail :one
SELECT * FROM investigators
WHERE $1 = email;

-- name: CreateInvestigator :one
INSERT INTO investigators(id, i_name, nickname, email, position, active)
VALUES(gen_random_uuid(), $1, $2, $3, $4, true)
RETURNING *;

-- name: UpdateInvestigator :exec
UPDATE investigators
SET i_name = $2,
    nickname = $3,
    email = $4,
    position = $5,
    active = $6
WHERE $1 = id;

-- name: GetInvestigatorByID :one
SELECT * FROM investigators
WHERE $1 = id;