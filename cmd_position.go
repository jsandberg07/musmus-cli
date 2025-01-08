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
	}

	return addPositionCmd
}

// TODO: extremely minor, but maybe add a check to see if ANY changes have been made and just discard if they havent
func getPositionFlags() map[string]Flag {
	PositionFlags := make(map[string]Flag)

	tFlag := Flag{
		symbol:      "t",
		description: "Changes the title",
		takesValue:  true,
	}
	PositionFlags["-"+tFlag.symbol] = tFlag

	aFlag := Flag{
		symbol:      "a",
		description: "Toggles if the role has permission to add or activate cage cards",
		takesValue:  false,
	}
	PositionFlags["-"+aFlag.symbol] = aFlag

	dFlag := Flag{
		symbol:      "d",
		description: "Toggles if the role has permission to deactivate cage cards",
		takesValue:  false,
	}
	PositionFlags["-"+dFlag.symbol] = dFlag

	oFlag := Flag{
		symbol:      "o",
		description: "Toggles if the role has permission to add or mark orders as recieved",
		takesValue:  false,
	}
	PositionFlags["-"+oFlag.symbol] = oFlag

	qFlag := Flag{
		symbol:      "q",
		description: "Toggles if the role has permission to run queries",
		takesValue:  false,
	}
	PositionFlags["-"+qFlag.symbol] = qFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Toggles if the role has permission to make changes to protocols including adding staff to them",
		takesValue:  false,
	}
	PositionFlags["-"+pFlag.symbol] = pFlag

	sFlag := Flag{
		symbol:      "s",
		description: "Toggles if the role has permission to add or remove staff or positions",
		takesValue:  false,
	}
	PositionFlags["-"+sFlag.symbol] = sFlag

	// ect as needed or remove the "-"+ for longer ones

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags for the command",
		takesValue:  false,
	}
	PositionFlags[helpFlag.symbol] = helpFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Display WIP permissions current settings",
		takesValue:  false,
	}
	PositionFlags[printFlag.symbol] = printFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the new position and exits",
		takesValue:  false,
	}
	PositionFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	PositionFlags[exitFlag.symbol] = exitFlag

	return PositionFlags

}

func addPositionFunction(cfg *Config, args []Argument) error {
	// get title before anything else, or exit early
	title, err := getStringPrompt(cfg, "Please enter the title for the new position,", checkIfPositionTitleUnique)
	if err != nil {
		return err
	}
	if title == "" {
		// user said exit, so cancel creating new position
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
		CanQuery:          false,
		CanChangeProtocol: false,
		CanAddStaff:       false,
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
		if reviewed.ChangesMade == true {
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

			case "-t":
				err := checkIfPositionTitleUnique(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				cpParams.Title = arg.value
				fmt.Printf("New title for new position set: %s\n", cpParams.Title)

			case "help":
				err := cmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}
			case "print":
				fmt.Println("Printing...")
				err := printCreatePermissions(&cpParams)
				if err != nil {
					fmt.Println("Error printing permissions")
				}
				reviewed.ChangesMade = false
				reviewed.Printed = true
			case "save":
				fmt.Println("Saving and exiting...")
				if reviewed.Printed == false {
					fmt.Println("Creating a role with these permissions:")
					err := printCreatePermissions(&cpParams)
					if err != nil {
						fmt.Println("Error printing permissions")
					}
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
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
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
	return errors.New("Position titles must be unique. Please try again.")
}

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

func getEditPositionCmd() Command {
	editPositionFlags := make(map[string]Flag)
	editPositionCmd := Command{
		name:        "edit",
		description: "Edit an existing position and set permissions",
		function:    editPositionFunction,
		flags:       editPositionFlags,
	}

	return editPositionCmd
}

// TODO: print all titles so people know what the names are
// flags are in addPosition
func editPositionFunction(cfg *Config, args []Argument) error {
	position, err := getStructPrompt(cfg, "Enter the title of the position to edit,", getPositionStruct)
	if err != nil {
		return err
	}
	nilPosition := database.Position{}
	if position == nilPosition {
		fmt.Println("Exiting...")
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
		if reviewed.ChangesMade == true {
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

			case "-t":
				err := checkIfPositionTitleUnique(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				upParams.Title = arg.value
				fmt.Printf("Position title set: %s\n", upParams.Title)

			case "help":
				err := cmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}

			case "print":
				fmt.Println("Printing...")
				err := printUpdatePermissions(&upParams)
				if err != nil {
					fmt.Println("Error printing permissions")
				}
				reviewed.ChangesMade = false
				reviewed.Printed = true

			case "save":
				fmt.Println("Saving and exiting...")
				if reviewed.Printed == false {
					fmt.Println("Creating a role with these permissions:")
					err := printUpdatePermissions(&upParams)
					if err != nil {
						fmt.Println("Error printing permissions")
					}
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
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

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

func getPositionStruct(cfg *Config, input string) (database.Position, error) {
	position, err := cfg.db.GetPositionByTitle(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Position{}, errors.New("Position not found. Please try again.")
	}
	if err != nil {
		// any other error
		fmt.Println("Error getting strain from db.")
		return database.Position{}, err
	}

	return position, nil
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
