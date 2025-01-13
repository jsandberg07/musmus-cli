package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddStrainCmd() Command {
	addStrainFlags := make(map[string]Flag)
	addStrainCmd := Command{
		name:        "add",
		description: "Used for adding strains to the database",
		function:    addStrainFunction,
		flags:       addStrainFlags,
	}

	return addStrainCmd
}

// via prompts
// so just save and exit and help for consistency
func getAddStrainFlags() map[string]Flag {
	addStrainFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the entered strain",
		takesValue:  false,
	}
	addStrainFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	addStrainFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	addStrainFlags[helpFlag.symbol] = helpFlag

	// ect as needed or remove the "-"+ for longer ones

	return addStrainFlags

}

// look into removing the args thing, might have to stay
func addStrainFunction(cfg *Config, args []Argument) error {

	name, err := getStringPrompt(cfg, "Enter strain name", checkFuncNil)
	if err != nil {
		return err
	}
	if name == "" {
		fmt.Println("Exiting...")
		return nil
	}

	vendor, err := getStringPrompt(cfg, "Enter vendor", checkFuncNil)
	if err != nil {
		return err
	}
	if vendor == "" {
		fmt.Println("Exiting...")
		return nil
	}

	code, err := getStringPrompt(cfg, "Enter strain code", checkIfStrainCodeUnique)
	if err != nil {
		return err
	}
	if code == "" {
		fmt.Println("Exiting...")
		return nil
	}

	asParams := database.AddStrainParams{
		SName:      name,
		Vendor:     vendor,
		VendorCode: code,
	}

	// get flags
	flags := getAddStrainFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
	fmt.Println("Please review the following information")
	fmt.Println("Enter 'save' to save the strain or 'exit' to leave without saving")
	printAddStrain(&asParams)
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

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "help":
				cmdHelp(flags)
			case "exit":
				fmt.Println("Exiting with saving...")
				exit = true
			case "save":
				fmt.Println("Saving...")
				strain, err := cfg.db.AddStrain(context.Background(), asParams)
				if err != nil {
					fmt.Println("Error saving strain")
					return err
				}
				if verbose {
					fmt.Println(strain)
				}
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

func checkIfStrainCodeUnique(cfg *Config, s string) error {
	_, err := cfg.db.GetStrainByCode(context.Background(), s)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error retrieving data from the DB")
		return err
	}
	if err == nil {
		// strain found, meaning input is not unique
		return errors.New("Strain by that ID already exists. Please try again.")
	}

	// strain is unique
	return nil

}

func printAddStrain(as *database.AddStrainParams) {
	fmt.Printf("Name: %s\n", as.SName)
	fmt.Printf("Vendor: %s\n", as.Vendor)
	fmt.Printf("Code: %s\n", as.VendorCode)
}

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
