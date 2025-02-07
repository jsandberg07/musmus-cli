-- name: GetAllOrders :many
SELECT * FROM orders;

-- name: GetAllExpectedOrders :many
SELECT * FROM orders
WHERE received = false;

-- name: CreateNewOrder :one
INSERT INTO orders(id, order_number, expected_date, protocol_id, investigator_id, strain_id, note, received)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, false)
RETURNING *;

-- name: GetOrderDateRange :many
SELECT * FROM orders
WHERE (expected_date BETWEEN $1 AND $2);

-- name: GetOrderExpectedToday :many
SELECT * FROM orders
WHERE (expected_date = $1) AND received = false;

-- name: MarkOrderReceived :one
UPDATE orders
SET received = true
WHERE id = $1
RETURNING *;


-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderByNumber :one
SELECT * FROM orders
WHERE order_number = $1;

-- name: UpdateOrder :exec
UPDATE orders
SET expected_date = $2,
    protocol_id = $3,
    investigator_id = $4,
    strain_id = $5,
    note = $6
WHERE $1 = id;