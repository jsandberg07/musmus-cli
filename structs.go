package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

// TODO: replacing print order with a float would unironcally make things more extendable because you spend less time
// adding 1 to each number you already have if you just do like 1.1, 1.11, ect. That's dumb but i love it.
type Flag struct {
	symbol      string
	description string
	takesValue  bool
	printOrder  int
}

type Argument struct {
	flag  string
	value string
}

type PositionParams struct {
	Title             string
	CanActivate       bool
	CanDeactivate     bool
	CanAddOrders      bool
	CanQuery          bool
	CanChangeProtocol bool
	CanAddStaff       bool
}

type Command struct {
	name        string
	description string
	function    func(cfg *Config) error
	flags       map[string]Flag
	printOrder  int
}

type Config struct {
	currentState         *State
	nextState            *State
	db                   *database.Queries
	loggedInInvestigator *database.Investigator
	loggedInPosition     *database.Position
}

type State struct {
	currentCommands map[string]Command
	cliMessage      string
}

type CageCard struct {
	CCid   int
	Date   time.Time
	Person string
}

// TODO: generally used the same in every place, add member funcs like "changed made" or "reset" or "printed" or "check"
type Reviewed struct {
	Printed     bool
	ChangesMade bool
}

type ccError struct {
	CCid int
	Err  string
}

type CageCardExport struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
	OrderNumber   sql.NullString
}

// TODO: optimize size (order by largest to smallest size)
type CageCardActivationParams struct {
	ccID         int
	date         time.Time
	allotment    int
	strain       database.Strain
	keepStrain   bool
	note         string
	keepNote     bool
	daysReminder int
	keepReminder bool
}

// sets defaults or 0 values for struct
func (s *CageCardActivationParams) init() {
	s.ccID = 0
	s.date = normalizeDate(time.Now())
	s.allotment = 0
	s.strain = database.Strain{ID: uuid.Nil}
	s.keepStrain = false
	s.note = ""
	s.keepNote = false
	s.daysReminder = 0
	s.keepReminder = false
}

// checks the keep properties and resets the ones that aren't marked as kept.
// Is this a bad name? It's like of close to the variable names so it might be confusing
func (s *CageCardActivationParams) keepCheck() {
	if !s.keepStrain {
		s.strain = database.Strain{ID: uuid.Nil}
	}
	if !s.keepNote {
		s.note = ""
	}
	if !s.keepReminder {
		s.daysReminder = 0
	}
}

type Permission int

const (
	PermissionActivateInactivate = iota
	PermissionDeactivateReactivate
	PermissionAddOrder
	PermissionReceiveOrder
	PermissionRunQueries
	PermissionProtocol
	PermissionStaff
	PermissionReminders
)

/*
Create a flag:
symbol, description, and if it takes a value
symbol is without the -
in the getCmd function, $flag := flag{}
add to commands map
add handling in the function itself. takes value are used later, doesnt sets a bool

Create a command:
write the new function (handle flagss)
create a new function getNewCmd() Command {}
flags map := make(map[string]flag)
newCmd := Command{name, description, function, flags}

Create a state:
getState() &State {map, cli message}
make sure to add help by default, prints what is available in current map
getStateMap() and put that in
*/
