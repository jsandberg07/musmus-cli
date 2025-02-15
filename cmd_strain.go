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
		printOrder:  1,
	}

	return addStrainCmd
}

func getAddStrainFlags() map[string]Flag {
	addStrainFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the entered strain",
		takesValue:  false,
		printOrder:  99,
	}
	addStrainFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	addStrainFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
		printOrder:  100,
	}
	addStrainFlags[helpFlag.symbol] = helpFlag

	return addStrainFlags
}

func addStrainFunction(cfg *Config) error {
	// no permission check, strains are generally for reference
	name, err := getStringPrompt(cfg, "Enter strain name", checkFuncNil)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	vendor, err := getStringPrompt(cfg, "Enter vendor", checkFuncNil)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	code, err := getStringPrompt(cfg, "Enter strain code", checkIfStrainCodeUnique)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	asParams := database.AddStrainParams{
		SName:      name,
		Vendor:     vendor,
		VendorCode: code,
	}

	flags := getAddStrainFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)

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
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

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
		printOrder:  2,
	}

	return editStrainCmd
}

func getEditStrainFlags() map[string]Flag {
	editStrainFlags := make(map[string]Flag)
	nFlag := Flag{
		symbol:      "-n",
		description: "Set name of strain",
		takesValue:  true,
		printOrder:  1,
	}
	editStrainFlags[nFlag.symbol] = nFlag

	vFlag := Flag{
		symbol:      "-v",
		description: "Sets vendor",
		takesValue:  true,
		printOrder:  2,
	}
	editStrainFlags[vFlag.symbol] = vFlag

	cFlag := Flag{
		symbol:      "-c",
		description: "Sets strain code",
		takesValue:  true,
		printOrder:  3,
	}
	editStrainFlags[cFlag.symbol] = cFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
		printOrder:  100,
	}
	editStrainFlags[helpFlag.symbol] = helpFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves current changes and exits",
		takesValue:  false,
		printOrder:  99,
	}
	editStrainFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	editStrainFlags[exitFlag.symbol] = exitFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current strain draft for review",
		takesValue:  false,
		printOrder:  99,
	}
	editStrainFlags[printFlag.symbol] = printFlag

	return editStrainFlags
}

func editStrainFunction(cfg *Config) error {
	// no permission check, strains are generally for reference
	strain, err := getStructPrompt(cfg, "Enter strain name or ID to edit", getStrainStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	flags := getEditStrainFlags()

	exit := false

	usParams := database.UpdateStrainParams(strain)

	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Use flags to set new values for strain.")
	fmt.Println("Enter 'help' for list of available flags.")

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

		if reviewed.ChangesMade {
			reviewed.Printed = false
		}

		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

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
				if reviewed.Printed {
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
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
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
		return database.Strain{}, errors.New("strain not found. Please try again")
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
