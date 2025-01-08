package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getAddInvestigatorCmd() Command {
	addInvestigatorFlags := make(map[string]Flag)
	addInvestigatorCmd := Command{
		name:        "add",
		description: "Used for adding a new investigator",
		function:    addInvestigatorFunction,
		flags:       addInvestigatorFlags,
	}

	return addInvestigatorCmd
}

func getAddInvestigatorFlags() map[string]Flag {
	addInvestigatorFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the investigator to the database",
		takesValue:  false,
	}
	addInvestigatorFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	addInvestigatorFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	addInvestigatorFlags[helpFlag.symbol] = helpFlag

	// ect as needed or remove the "-"+ for longer ones

	return addInvestigatorFlags

}

// prompt step by step cause some of the fields could be null and whatever
// then review and you can change them with flags after
// and always allow you to exit early or whatever
// name, position, email, nickname, assume active
// look into removing the args thing, might have to stay
func addInvestigatorFunction(cfg *Config, args []Argument) error {
	name, err := getStringPrompt(cfg, "Enter name of new investigator", checkIfInvestigatorNameUnique)
	if err != nil {
		return err
	}
	if name == "" {
		fmt.Println("Exiting...")
		return nil
	}

	position, err := getStructPrompt(cfg, "Enter position of new investigator", getPositionStruct)
	if err != nil {
		return err
	}
	nilPosition := database.Position{}
	if position == nilPosition {
		fmt.Println("Exiting...")
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

	// get flags
	flags := getAddInvestigatorFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("New investigator will be added with the following info:")
	printNewInvestigator(&ciParams, &position)

	// TODO: add the ability to use flags to edit values after the fact
	fmt.Println("Enter 'save' to add the investigator to the database")
	fmt.Println("Or 'exit' to exit without saving")

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

		// but normal loop now
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
				err := cmdHelp(flags)
				if err != nil {
					fmt.Printf("Error printing help: %s\n", err)
				}
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

func getNewInvestigatorPosition(cfg *Config) (database.Position, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the position of the new investigator, or exit to cancel")
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
			return database.Position{ID: uuid.Nil}, nil
		}

		position, err := cfg.db.GetPositionByTitle(context.Background(), input)
		if err != nil && err.Error() == "sql: no rows in result set" {
			fmt.Println("Position by that title not found. Please try again")
			continue
		}
		if err != nil {
			fmt.Println("Error getting position from database")
			return database.Position{ID: uuid.Nil}, err
		}

		return position, nil

	}
}

func getNewInvestigatorExtraInfo(ciParams *database.CreateInvestigatorParams) error {
	fmt.Println("The next info isn't required. Skip by entering nothing.")
	fields := []string{"a nickname", "an email"}
	inputs := make([]string, len(fields))
	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
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

/* removed in refactor

func getNewInvestigatorName(cfg *Config) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the name of the investigator, or exit to cancel")
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

		err = checkIfInvestigatorNameUnique(cfg, input)
		if err != nil {
			return "", err
		}

		return input, nil


	}
}
*/
