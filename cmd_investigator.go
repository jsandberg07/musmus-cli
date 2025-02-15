package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/jsandberg07/clitest/internal/database"
)

func getEditInvestigatorCmd() Command {
	editInvestigatorFlags := make(map[string]Flag)
	editInvestigatorCmd := Command{
		name:        "edit",
		description: "Used for editing existing investigators",
		function:    editInvestigatorFunction,
		flags:       editInvestigatorFlags,
		printOrder:  2,
	}

	return editInvestigatorCmd
}

func getEditInvestigatorFlags() map[string]Flag {
	editInvestigatorFlags := make(map[string]Flag)
	iFlag := Flag{
		symbol:      "-i",
		description: "Changes proper 'investigator' name",
		takesValue:  true,
		printOrder:  1,
	}
	editInvestigatorFlags[iFlag.symbol] = iFlag

	nFlag := Flag{
		symbol:      "-n",
		description: "Changes nickname. Enter 'delete' to remove nickname",
		takesValue:  true,
		printOrder:  2,
	}
	editInvestigatorFlags[nFlag.symbol] = nFlag

	pFlag := Flag{
		symbol:      "-p",
		description: "Changes position",
		takesValue:  true,
		printOrder:  3,
	}
	editInvestigatorFlags[pFlag.symbol] = pFlag

	aFlag := Flag{
		symbol:      "-a",
		description: "Toggles if investigator is active on the protocols",
		takesValue:  false,
		printOrder:  4,
	}
	editInvestigatorFlags[aFlag.symbol] = aFlag

	eFlag := Flag{
		symbol:      "-e",
		description: "Changes email. Enter 'delete' to remove email",
		takesValue:  true,
		printOrder:  5,
	}
	editInvestigatorFlags[eFlag.symbol] = eFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current data for review",
		takesValue:  false,
		printOrder:  100,
	}
	editInvestigatorFlags[printFlag.symbol] = printFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
		printOrder:  100,
	}
	editInvestigatorFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	editInvestigatorFlags[exitFlag.symbol] = exitFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves changes made",
		takesValue:  false,
		printOrder:  99,
	}
	editInvestigatorFlags[saveFlag.symbol] = saveFlag

	return editInvestigatorFlags

}

func editInvestigatorFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionStaff)
	if err != nil {
		return err
	}
	investigator, err := getStructPrompt(cfg, "Enter the name of the investigator you'd like to edit,", getInvestigatorStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	position, err := cfg.db.GetUserPosition(context.Background(), investigator.Position)
	if err != nil {
		fmt.Printf("Error getting position for investigator: %s\n", err)
		os.Exit(1)
	}

	flags := getEditInvestigatorFlags()

	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}
	uiParam := database.UpdateInvestigatorParams{
		ID:       investigator.ID,
		IName:    investigator.IName,
		Nickname: investigator.Nickname,
		Email:    investigator.Email,
		Position: position.ID,
		Active:   investigator.Active,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Making changes to %s\n", investigator.IName)
	fmt.Println("If entering a value that uses spaces (like between a first and last name) use an underscore")
	fmt.Println("It will be changed to a space after")

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
			case "print":
				printEditInvestigator(&uiParam, &position)
				reviewed.ChangesMade = false
				reviewed.Printed = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true

			case "save":
				if reviewed.Printed {
					fmt.Println("Saving...")
				} else {
					fmt.Println("Saving with the following information")
					printEditInvestigator(&uiParam, &position)
				}

				err := cfg.db.UpdateInvestigator(context.Background(), uiParam)
				if err != nil {
					fmt.Println("Error updating investigator in DB")
					return err
				}
				fmt.Println("Investigator has been added. They will be asked to update their password on first login")
				exit = true

			case "help":
				cmdHelp(flags)

			case "-i":
				err := checkIfInvestigatorNameUnique(cfg, arg.value)
				if err != nil {
					return err
				}
				uiParam.IName = arg.value
				reviewed.ChangesMade = true

			case "-n":
				if arg.value == "delete" {
					uiParam.Nickname = sql.NullString{Valid: false}
				} else {
					uiParam.Nickname = sql.NullString{Valid: true, String: arg.value}
				}
				reviewed.ChangesMade = true

			case "-p":
				position, err = getPositionByFlag(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					continue
				}
				uiParam.Position = position.ID
				reviewed.ChangesMade = true

			case "-a":
				uiParam.Active = !uiParam.Active
				if uiParam.Active {
					fmt.Println("Investigator is flagged as active")
				} else {
					fmt.Println("Invetstigator is flagged as inactive")
				}
				reviewed.ChangesMade = true

			case "-e":
				if arg.value == "delete" {
					uiParam.Email = sql.NullString{Valid: false}
				} else {
					uiParam.Email = sql.NullString{Valid: true, String: arg.value}
				}
				reviewed.ChangesMade = true

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

func printEditInvestigator(ui *database.UpdateInvestigatorParams, p *database.Position) {
	nullfield := sql.NullString{
		Valid: false,
	}
	fmt.Printf("Name: %s\n", ui.IName)
	fmt.Printf("Position: %s\n", p.Title)
	if ui.Email != nullfield {
		fmt.Printf("Email: %s\n", ui.Email.String)
	}
	if ui.Nickname != nullfield {
		fmt.Printf("Nickname: %s\n", ui.Nickname.String)
	}
}

func getAddInvestigatorCmd() Command {
	addInvestigatorFlags := make(map[string]Flag)
	addInvestigatorCmd := Command{
		name:        "add",
		description: "Used for adding a new investigator",
		function:    addInvestigatorFunction,
		flags:       addInvestigatorFlags,
		printOrder:  1,
	}

	return addInvestigatorCmd
}

func getAddInvestigatorFlags() map[string]Flag {
	addInvestigatorFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the investigator to the database",
		takesValue:  false,
		printOrder:  100,
	}
	addInvestigatorFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	addInvestigatorFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
		printOrder:  100,
	}
	addInvestigatorFlags[helpFlag.symbol] = helpFlag

	return addInvestigatorFlags
}

func addInvestigatorFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionStaff)
	if err != nil {
		return err
	}
	name, err := getStringPrompt(cfg, "Enter name of new investigator", checkIfInvestigatorNameUnique)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	position, err := getStructPrompt(cfg, "Enter position of new investigator", getPositionStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	ciParams := database.CreateInvestigatorParams{
		IName:    name,
		Position: position.ID,
	}
	err = getNewInvestigatorExtraInfo(&ciParams)
	if err != nil {
		return err
	}

	flags := getAddInvestigatorFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("New investigator will be added with the following info:")
	printNewInvestigator(&ciParams, &position)

	fmt.Println("Enter 'save' to add the investigator to the database")
	fmt.Println("Or 'exit' to exit without saving")

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
			case "save":
				fmt.Println("Saving...")
				investigator, err := cfg.db.CreateInvestigator(context.Background(), ciParams)
				if err != nil {
					fmt.Println("Error creating investigator")
					return err
				}
				exit = true
				if verbose {
					fmt.Println(investigator)
				}
			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true
			case "help":
				cmdHelp(flags)

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

// separate because it allows empty inputs
func getNewInvestigatorExtraInfo(ciParams *database.CreateInvestigatorParams) error {
	fmt.Println("The next info isn't required. Skip by entering nothing.")
	fields := []string{"a nickname", "an email"}
	inputs := make([]string, len(fields))
	reader := bufio.NewReader(os.Stdin)

	for i, field := range fields {
		fmt.Printf("Please enter %s\n", field)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		inputs[i] = input
	}

	if inputs[0] == "" {
		ciParams.Nickname = sql.NullString{Valid: false}
	} else {
		ciParams.Nickname = sql.NullString{Valid: true, String: inputs[0]}
	}

	if inputs[1] == "" {
		ciParams.Email = sql.NullString{Valid: false}
	} else {
		ciParams.Email = sql.NullString{Valid: true, String: inputs[1]}
	}

	return nil
}

func printNewInvestigator(ci *database.CreateInvestigatorParams, p *database.Position) {
	nullfield := sql.NullString{
		Valid: false,
	}
	fmt.Printf("Name: %s\n", ci.IName)
	fmt.Printf("Position: %s\n", p.Title)
	if ci.Email != nullfield {
		fmt.Printf("Email: %s\n", ci.Email.String)
	}
	if ci.Nickname != nullfield {
		fmt.Printf("Nickname: %s\n", ci.Nickname.String)
	}
}
