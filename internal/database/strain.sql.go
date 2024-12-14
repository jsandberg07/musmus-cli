// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: strain.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addStrain = `-- name: AddStrain :one
INSERT INTO strains(id, s_name, vendor, vendor_code)
VALUES(gen_random_uuid(), $1, $2, $3)
RETURNING id, s_name, vendor, vendor_code
`

type AddStrainParams struct {
	SName      string
	Vendor     string
	VendorCode string
}

func (q *Queries) AddStrain(ctx context.Context, arg AddStrainParams) (Strain, error) {
	row := q.db.QueryRowContext(ctx, addStrain, arg.SName, arg.Vendor, arg.VendorCode)
	var i Strain
	err := row.Scan(
		&i.ID,
		&i.SName,
		&i.Vendor,
		&i.VendorCode,
	)
	return i, err
}

const getStrainByCode = `-- name: GetStrainByCode :one
SELECT id, s_name, vendor, vendor_code FROM strains
WHERE $1 = vendor_code
`

func (q *Queries) GetStrainByCode(ctx context.Context, vendorCode string) (Strain, error) {
	row := q.db.QueryRowContext(ctx, getStrainByCode, vendorCode)
	var i Strain
	err := row.Scan(
		&i.ID,
		&i.SName,
		&i.Vendor,
		&i.VendorCode,
	)
	return i, err
}

const getStrainByID = `-- name: GetStrainByID :one
SELECT id, s_name, vendor, vendor_code FROM strains
WHERE $1 = id
`

func (q *Queries) GetStrainByID(ctx context.Context, id uuid.UUID) (Strain, error) {
	row := q.db.QueryRowContext(ctx, getStrainByID, id)
	var i Strain
	err := row.Scan(
		&i.ID,
		&i.SName,
		&i.Vendor,
		&i.VendorCode,
	)
	return i, err
}

const getStrainByName = `-- name: GetStrainByName :one
SELECT id, s_name, vendor, vendor_code FROM strains
WHERE $1 = vendor_code OR $1 = s_name
`

func (q *Queries) GetStrainByName(ctx context.Context, vendorCode string) (Strain, error) {
	row := q.db.QueryRowContext(ctx, getStrainByName, vendorCode)
	var i Strain
	err := row.Scan(
		&i.ID,
		&i.SName,
		&i.Vendor,
		&i.VendorCode,
	)
	return i, err
}

const getStrains = `-- name: GetStrains :many
SELECT id, s_name, vendor, vendor_code FROM strains
ORDER BY vendor DESC
`

func (q *Queries) GetStrains(ctx context.Context) ([]Strain, error) {
	rows, err := q.db.QueryContext(ctx, getStrains)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Strain
	for rows.Next() {
		var i Strain
		if err := rows.Scan(
			&i.ID,
			&i.SName,
			&i.Vendor,
			&i.VendorCode,
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
