-- name: GetPositions :many
SELECT * FROM positions;

-- name: GetUserPosition :one
SELECT * FROM positions
WHERE $1 = id;

-- name: CreatePosition :one
INSERT INTO positions(id, title, can_activate, can_deactivate, can_add_orders, can_receive_orders, can_query, can_change_protocol, can_add_staff, can_add_reminders)
VALUES(gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdatePosition :exec
UPDATE positions
SET title = $1,
    can_activate = $2,
    can_deactivate = $3,
    can_add_orders = $4,
    can_receive_orders = $5,
    can_query = $6,
    can_change_protocol = $7,
    can_add_staff = $8,
    can_add_reminders = $9
WHERE $10 = id;

-- name: GetPositionByTitle :one
SELECT * FROM positions
WHERE $1 = title;