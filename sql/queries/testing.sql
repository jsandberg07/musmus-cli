-- name: ResetDatabase :exec
TRUNCATE orders, reminders, 
    cage_cards, 
    strains, 
    added_to_protocol, 
    investigators, 
    positions, 
    settings CASCADE;