-- name: GetProtocols :many
SELECT * FROM protocols
ORDER BY p_number DESC;

-- name: GetProtocolByNumber :one
SELECT * FROM protocols
WHERE $1 = p_number;

-- name: AddBalance :exec
UPDATE protocols SET balance = (balance + $2)
WHERE $1 = p_number;

-- name: UpdateAllocated :exec
UPDATE protocols SET allocated = $2
WHERE $1 = p_number;

-- name: CreateProtocol :one
INSERT INTO protocols(id, p_number, primary_investigator, title, allocated, balance, expiration_date, is_active, previous_protocol)
VALUES(gen_random_uuid(), $1, $2, $3, $4, $5, $6, true, $7)
RETURNING *;

-- name: UpdateProtocol :exec
UPDATE protocols
SET p_number = $2,
    primary_investigator = $3,
    title = $4,
    allocated = $5,
    balance = $6,
    expiration_date = $7
WHERE $1 = id;
