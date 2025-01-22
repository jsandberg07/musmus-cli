-- name: GetTodayReminders :many
SELECT * FROM reminders
WHERE r_date = $1;

-- name: GetAllReminders :many
SELECT * FROM reminders;

-- name: GetRemindersDateRange :many
SELECT * FROM reminders
WHERE (r_date BETWEEN $1 AND $2);

-- name: AddReminder :one
INSERT INTO reminders(id, r_date, r_cc_id, investigator_id, note)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: DeleteReminder :exec
DELETE FROM reminders
WHERE $1 = id;