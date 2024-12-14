-- name: GetStrains :many
SELECT * FROM strains
ORDER BY vendor DESC;

-- name: AddStrain :one
INSERT INTO strains(id, s_name, vendor, vendor_code)
VALUES(gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetStrainByName :one
SELECT * FROM strains
WHERE $1 = vendor_code OR $1 = s_name;

-- name: GetStrainByID :one
SELECT * FROM strains
WHERE $1 = id;

-- name: GetStrainByCode :one
SELECT * FROM strains
WHERE $1 = vendor_code;

-- name: UpdateStrain :exec
UPDATE strains
SET s_name = $2,
    vendor = $3,
    vendor_code = $4
WHERE $1 = id;