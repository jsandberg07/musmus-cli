package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jsandberg07/clitest/internal/database"
)

// genuinely not sure if anybody would ever use this but it's a non reliant way to test out some refactoring
// we have make sure string is ok for non dupe values
// now get return a data struct via an input
// literally undo several weeks worth of work because you just sit down and type

func getEditStrainCmd() Command {
	editStrainFlags := make(map[string]Flag)
	editStrainCmd := Command{
		name:        "edit",
		description: "Used for editing existing strains",
		function:    editStrainFunction,
		flags:       editStrainFlags,
	}

	return editStrainCmd
}

// save, exit, print, [n]ame, [v]endor, [c]ode
func getEditStrainFlags() map[string]Flag {
	editStrainFlags := make(map[string]Flag)
	nFlag := Flag{
		symbol:      "n",
		description: "Set name of strain",
		takesValue:  true,
	}
	editStrainFlags["-"+nFlag.symbol] = nFlag

	vFlag := Flag{
		symbol:      "v",
		description: "Sets vendor",
		takesValue:  true,
	}
	editStrainFlags["-"+vFlag.symbol] = vFlag

	cFlag := Flag{
		symbol:      "c",
		description: "Sets strain code",
		takesValue:  true,
	}
	editStrainFlags["-"+cFlag.symbol] = cFlag

	// ect as needed or remove the "-"+ for longer ones

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	editStrainFlags[helpFlag.symbol] = helpFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves current changes and exits",
		takesValue:  false,
	}
	editStrainFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	editStrainFlags[exitFlag.symbol] = exitFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current strain draft for review",
		takesValue:  false,
	}
	editStrainFlags[printFlag.symbol] = printFlag

	return editStrainFlags

}

// look into removing the args thing, might have to stay
func editStrainFunction(cfg *Config, args []Argument) error {
	nilStrain := database.Strain{}
	strain, err := getStructPrompt(cfg, "Enter strain name or ID to edit", getStrainStruct)
	if err != nil {
		return err
	}
	if strain == nilStrain {
		fmt.Println("Exiting...")
	}

	// get flags
	flags := getEditStrainFlags()

	// set defaults
	exit := false
	usParams := database.UpdateStrainParams{
		ID:         strain.ID,
		SName:      strain.SName,
		Vendor:     strain.Vendor,
		VendorCode: strain.VendorCode,
	}

	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Use flags to set new values for strain.")
	fmt.Println("Enter 'help' for list of available flags.")

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

		// save, exit, print, [n]ame, [v]endor, [c]ode
		for _, arg := range args {
			switch arg.flag {
			case "-n":
				usParams.SName = arg.value
				reviewed.ChangesMade = true

			case "-v":
				usParams.Vendor = arg.value
				reviewed.ChangesMade = true

			case "-c":
				// if same as initial, will throw an error for duplicate
				if arg.value == strain.VendorCode {
					usParams.VendorCode = arg.value
					break
				}
				err := checkIfStrainCodeUnique(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				usParams.VendorCode = arg.value
				reviewed.ChangesMade = true

			case "help":
				cmdHelp(flags)

			case "save":
				if reviewed.Printed == false {
					fmt.Println("Updating strain with the following info:")
					printUpdateStrain(&usParams)
				}
				err := cfg.db.UpdateStrain(context.Background(), usParams)
				if err != nil {
					fmt.Println("Error saving strain to db")
					return err
				}
				exit = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true

			case "print":
				printUpdateStrain(&usParams)
				reviewed.ChangesMade = false
				reviewed.Printed = true

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

func getStrainStruct(cfg *Config, input string) (database.Strain, error) {
	strain, err := cfg.db.GetStrainByName(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Strain{}, errors.New("Strain not found. Please try again.")
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting strain from DB.")
		return database.Strain{}, err
	}

	return strain, nil
}

func printUpdateStrain(us *database.UpdateStrainParams) {
	fmt.Printf("Name: %s\n", us.SName)
	fmt.Printf("Vendor: %s\n", us.Vendor)
	fmt.Printf("Code: %s\n", us.VendorCode)
}
