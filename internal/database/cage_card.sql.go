// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: cage_card.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addCageCard = `-- name: AddCageCard :one
INSERT INTO cage_cards(cc_id, protocol_id, investigator_id)
VALUES ($1, $2, $3)
RETURNING cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by
`

type AddCageCardParams struct {
	CcID           int32
	ProtocolID     uuid.UUID
	InvestigatorID uuid.UUID
}

func (q *Queries) AddCageCard(ctx context.Context, arg AddCageCardParams) (CageCard, error) {
	row := q.db.QueryRowContext(ctx, addCageCard, arg.CcID, arg.ProtocolID, arg.InvestigatorID)
	var i CageCard
	err := row.Scan(
		&i.CcID,
		&i.ProtocolID,
		&i.ActivatedOn,
		&i.DeactivatedOn,
		&i.InvestigatorID,
		&i.Strain,
		&i.Notes,
		&i.ActivatedBy,
		&i.DeactivatedBy,
	)
	return i, err
}

const addNote = `-- name: AddNote :exec
UPDATE cage_cards
SET notes = $2
WHERE cc_id = $1
`

type AddNoteParams struct {
	CcID  int32
	Notes sql.NullString
}

func (q *Queries) AddNote(ctx context.Context, arg AddNoteParams) error {
	_, err := q.db.ExecContext(ctx, addNote, arg.CcID, arg.Notes)
	return err
}

const deactivateCageCard = `-- name: DeactivateCageCard :one
UPDATE cage_cards
SET deactivated_on = $2,
    deactivated_by = $3
WHERE cc_id = $1
RETURNING cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by
`

type DeactivateCageCardParams struct {
	CcID          int32
	DeactivatedOn sql.NullTime
	DeactivatedBy uuid.NullUUID
}

func (q *Queries) DeactivateCageCard(ctx context.Context, arg DeactivateCageCardParams) (CageCard, error) {
	row := q.db.QueryRowContext(ctx, deactivateCageCard, arg.CcID, arg.DeactivatedOn, arg.DeactivatedBy)
	var i CageCard
	err := row.Scan(
		&i.CcID,
		&i.ProtocolID,
		&i.ActivatedOn,
		&i.DeactivatedOn,
		&i.InvestigatorID,
		&i.Strain,
		&i.Notes,
		&i.ActivatedBy,
		&i.DeactivatedBy,
	)
	return i, err
}

const getActivationDate = `-- name: GetActivationDate :one
SELECT activated_on FROM cage_cards
WHERE $1 = cc_id
`

func (q *Queries) GetActivationDate(ctx context.Context, ccID int32) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, getActivationDate, ccID)
	var activated_on sql.NullTime
	err := row.Scan(&activated_on)
	return activated_on, err
}

const getActiveTestCards = `-- name: GetActiveTestCards :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE cage_cards.activated_on IS NOT NULL and cage_cards.deactivated_on IS NULL
ORDER BY cage_cards.cc_id ASC
`

type GetActiveTestCardsRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetActiveTestCards(ctx context.Context) ([]GetActiveTestCardsRow, error) {
	rows, err := q.db.QueryContext(ctx, getActiveTestCards)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetActiveTestCardsRow
	for rows.Next() {
		var i GetActiveTestCardsRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getAllActiveCageCards = `-- name: GetAllActiveCageCards :many
