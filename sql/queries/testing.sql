-- name: ResetDatabase :exec
TRUNCATE cage_cards, 
    strains, 
    added_to_protocol, 
    investigators, 
    positions, 
    settings CASCADE;