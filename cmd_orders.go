package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddOrderCmd() Command {
	addOrderFlags := make(map[string]Flag)
	addOrderCmd := Command{
		name:        "add",
		description: "Used for adding orders",
		function:    addOrderFunction,
		flags:       addOrderFlags,
	}

	return addOrderCmd
}

// prompts so just save, exit, print
func getAddOrderFlags() map[string]Flag {
	addOrderFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the current order",
		takesValue:  false,
	}
	addOrderFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	addOrderFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints flags available for current command",
		takesValue:  false,
	}
	addOrderFlags[helpFlag.symbol] = helpFlag

	return addOrderFlags

}

// look into removing the args thing, might have to stay
func addOrderFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getAddOrderFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	orderNumber, err := getStringPrompt(cfg, "Enter order number", checkIfOrderNumberUnique)
	if err != nil {
		return err
	}
	if orderNumber == "" {
		fmt.Println("Exiting...")
	}

	date, err := getDatePrompt("Enter expected date")
	if err != nil {
		return err
	}
	nilDate := time.Time{}
	if date == nilDate {
		fmt.Println("Exiting...")
		return nil
	}

	investigator, err := getStructPrompt(cfg, "Enter investigator receiving order", getInvestigatorStruct)
	if err != nil {
		return err
	}
	nilInvestigator := database.Investigator{}
	if investigator == nilInvestigator {
		fmt.Println("Exiting...")
		return nil
	}

	strain, err := getStructPrompt(cfg, "Enter strain of order", getStrainStruct)
	if err != nil {
		return err
	}
	nilStrain := database.Strain{}
	if strain == nilStrain {
		fmt.Println("Exiting...")
		return nil
	}

	note, err := getStringPrompt(cfg, "Optionally enter a note. Will be applied to all cage cards from order", checkFuncNil)
	if err != nil {
		return err
	}

	// working here: create the params, set the flags, remember to check if note valid or not

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
			case "-X":
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

func checkIfOrderNumberUnique(cfg *Config, input string) error {
	_, err := cfg.db.GetOrderByNumber(context.Background(), input)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}

	if err == nil {
		// is not unique
		return errors.New("order number is not unique. Please try again")
	}

	// is unique
	return nil

}

func printNewOrder(o *database.CreateNewOrderParams, i *database.Investigator, s *database.Strain) {
	fmt.Printf("* Number - %s\n", o.OrderNumber)
	fmt.Println("* Date - %v\n", o.ExpectedDate)
	fmt.Println("* Investigator - %s\n", i.IName)
	fmt.Println("* Strain - %s\n", s.SName)
	if o.Note.Valid {
		fmt.Printf("* Note - %s", o.Note.String)
	}
}
