// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AddedToProtocol struct {
	ID             uuid.UUID
	InvestigatorID uuid.UUID
	ProtocolID     uuid.UUID
}

type CageCard struct {
	CcID           int32
	ProtocolID     uuid.UUID
	ActivatedOn    sql.NullTime
	DeactivatedOn  sql.NullTime
	InvestigatorID uuid.UUID
	Strain         uuid.NullUUID
	Notes          sql.NullString
	ActivatedBy    uuid.NullUUID
	DeactivatedBy  uuid.NullUUID
	OrderID        uuid.NullUUID
}

type Investigator struct {
	ID             uuid.UUID
	IName          string
	Nickname       sql.NullString
	Email          sql.NullString
	Position       uuid.UUID
	Active         bool
	HashedPassword sql.NullString
}

type Order struct {
	ID             uuid.UUID
	OrderNumber    string
	ExpectedDate   time.Time
	ProtocolID     uuid.UUID
	InvestigatorID uuid.UUID
	StrainID       uuid.UUID
	Note           sql.NullString
	Received       bool
}

type Position struct {
	ID                uuid.UUID
	Title             string
	CanActivate       bool
	CanDeactivate     bool
	CanAddOrders      bool
	CanReceiveOrders  bool
	CanQuery          bool
	CanChangeProtocol bool
	CanAddStaff       bool
	CanAddReminders   bool
	IsAdmin           bool
}

type Protocol struct {
	ID                  uuid.UUID
	PNumber             string
	PrimaryInvestigator uuid.UUID
	Title               string
	Allocated           int32
	Balance             int32
	ExpirationDate      time.Time
	IsActive            bool
	PreviousProtocol    uuid.NullUUID
}

type Reminder struct {
	ID             uuid.UUID
	RDate          time.Time
	RCcID          int32
	InvestigatorID uuid.UUID
	Note           string
}

type Setting struct {
	ID               int32
	SettingsComplete bool
	OnlyActivateSelf bool
	TestDataLoaded   bool
}

type Strain struct {
	ID         uuid.UUID
	SName      string
	Vendor     string
	VendorCode string
}
