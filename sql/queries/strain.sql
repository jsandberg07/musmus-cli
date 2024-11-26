-- name: GetStrains :many
SELECT * FROM strains
ORDER BY vendor DESC;

-- name: AddStrain :one
INSERT INTO strains(id, s_name, vendor, vendor_code)
VALUES(gen_random_uuid(), s_name, vendor, vendor_code)
RETURNING *;

-- name: GetStrainByName :one
SELECT * FROM strains
WHERE $1 = vendor_code OR $1 = s_name;

-- name: getStrainByID :one
SELECT * FROM strains
WHERE $1 = id;

