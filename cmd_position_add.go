package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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
	// [a]ctivate, [d]eact, add [o]rders, [q]uery, [p]rotocol, [s]taff, [n]ame
	// START with setting a name
	// and previewing
	// check to see if previewed before saving, print anyway
	// and saving
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
		symbol:      "review",
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
	// get name before anything else, or exit early
	name, err := getNewPositionTitle(cfg)
	if err != nil {
		return err
	}
	if name == "" {
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
		Title:             name,
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
				newName := checkIfPositionTitleUnique(cfg, arg.value)
				if newName != "" {
					cpParams.Title = newName
					fmt.Printf("New title for new position set: %s\n", cpParams.Title)
				}
			case "help":
				err := cmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}
			case "review":
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

func checkIfPositionTitleUnique(cfg *Config, input string) string {
	_, err := cfg.db.GetPositionByTitle(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// no position by that name was found
		return input
	}
	if err != nil {
		// any other DB error so exit
		fmt.Printf("Error checking database for title: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Position with that name was found. Please try again.")
	return ""
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