SELECT cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by FROM cage_cards
WHERE activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC
`

func (q *Queries) GetAllActiveCageCards(ctx context.Context) ([]CageCard, error) {
	rows, err := q.db.QueryContext(ctx, getAllActiveCageCards)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CageCard
	for rows.Next() {
		var i CageCard
		if err := rows.Scan(
			&i.CcID,
			&i.ProtocolID,
			&i.ActivatedOn,
			&i.DeactivatedOn,
			&i.InvestigatorID,
			&i.Strain,
			&i.Notes,
			&i.ActivatedBy,
			&i.DeactivatedBy,
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

const getAllCageCards = `-- name: GetAllCageCards :many
SELECT cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by FROM cage_cards
ORDER BY cc_id ASC
`

func (q *Queries) GetAllCageCards(ctx context.Context) ([]CageCard, error) {
	rows, err := q.db.QueryContext(ctx, getAllCageCards)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CageCard
	for rows.Next() {
		var i CageCard
		if err := rows.Scan(
			&i.CcID,
			&i.ProtocolID,
			&i.ActivatedOn,
			&i.DeactivatedOn,
			&i.InvestigatorID,
			&i.Strain,
			&i.Notes,
			&i.ActivatedBy,
			&i.DeactivatedBy,
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

const getCageCardByID = `-- name: GetCageCardByID :one
SELECT cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by FROM cage_cards
WHERE $1 = cc_id
`

func (q *Queries) GetCageCardByID(ctx context.Context, ccID int32) (CageCard, error) {
	row := q.db.QueryRowContext(ctx, getCageCardByID, ccID)
	var i CageCard
	err := row.Scan(
		&i.CcID,
		&i.ProtocolID,
		&i.ActivatedOn,
		&i.DeactivatedOn,
		&i.InvestigatorID,
		&i.Strain,
		&i.Notes,
		&i.ActivatedBy,
		&i.DeactivatedBy,
	)
	return i, err
}

const getCageCardsActive = `-- name: GetCageCardsActive :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE cage_cards.activated_on IS NOT NULL and cage_cards.deactivated_on IS NULL
ORDER BY cage_cards.cc_id ASC
`

type GetCageCardsActiveRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetCageCardsActive(ctx context.Context) ([]GetCageCardsActiveRow, error) {
	rows, err := q.db.QueryContext(ctx, getCageCardsActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCageCardsActiveRow
	for rows.Next() {
		var i GetCageCardsActiveRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getCageCardsAll = `-- name: GetCageCardsAll :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
ORDER BY cage_cards.cc_id ASC
`

type GetCageCardsAllRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetCageCardsAll(ctx context.Context) ([]GetCageCardsAllRow, error) {
	rows, err := q.db.QueryContext(ctx, getCageCardsAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCageCardsAllRow
	for rows.Next() {
		var i GetCageCardsAllRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getCageCardsByInvestigator = `-- name: GetCageCardsByInvestigator :many
SELECT cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by FROM cage_cards
WHERE $1 = investigator_id
AND activated_on IS NOT NULL AND deactivated_on IS NULL
ORDER BY cc_id ASC
`

func (q *Queries) GetCageCardsByInvestigator(ctx context.Context, investigatorID uuid.UUID) ([]CageCard, error) {
	rows, err := q.db.QueryContext(ctx, getCageCardsByInvestigator, investigatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CageCard
	for rows.Next() {
		var i CageCard
		if err := rows.Scan(
			&i.CcID,
			&i.ProtocolID,
			&i.ActivatedOn,
			&i.DeactivatedOn,
			&i.InvestigatorID,
			&i.Strain,
			&i.Notes,
			&i.ActivatedBy,
			&i.DeactivatedBy,
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

const getCageCardsInvestigator = `-- name: GetCageCardsInvestigator :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on <= $1) AND (deactivated_on >= $2 OR deactivated_on IS NULL)
AND investigators.i_name = $3
ORDER BY cage_cards.cc_id ASC
`

type GetCageCardsInvestigatorParams struct {
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
	IName         string
}

type GetCageCardsInvestigatorRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetCageCardsInvestigator(ctx context.Context, arg GetCageCardsInvestigatorParams) ([]GetCageCardsInvestigatorRow, error) {
	rows, err := q.db.QueryContext(ctx, getCageCardsInvestigator, arg.ActivatedOn, arg.DeactivatedOn, arg.IName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCageCardsInvestigatorRow
	for rows.Next() {
		var i GetCageCardsInvestigatorRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getCageCardsProtocol = `-- name: GetCageCardsProtocol :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on <= $1) AND (deactivated_on >= $2 OR deactivated_on IS NULL)
AND protocols.p_number = $3
ORDER BY cage_cards.cc_id ASC
`

type GetCageCardsProtocolParams struct {
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
	PNumber       string
}

type GetCageCardsProtocolRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetCageCardsProtocol(ctx context.Context, arg GetCageCardsProtocolParams) ([]GetCageCardsProtocolRow, error) {
	rows, err := q.db.QueryContext(ctx, getCageCardsProtocol, arg.ActivatedOn, arg.DeactivatedOn, arg.PNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCageCardsProtocolRow
	for rows.Next() {
		var i GetCageCardsProtocolRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getCardsDateRange = `-- name: GetCardsDateRange :many
