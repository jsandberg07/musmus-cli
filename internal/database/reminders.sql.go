// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reminders.sql

package database

import (
	"context"
	"time"
)

const getAllReminders = `-- name: GetAllReminders :many
SELECT id, r_date, r_cc_id, investigator_id, note FROM reminders
`

func (q *Queries) GetAllReminders(ctx context.Context) ([]Reminder, error) {
	rows, err := q.db.QueryContext(ctx, getAllReminders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reminder
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(
			&i.ID,
			&i.RDate,
			&i.RCcID,
			&i.InvestigatorID,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRemindersDateRange = `-- name: GetRemindersDateRange :many
SELECT id, r_date, r_cc_id, investigator_id, note FROM reminders
WHERE (r_date BETWEEN $1 AND $2)
`

type GetRemindersDateRangeParams struct {
	RDate   time.Time
	RDate_2 time.Time
}

func (q *Queries) GetRemindersDateRange(ctx context.Context, arg GetRemindersDateRangeParams) ([]Reminder, error) {
	rows, err := q.db.QueryContext(ctx, getRemindersDateRange, arg.RDate, arg.RDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reminder
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(
			&i.ID,
			&i.RDate,
			&i.RCcID,
			&i.InvestigatorID,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTodayReminders = `-- name: GetTodayReminders :many
SELECT id, r_date, r_cc_id, investigator_id, note FROM reminders
WHERE r_date = $1
`

func (q *Queries) GetTodayReminders(ctx context.Context, rDate time.Time) ([]Reminder, error) {
	rows, err := q.db.QueryContext(ctx, getTodayReminders, rDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reminder
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(
			&i.ID,
			&i.RDate,
			&i.RCcID,
			&i.InvestigatorID,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
