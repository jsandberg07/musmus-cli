package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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

// idea for how i should have created more reusable functions for all the other data types
// more generic get string with a prompt, instead of separate functions for everything
// just pass in "get investigator name to edit" instead of new function just to say "to edit"
// write the same program 3 times and you'll realize what you want you want to refactor
func getStringPrompt(cfg *Config, prompt string, checkFunc func(*Config, string) error) (string, error) {
	fmt.Println(prompt + " or exit to cancel")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			return "", nil
		}

		// then have check if unique or check if not unique after
		err = checkFunc(cfg, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return input, nil

	}
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

func checkFuncNil(cfg *Config, s string) error {
	// look into optional 1st order func params
	return nil
}

func printAddStrain(as *database.AddStrainParams) {
	fmt.Printf("Name: %s\n", as.SName)
	fmt.Printf("Vendor: %s\n", as.Vendor)
	fmt.Printf("Code: %s\n", as.VendorCode)
}
