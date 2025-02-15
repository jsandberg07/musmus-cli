package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddPositionCmd() Command {
	addPositionFlags := make(map[string]Flag)
	addPositionCmd := Command{
		name:        "add",
		description: "Create a new position and set permissions",
		function:    addPositionFunction,
		flags:       addPositionFlags,
		printOrder:  1,
	}

	return addPositionCmd
}

// TODO: extremely minor, but maybe add a check to see if ANY changes have been made and just discard if they havent
func getPositionFlags() map[string]Flag {
	PositionFlags := make(map[string]Flag)

	tFlag := Flag{
		symbol:      "-t",
		description: "Changes the title",
		takesValue:  true,
		printOrder:  1,
	}
	PositionFlags[tFlag.symbol] = tFlag

	aFlag := Flag{
		symbol:      "-a",
		description: "Toggles if the role has permission to add or activate cage cards",
		takesValue:  false,
		printOrder:  2,
	}
	PositionFlags[aFlag.symbol] = aFlag

	dFlag := Flag{
		symbol:      "-d",
		description: "Toggles if the role has permission to deactivate cage cards",
		takesValue:  false,
		printOrder:  3,
	}
	PositionFlags[dFlag.symbol] = dFlag

	oFlag := Flag{
		symbol:      "-o",
		description: "Toggles if the role has permission to add orders",
		takesValue:  false,
		printOrder:  4,
	}
	PositionFlags[oFlag.symbol] = oFlag

	rFlag := Flag{
		symbol:      "-r",
		description: "Toggles if the role has permission to mark orders as received",
		takesValue:  false,
		printOrder:  4,
	}
	PositionFlags[rFlag.symbol] = rFlag

	qFlag := Flag{
		symbol:      "-q",
		description: "Toggles if the role has permission to run queries",
		takesValue:  false,
		printOrder:  5,
	}
	PositionFlags[qFlag.symbol] = qFlag

	pFlag := Flag{
		symbol:      "-p",
		description: "Toggles if the role has permission to make changes to protocols including adding staff to them",
		takesValue:  false,
		printOrder:  6,
	}
	PositionFlags[pFlag.symbol] = pFlag

	sFlag := Flag{
		symbol:      "-s",
		description: "Toggles if the role has permission to add or remove staff or positions",
		takesValue:  false,
		printOrder:  7,
	}
	PositionFlags[sFlag.symbol] = sFlag

	mFlag := Flag{
		symbol:      "-m",
		description: "Toggles if the role has permission to add reminders",
		takesValue:  false,
		printOrder:  8,
	}
	PositionFlags[mFlag.symbol] = mFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags for the command",
		takesValue:  false,
		printOrder:  100,
	}
	PositionFlags[helpFlag.symbol] = helpFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Display WIP permissions current settings",
		takesValue:  false,
		printOrder:  99,
	}
	PositionFlags[printFlag.symbol] = printFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the new position and exits",
		takesValue:  false,
		printOrder:  99,
	}
	PositionFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	PositionFlags[exitFlag.symbol] = exitFlag

	return PositionFlags

}

func addPositionFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionStaff)
	if err != nil {
		return err
	}
	// get title before anything else, or exit early
	title, err := getStringPrompt(cfg, "Please enter the title for the new position,", checkIfPositionTitleUnique)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	// get flags
	flags := getPositionFlags()

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// set defaults
	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}
	cpParams := database.CreatePositionParams{
		Title:             title,
		CanActivate:       false,
		CanDeactivate:     false,
		CanAddOrders:      false,
		CanReceiveOrders:  false,
		CanQuery:          false,
		CanChangeProtocol: false,
		CanAddStaff:       false,
		CanAddReminders:   false,
	}

	fmt.Println("Use flags to toggle permission. All permissions default to false. Multiple flags can be passed in at once.\nUse help to see flags and what permissions they toggle")

	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here
		if reviewed.ChangesMade {
			reviewed.Printed = false
		}

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// [a]ctivate, [d]eact, add [o]rders, [q]uery, [p]rotocol, [s]taff
		// print, save, exit
		for _, arg := range args {
			switch arg.flag {
			case "-a":
				cpParams.CanActivate = !cpParams.CanActivate
				reviewed.ChangesMade = true
			case "-d":
				cpParams.CanDeactivate = !cpParams.CanDeactivate
				reviewed.ChangesMade = true
			case "-o":
				cpParams.CanAddOrders = !cpParams.CanAddOrders
				reviewed.ChangesMade = true
			case "-q":
				cpParams.CanQuery = !cpParams.CanQuery
				reviewed.ChangesMade = true
			case "-p":
				cpParams.CanChangeProtocol = !cpParams.CanChangeProtocol
				reviewed.ChangesMade = true
			case "-s":
				cpParams.CanAddStaff = !cpParams.CanAddStaff
				reviewed.ChangesMade = true
			case "-m":
				cpParams.CanAddReminders = !cpParams.CanAddReminders
				reviewed.ChangesMade = true
			case "-r":
				cpParams.CanReceiveOrders = !cpParams.CanReceiveOrders
				reviewed.ChangesMade = true

			case "-t":
				err := checkIfPositionTitleUnique(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				cpParams.Title = arg.value
				fmt.Printf("New title for new position set: %s\n", cpParams.Title)

			case "help":
				cmdHelp(flags)

			case "print":
				fmt.Println("Printing...")
				printPermissions(&cpParams, nil)
				reviewed.ChangesMade = false
				reviewed.Printed = true
			case "save":
				fmt.Println("Saving and exiting...")
				if !reviewed.Printed {
					fmt.Println("Creating a role with these permissions:")
					printPermissions(&cpParams, nil)
				}
				newPosition, err := cfg.db.CreatePosition(context.Background(), cpParams)
				if err != nil {
					fmt.Println("Error creating new position")
					return err
				}
				if verbose {
					fmt.Println(newPosition)
				}
				exit = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true
			default:
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func checkIfPositionTitleUnique(cfg *Config, input string) error {
	_, err := cfg.db.GetPositionByTitle(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// no position by that name was found
		return nil
	}
	if err != nil {
		// any other DB error so exit
		fmt.Printf("Error checking database for title: %s\n", err)
		os.Exit(1)
	}
	return errors.New("position titles must be unique. Please try again")
}

func getEditPositionCmd() Command {
	editPositionFlags := make(map[string]Flag)
	editPositionCmd := Command{
		name:        "edit",
		description: "Edit an existing position and set permissions",
		function:    editPositionFunction,
		flags:       editPositionFlags,
		printOrder:  2,
	}

	return editPositionCmd
}

// TODO: print all titles so people know what the names are
// flags are in addPosition
func editPositionFunction(cfg *Config) error {
	// permission check
	err := checkPermission(cfg.loggedInPosition, PermissionStaff)
	if err != nil {
		return err
	}
	position, err := getStructPrompt(cfg, "Enter the title of the position to edit,", getPositionStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	// get flags
	flags := getPositionFlags()

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// set defaults
	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}
	upParams := database.UpdatePositionParams{
		Title:             position.Title,
		CanActivate:       position.CanActivate,
		CanDeactivate:     position.CanDeactivate,
		CanAddOrders:      position.CanAddOrders,
		CanQuery:          position.CanQuery,
		CanChangeProtocol: position.CanChangeProtocol,
		CanAddStaff:       position.CanAddStaff,
		ID:                position.ID,
	}

	fmt.Println("Use flags to toggle permission. All permissions default to false. Multiple flags can be passed in at once.\nUse help to see flags and what permissions they toggle")

	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here
		if reviewed.ChangesMade {
			reviewed.Printed = false
		}

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// [a]ctivate, [d]eact, add [o]rders, [q]uery, [p]rotocol, [s]taff
		// print, save, exit
		for _, arg := range args {
			switch arg.flag {
			case "-a":
				upParams.CanActivate = !upParams.CanActivate
				reviewed.ChangesMade = true
			case "-d":
				upParams.CanDeactivate = !upParams.CanDeactivate
				reviewed.ChangesMade = true
			case "-o":
				upParams.CanAddOrders = !upParams.CanAddOrders
				reviewed.ChangesMade = true
			case "-q":
				upParams.CanQuery = !upParams.CanQuery
				reviewed.ChangesMade = true
			case "-p":
				upParams.CanChangeProtocol = !upParams.CanChangeProtocol
				reviewed.ChangesMade = true
			case "-s":
				upParams.CanAddStaff = !upParams.CanAddStaff
				reviewed.ChangesMade = true
			case "-m":
				upParams.CanAddReminders = !upParams.CanAddReminders
				reviewed.ChangesMade = true
			case "-r":
				upParams.CanReceiveOrders = !upParams.CanReceiveOrders
				reviewed.ChangesMade = true

			case "-t":
				err := checkIfPositionTitleUnique(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				upParams.Title = arg.value
				fmt.Printf("Position title set: %s\n", upParams.Title)

			case "help":
				cmdHelp(flags)

			case "print":
				fmt.Println("Printing...")
				printPermissions(nil, &upParams)
				reviewed.ChangesMade = false
				reviewed.Printed = true

			case "save":
				fmt.Println("Saving and exiting...")
				if !reviewed.Printed {
					fmt.Println("Creating a role with these permissions:")
					printPermissions(nil, &upParams)
				}
				err := cfg.db.UpdatePosition(context.Background(), upParams)
				if err != nil {
					fmt.Println("Error updating position")
					return err
				}
				exit = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true
			default:
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func getPositionStruct(cfg *Config, input string) (database.Position, error) {
	position, err := cfg.db.GetPositionByTitle(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Position{}, errors.New("position not found. Please try again")
	}
	if err != nil {
		// any other error
		fmt.Println("Error getting strain from db.")
		return database.Position{}, err
	}

	return position, nil
}

/*
the prospective printing method
func (c createposition) getPrint() printing struct { return p(c) } as a member
func (u updateposition) getPrint() printing struct { return conversion-UUID }
BUT you can't extend functions from another package + they're updating so it'll probably get overwritten on generate
*/

// don't pass in both, just one, and it'll convert and print it
// extremely hacky way to print both structs as theyre not identical in places that aren't printed anyway.
// can't extend structs from another package
func printPermissions(c *database.CreatePositionParams, u *database.UpdatePositionParams) {
	if c == nil && u == nil {
		fmt.Println("Error printing permissions: both params nil")
		return
	}
	if c != nil && u != nil {
		fmt.Println("Error printing permissions: both params NOT nil")
	}

	type PrintPosition struct {
		Title             string
		CanActivate       bool
		CanDeactivate     bool
		CanAddOrders      bool
		CanReceiveOrders  bool
		CanQuery          bool
		CanChangeProtocol bool
		CanAddStaff       bool
		CanAddReminders   bool
	}
	var p PrintPosition

	if c != nil {
		p = PrintPosition(*c)
	} else {
		p.Title = u.Title
		p.CanActivate = u.CanActivate
		p.CanDeactivate = u.CanActivate
		p.CanAddOrders = u.CanAddOrders
		p.CanReceiveOrders = u.CanReceiveOrders
		p.CanQuery = u.CanQuery
		p.CanChangeProtocol = u.CanChangeProtocol
		p.CanAddStaff = u.CanAddStaff
		p.CanAddReminders = u.CanAddReminders
	}

	granted := []string{}
	denied := []string{}
	as := "Activate cage cards"
	if p.CanActivate {
		granted = append(granted, as)
	} else {
		denied = append(denied, as)
	}
	ds := "Deactivate cage cards"
	if p.CanDeactivate {
		granted = append(granted, ds)
	} else {
		denied = append(denied, ds)
	}
	os := "Add orders"
	if p.CanAddOrders {
		granted = append(granted, os)
	} else {
		denied = append(denied, os)
	}
	rs := "Recieve orders"
	if p.CanReceiveOrders {
		granted = append(granted, rs)
	} else {
		denied = append(denied, rs)
	}
	qs := "Run queries"
	if p.CanQuery {
		granted = append(granted, qs)
	} else {
		denied = append(denied, qs)
	}
	ps := "Adjust protocols"
	if p.CanChangeProtocol {
		granted = append(granted, ps)
	} else {
		denied = append(denied, ps)
	}
	ss := "Make changes to staff"
	if p.CanAddStaff {
		granted = append(granted, ss)
	} else {
		denied = append(denied, ss)
	}
	ms := "Add reminders"
	if p.CanAddReminders {
		granted = append(granted, ms)
	} else {
		denied = append(denied, ms)
	}
	fmt.Printf("* %s\n", p.Title)
	if len(granted) == 0 {
		fmt.Println("No permissions granted.")
		return
	}
	if len(denied) == 0 {
		fmt.Println("All permissions granted.")
		return
	}
	fmt.Println("* Allowed permissions:")
	for _, perm := range granted {
		fmt.Println(perm)
	}
	fmt.Println("* Denied permissions:")
	for _, den := range denied {
		fmt.Println(den)
	}
}

// pass in currently logged in user's position from cfg (it's literally stored there)
func checkPermission(i *database.Position, p Permission) error {
	// we dont need to contact the db, the position is already loaded in cfg but check if nil
	if i == nil {
		return errors.New("could not get position to verify permissions")
	}
	var pMsg string
	permitted := true

	switch p {
	case PermissionActivateInactivate:
		if !i.CanActivate {
			permitted = false
			pMsg = "add, activate or inactivate cage cards"
		}
	case PermissionDeactivateReactivate:
		if !i.CanActivate {
			permitted = false
			pMsg = "deactivate or reactivate cage cards"
		}
	case PermissionAddOrder:
		if !i.CanActivate {
			permitted = false
			pMsg = "add orders"
		}
	case PermissionReceiveOrder:
		if !i.CanActivate {
			permitted = false
			pMsg = "receive orders"
		}
	case PermissionRunQueries:
		if !i.CanActivate {
			permitted = false
			pMsg = "run queries"
		}
	case PermissionProtocol:
		if !i.CanActivate {
			permitted = false
			pMsg = "adjust protocols"
		}
	case PermissionStaff:
		if !i.CanActivate {
			permitted = false
			pMsg = "edit staff"
		}
	case PermissionReminders:
		if !i.CanActivate {
			permitted = false
			pMsg = "add reminders"
		}
	default:
		return errors.New("default in check permissions. unknown permission")
	}
	if !permitted {
		msg := fmt.Sprintf("position %s is not permitted to %s", i.Title, pMsg)
		return errors.New(msg)
	} else {
		return nil
	}
}

/* removed as part of refactor
func getPositionToEdit(cfg *Config) (database.Position, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("* Please enter the title of the position to edit, or exit to cancel")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error getting input: %s", err)
			os.Exit(1)
		}

		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input entered. Title must be at least one character.")
			continue
		}

		if input == "exit" || input == "cancel" {
			fmt.Println("Exiting...")
			return database.Position{}, nil
		}

		position, err := cfg.db.GetPositionByTitle(context.Background(), input)
		if err != nil && err.Error() == "sql: no rows in result set" {
			// no position by that name was found
			fmt.Println("No position by that title was found. Please try again.")
			continue
		}
		if err != nil {
			// any other DB error so exit
			fmt.Printf("Error checking database for title: %s\n", err)
			os.Exit(1)
		}

		return position, nil
	}
}
*/

/* removed in refactor
func getNewPositionTitle(cfg *Config) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("* Please enter the title for the new position, or exit to cancel")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error getting new title: %s", err)
			os.Exit(1)
		}

		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input entered. Title must be at least one character.")
			continue
		}

		if input == "exit" || input == "cancel" {
			fmt.Println("Exiting...")
			return "", nil
		}

		name := checkIfPositionTitleUnique(cfg, input)
		if name != "" {
			return name, nil
		}
	}

}
*/

/* removed by combining both print functions into one really bad one
// TODO: struct is the same as update position.
// isn't there a way to create a func that works with 2 types?
func printUpdatePermissions(cp *database.UpdatePositionParams) error {
	granted := []string{}
	denied := []string{}
	as := "Activate cage cards"
	if cp.CanActivate {
		granted = append(granted, as)
	} else {
		denied = append(denied, as)
	}
	ds := "Deactivate cage cards"
	if cp.CanDeactivate {
		granted = append(granted, ds)
	} else {
		denied = append(denied, ds)
	}
	os := "Add and recieve orders"
	if cp.CanAddOrders {
		granted = append(granted, os)
	} else {
		denied = append(denied, os)
	}
	qs := "Run queries"
	if cp.CanQuery {
		granted = append(granted, qs)
	} else {
		denied = append(denied, qs)
	}
	ps := "Adjust protocols"
	if cp.CanChangeProtocol {
		granted = append(granted, ps)
	} else {
		denied = append(denied, ps)
	}
	ss := "Make changes to staff"
	if cp.CanAddStaff {
		granted = append(granted, ss)
	} else {
		denied = append(denied, ss)
	}
	fmt.Printf("* %s\n", cp.Title)
	if len(granted) == 0 {
		fmt.Println("No permissions granted.")
		return nil
	}
	if len(denied) == 0 {
		fmt.Println("All permissions granted.")
		return nil
	}
	fmt.Println("* Allowed permissions:")
	for _, perm := range granted {
		fmt.Println(perm)
	}
	fmt.Println("* Denied permissions:")
	for _, den := range denied {
		fmt.Println(den)
	}

	return nil
}

// TODO: there are two print permission functions that are identical, but the strcuts differ in one stores the UUID.
// can't convert the structs so easily because of it. Has to be a way to DRY this up using interfaces then.
// uuh extremely hacky have func print(*struct A, *struct B) {if struct A == nil, printable = B, else printable = B, then print printable}
// given the weirdness of how theyre printing, this works and reuses what code i have now
func printCreatePermissions(cp *database.CreatePositionParams) error {
	granted := []string{}
	denied := []string{}
	as := "Activate cage cards"
	if cp.CanActivate {
		granted = append(granted, as)
	} else {
		denied = append(denied, as)
	}
	ds := "Deactivate cage cards"
	if cp.CanDeactivate {
		granted = append(granted, ds)
	} else {
		denied = append(denied, ds)
	}
	os := "Add and recieve orders"
	if cp.CanAddOrders {
		granted = append(granted, os)
	} else {
		denied = append(denied, os)
	}
	qs := "Run queries"
	if cp.CanQuery {
		granted = append(granted, qs)
	} else {
		denied = append(denied, qs)
	}
	ps := "Adjust protocols"
	if cp.CanChangeProtocol {
		granted = append(granted, ps)
	} else {
		denied = append(denied, ps)
	}
	ss := "Make changes to staff"
	if cp.CanAddStaff {
		granted = append(granted, ss)
	} else {
		denied = append(denied, ss)
	}
	fmt.Printf("* %s\n", cp.Title)
	if len(granted) == 0 {
		fmt.Println("No permissions granted.")
		return nil
	}
	if len(denied) == 0 {
		fmt.Println("All permissions granted.")
		return nil
	}
	fmt.Println("* Allowed permissions:")
	for _, perm := range granted {
		fmt.Println(perm)
	}
	fmt.Println("* Denied permissions:")
	for _, den := range denied {
		fmt.Println(den)
	}

	return nil
}
*/
