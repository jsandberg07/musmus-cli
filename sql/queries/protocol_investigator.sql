-- name: AddInvestigatorToProtocol :one
INSERT INTO added_to_protocol(investigator_id, protocol_id)
VALUES($1, $2)
RETURNING *;

-- name: RemoveInvestigatorFromProtocol :exec
DELETE FROM added_to_protocol
WHERE $1 = investigator_id AND $2 = protocol_id;

