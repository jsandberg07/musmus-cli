package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getEditInvestigatorCmd() Command {
	editInvestigatorFlags := make(map[string]Flag)
	editInvestigatorCmd := Command{
		name:        "edit",
		description: "Used for editing existing investigators",
		function:    editInvestigatorFunction,
		flags:       editInvestigatorFlags,
	}

	return editInvestigatorCmd
}

// print, exit, save, help, [i]nvestigator name, [n]ickname, [p]osition, [a]ctive, [e]mail
// TODO: a function that parses things with QUOTES and a space properly JK JUST USE UNDERSCORES
func getEditInvestigatorFlags() map[string]Flag {
	editInvestigatorFlags := make(map[string]Flag)
	iFlag := Flag{
		symbol:      "i",
		description: "Changes proper 'investigator' name",
		takesValue:  true,
	}
	editInvestigatorFlags["-"+iFlag.symbol] = iFlag

	nFlag := Flag{
		symbol:      "n",
		description: "Changes nickname. Enter 'delete' to remove nickname",
		takesValue:  true,
	}
	editInvestigatorFlags["-"+nFlag.symbol] = nFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Changes position",
		takesValue:  true,
	}
	editInvestigatorFlags["-"+pFlag.symbol] = pFlag

	aFlag := Flag{
		symbol:      "a",
		description: "Toggles if investigator is active on the protocols",
		takesValue:  false,
	}
	editInvestigatorFlags["-"+aFlag.symbol] = aFlag

	eFlag := Flag{
		symbol:      "e",
		description: "Changes email. Enter 'delete' to remove email",
		takesValue:  true,
	}
	editInvestigatorFlags["-"+eFlag.symbol] = eFlag

	// ect as needed or remove the "-"+ for longer ones

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current data for review",
		takesValue:  false,
	}
	editInvestigatorFlags[printFlag.symbol] = printFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
	}
	editInvestigatorFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	editInvestigatorFlags[exitFlag.symbol] = exitFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves changes made",
		takesValue:  false,
	}
	editInvestigatorFlags[saveFlag.symbol] = saveFlag

	return editInvestigatorFlags

}

// ask for a name, then pass in flags for everything
// then print old to new
// look into removing the args thing, might have to stay
func editInvestigatorFunction(cfg *Config, args []Argument) error {
	investigator, err := getStructPrompt(cfg, "Enter the name of the investigator you'd like to edit,", getInvestigatorStruct)
	if err != nil {
		return err
	}
	nilInvestigator := database.Investigator{}
	if investigator == nilInvestigator {
		fmt.Println("Exiting...")
		return nil
	}

	position, err := cfg.db.GetUserPosition(context.Background(), investigator.Position)
	if err != nil {
		fmt.Printf("Error getting position for investigator: %s\n", err)
		os.Exit(1)
	}

	// get flags
	flags := getEditInvestigatorFlags()

	// set defaults
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

	// the reader
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Making changes to %s\n", investigator.IName)
	fmt.Println("If entering a value that uses spaces (like between a first and last name) use an underscore")
	fmt.Println("It will be changed to a space after")

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

		// print, exit, save, help, [i]nvestigator name, [n]ickname, [p]osition, [a]ctive, [e]mail
		for _, arg := range args {
			switch arg.flag {
			case "print":
				// need some generic type for printing
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
					return err
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
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
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

func getPositionByFlag(cfg *Config, title string) (database.Position, error) {
	position, err := cfg.db.GetPositionByTitle(context.Background(), title)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("Error getting position from DB")
		return database.Position{ID: uuid.Nil}, err
	}
	if err.Error() == "sql: no rows in result set" {
		fmt.Println("No position by that title found")
		return database.Position{ID: uuid.Nil}, err
	}

	return position, nil

}

func checkIfInvestigatorNameUnique(cfg *Config, name string) error {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("Error getting name from DB")
		return err
	}
	if len(investigators) != 0 {
		fmt.Println("Investigator name is not unique. Please consider adding a nickname to both investigators.")
	}
	return nil
}

func getInvestigatorStruct(cfg *Config, input string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	// TODO: does returning many even throw the "no rows in result set" error?
	if len(investigators) == 0 {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	if len(investigators) > 1 {
		return database.Investigator{}, errors.New("vague investigator name. Please try again")
	}
	if err != nil {
		// any other error
		return database.Investigator{}, err
	}

	return investigators[0], nil
}

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
				cmdHelp(flags)

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

/* removed from refactor


func getInvestigatorByName(cfg *Config) (database.Investigator, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the name or nickname of the investigator you would like to edit, or exit to cancel")
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
			return database.Investigator{ID: uuid.Nil}, nil
		}

		investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
		if err != nil && err.Error() != "sql: no rows in result set" {
			// error that isnt related to no rows returned
			fmt.Println("Error checking db for investigator")
			return database.Investigator{ID: uuid.Nil}, err
		}
		if len(investigators) == 0 {
			fmt.Println("No investigator by that name or nickname found. Please try again.")
			continue
		}
		if len(investigators) > 1 {
			fmt.Println("Vague investigator name. Consider entering a nickname instead.")
			continue
		}
		return investigators[0], nil
	}
}



*/

/* removed by refactor???
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
*/