SELECT cage_cards.cc_id, investigators.i_name, protocols.p_number, strains.s_name, cage_cards.activated_on, cage_cards.deactivated_on
FROM cage_cards
INNER JOIN investigators ON cage_cards.investigator_id = investigators.id
INNER JOIN protocols ON cage_cards.protocol_id = protocols.id
LEFT JOIN strains ON cage_cards.strain = strains.id
WHERE (activated_on IS NOT NULL AND activated_on <= $1) AND (deactivated_on >= $2 OR deactivated_on IS NULL)
ORDER BY cage_cards.cc_id ASC
`

type GetCardsDateRangeParams struct {
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

type GetCardsDateRangeRow struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

func (q *Queries) GetCardsDateRange(ctx context.Context, arg GetCardsDateRangeParams) ([]GetCardsDateRangeRow, error) {
	rows, err := q.db.QueryContext(ctx, getCardsDateRange, arg.ActivatedOn, arg.DeactivatedOn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCardsDateRangeRow
	for rows.Next() {
		var i GetCardsDateRangeRow
		if err := rows.Scan(
			&i.CcID,
			&i.IName,
			&i.PNumber,
			&i.SName,
			&i.ActivatedOn,
			&i.DeactivatedOn,
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

const getDeactivationDate = `-- name: GetDeactivationDate :one
SELECT deactivated_on FROM cage_cards
WHERE $1 = cc_id
`

func (q *Queries) GetDeactivationDate(ctx context.Context, ccID int32) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, getDeactivationDate, ccID)
	var deactivated_on sql.NullTime
	err := row.Scan(&deactivated_on)
	return deactivated_on, err
}

const inactivateCageCard = `-- name: InactivateCageCard :exec
UPDATE cage_cards
SET activated_on = NULL
WHERE $1 = cc_id
`

func (q *Queries) InactivateCageCard(ctx context.Context, ccID int32) error {
	_, err := q.db.ExecContext(ctx, inactivateCageCard, ccID)
	return err
}

const newActivateCageCard = `-- name: NewActivateCageCard :exec
UPDATE cage_cards
SET activated_on = $2,
    activated_by = $3
WHERE cc_id = $1
`

type NewActivateCageCardParams struct {
	CcID        int32
	ActivatedOn sql.NullTime
	ActivatedBy uuid.NullUUID
}

func (q *Queries) NewActivateCageCard(ctx context.Context, arg NewActivateCageCardParams) error {
	_, err := q.db.ExecContext(ctx, newActivateCageCard, arg.CcID, arg.ActivatedOn, arg.ActivatedBy)
	return err
}

const reactivateCageCard = `-- name: ReactivateCageCard :exec
UPDATE cage_cards
SET deactivated_on = NULL
WHERE $1 = cc_id
`

func (q *Queries) ReactivateCageCard(ctx context.Context, ccID int32) error {
	_, err := q.db.ExecContext(ctx, reactivateCageCard, ccID)
	return err
}

const trueActivateCageCard = `-- name: TrueActivateCageCard :one
UPDATE cage_cards
SET activated_on = $2,
    activated_by = $3,
    strain = $4,
    notes = $5
WHERE cc_id = $1
RETURNING cc_id, protocol_id, activated_on, deactivated_on, investigator_id, strain, notes, activated_by, deactivated_by
`

type TrueActivateCageCardParams struct {
	CcID        int32
	ActivatedOn sql.NullTime
	ActivatedBy uuid.NullUUID
	Strain      uuid.NullUUID
	Notes       sql.NullString
}

func (q *Queries) TrueActivateCageCard(ctx context.Context, arg TrueActivateCageCardParams) (CageCard, error) {
	row := q.db.QueryRowContext(ctx, trueActivateCageCard,
		arg.CcID,
		arg.ActivatedOn,
		arg.ActivatedBy,
		arg.Strain,
		arg.Notes,
	)
	var i CageCard
	err := row.Scan(
		&i.CcID,
		&i.ProtocolID,
		&i.ActivatedOn,
		&i.DeactivatedOn,
		&i.InvestigatorID,
		&i.Strain,
		&i.Notes,
		&i.ActivatedBy,
		&i.DeactivatedBy,
	)
	return i, err
}
