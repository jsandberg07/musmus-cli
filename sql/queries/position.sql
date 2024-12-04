-- name: GetPositions :many
SELECT * FROM positions;

-- name: GetUserPosition :one
SELECT * FROM positions
WHERE $1 = id;

-- name: CreatePosition :one
INSERT INTO positions(id, title, can_activate, can_deactivate, can_add_orders, can_query, can_change_protocol, can_add_staff)
VALUES(gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdatePosition :exec
UPDATE positions
SET can_activate = $2,
    can_deactivate = $3,
    can_add_orders = $4,
    can_query = $5,
    can_change_protocol = $6,
    can_add_staff = $7
WHERE $1 = title;

-- name: GetPositionByTitle :one
SELECT * FROM positions
WHERE $1 = title;