package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jsandberg07/clitest/internal/database"
)

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

// flags are in addPosition
func editPositionFunction(cfg *Config, args []Argument) error {
	// get positon to edit before anything else, or exit early
	position, err := getPositionToEdit(cfg)
	if err != nil {
		return err
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
				newName := checkIfPositionTitleUnique(cfg, arg.value)
				if newName != "" {
					upParams.Title = newName
					fmt.Printf("Position title set: %s\n", upParams.Title)
				}
			case "help":
				err := cmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}
			case "review":
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

// TODO: print all titles so people know what the names are
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
