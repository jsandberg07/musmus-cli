-- name: GetAllTodayReminders :many
SELECT * FROM reminders
WHERE r_date = $1;

-- name: GetUserTodayReminders :many
SELECT * FROM reminders
WHERE r_date = $1 AND investigator_id = $2;

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

-- name: GetRemindersByCC :many
SELECT * FROM reminders
WHERE r_cc_id = $1
ORDER BY r_date;

-- name: GetUserDayReminder :many
SELECT * FROM reminders
WHERE investigator_id = $1
AND r_date = $2;