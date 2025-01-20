-- name: GetTodayReminders :many
SELECT * FROM reminders
WHERE r_date = $1;

-- name: GetAllReminders :many
SELECT * FROM reminders;

-- name: GetRemindersDateRange :many
SELECT * FROM reminders
WHERE (r_date BETWEEN $1 AND $2);