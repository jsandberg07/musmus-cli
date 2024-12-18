package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getCCDeactivationCmd() Command {
	deactivationFlags := make(map[string]Flag)
	deactivationCmd := Command{
		name:        "deactivate",
		description: "Used for deactivating cage cards",
		function:    deactivateFunction,
		flags:       deactivationFlags,
	}

	return deactivationCmd
}

// cc, pop, process, exit, list of errors (both previously deact and not activated)
func getDeactivationFlags() map[string]Flag {
	deactivationFlags := make(map[string]Flag)
	ccFlag := Flag{
		symbol:      "cc",
		description: "Allows entering multiple cage cards in one pass",
		takesValue:  true,
	}
	deactivationFlags["-"+ccFlag.symbol] = ccFlag

	// ect as needed or remove the "-"+ for longer ones

	popFlag := Flag{
		symbol:      "pop",
		description: "Removes the last added cage card",
		takesValue:  false,
	}
	deactivationFlags[popFlag.symbol] = popFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing cage cards",
		takesValue:  false,
	}
	deactivationFlags[exitFlag.symbol] = exitFlag

	processFlag := Flag{
		symbol:      "process",
		description: "Processes cage cards and then exits",
		takesValue:  false,
	}
	deactivationFlags[processFlag.symbol] = processFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags for current command",
		takesValue:  false,
	}
	deactivationFlags[helpFlag.symbol] = helpFlag

	return deactivationFlags

}

// look into removing the args thing, might have to stay
func deactivateFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getDeactivationFlags()

	// set defaults
	exit := false
	cardsToDeactivate := []database.DeactivateCageCardParams{}
	date := time.Now()

	fmt.Println("Enter cards to deactivate.")
	// the reader
	reader := bufio.NewReader(os.Stdin)

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

		// try to run as a number, and add it to the list of cards to activate using the current values
		if len(inputs) == 1 {
			cc, err := strconv.Atoi(inputs[0])
			if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
				// an error occured and it was not from passing a word in to atoi
				fmt.Println("Error convering input to cage card number")
				fmt.Println(err)
				continue
			}

			// a misread on cc means the value 0 init
			if cc != 0 {
				tDccp := getCCToDeactivate(cc, &date, cfg.loggedInInvestigator)
				cardsToDeactivate = append(cardsToDeactivate, tDccp)
				fmt.Printf("%v card added\n", cc)
				continue
			}
		}

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// cc, pop, process, exit, list of errors (both previously deact and not activated)
		// TODO: make sure activated cards check if previously deact too
		for _, arg := range args {
			switch arg.flag {
			case "-cc":
				cc, err := strconv.Atoi(inputs[0])
				if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
					// an error occured and it was not from passing a word in to atoi
					fmt.Println("Error convering input to cage card number")
					fmt.Println(err)
					continue
				}

				// a misread on cc means the value 0 init
				if cc != 0 {
					tDccp := getCCToDeactivate(cc, &date, cfg.loggedInInvestigator)
					cardsToDeactivate = append(cardsToDeactivate, tDccp)
					fmt.Printf("%v card added\n", cc)
					continue
				}

			case "help":
				cmdHelp(flags)

			case "pop":
				length := len(cardsToDeactivate)
				if length == 0 {
					fmt.Println("No cards have been entered")
					break
				}
				popped := cardsToDeactivate[length-1]
				fmt.Printf("Popped %v\n", popped.CcID)
				cardsToDeactivate = cardsToDeactivate[0 : length-1]

			case "process":
				exit = true

			case "exit":
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

// is it more expensive to pass an int by pointer and deref or just pass by value
func getCCToDeactivate(cc int, date *time.Time, deactivatedBy *database.Investigator) database.DeactivateCageCardParams {
	tdate := sql.NullTime{Valid: true, Time: *date}
	tdeactivatedBy := uuid.NullUUID{Valid: true, UUID: deactivatedBy.ID}

	tdccp := database.DeactivateCageCardParams{
		CcID:          int32(cc),
		DeactivatedOn: tdate,
		DeactivatedBy: tdeactivatedBy,
	}

	return tdccp
}

func deactivateCageCards(cfg *Config, ctd []database.DeactivateCageCardParams) error {
	if len(ctd) == 0 {
		return errors.New("Oops! No cards!")
	}
	deactivationErrors := []ccError{}
	totalDeactivated := 0

	for _, cc := range ctd {
		// check if card is in the system

		// check if not active or previously deact

	}
}
