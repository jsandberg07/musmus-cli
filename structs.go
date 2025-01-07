package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

type Flag struct {
	symbol      string
	description string
	takesValue  bool
}

type Argument struct {
	flag  string
	value string
}

type Command struct {
	name        string
	description string
	function    func(cfg *Config, args []Argument) error
	flags       map[string]Flag
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

type Reviewed struct {
	Printed     bool
	ChangesMade bool
}

type ccError struct {
	CCid int
	Err  string
}

type addAllotment struct {
	protocolID uuid.UUID
	total      int
}

type CageCardExport struct {
	CcID          int32
	IName         string
	PNumber       string
	SName         sql.NullString
	ActivatedOn   sql.NullTime
	DeactivatedOn sql.NullTime
}

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
