package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// return to inactive

func getCCInactivateCmd() Command {
	inactivateFlags := make(map[string]Flag)
	inactivateCmd := Command{
		name:        "inactivate",
		description: "Used for returning cards to inactive status",
		function:    inactivateFunction,
		flags:       inactivateFlags,
		printOrder:  5,
	}

	return inactivateCmd
}

func getInactivateFlags() map[string]Flag {
	// cc, process, exit, help
	InactivateFlags := make(map[string]Flag)
	ccFlag := Flag{
		symbol:      "-cc",
		description: "Adds CC to queue to be reactivated",
		takesValue:  true,
		printOrder:  1,
	}
	InactivateFlags[ccFlag.symbol] = ccFlag

	// ect as needed or remove the "-"+ for longer ones

	// maybe remove this
	processFlag := Flag{
		symbol:      "process",
		description: "Reactivates card in queue",
		takesValue:  false,
		printOrder:  2,
	}
	InactivateFlags[processFlag.symbol] = processFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing queue",
		takesValue:  false,
		printOrder:  100,
	}
	InactivateFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available flags for this command",
		takesValue:  false,
		printOrder:  100,
	}
	InactivateFlags[helpFlag.symbol] = helpFlag

	return InactivateFlags

}

// look into removing the args thing, might have to stay
func inactivateFunction(cfg *Config) error {
	// permission check
	err := checkPermission(cfg.loggedInPosition, PermissionActivateInactivate)
	if err != nil {
		return err
	}
	// get flags
	flags := getInactivateFlags()

	// set defaults
	exit := false
	cardsToInactivate := []int{}

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter cards to inactivate")

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
				cardsToInactivate = append(cardsToInactivate, cc)
				continue
			}
		}

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-cc":
				cc, err := strconv.Atoi(arg.value)
				if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
					// an error occured and it was not from passing a word in to atoi
					fmt.Println("Error convering input to cage card number")
					fmt.Println(err)
					continue
				}
				cardsToInactivate = append(cardsToInactivate, cc)
				fmt.Printf("%v card added\n", cc)

			case "process":
				err := inactivateCageCards(cfg, cardsToInactivate)
				if err != nil {
					return err
				}
				exit = true

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

func inactivateCageCards(cfg *Config, cti []int) error {
	if len(cti) == 0 {
		return errors.New("oops! No cards")
	}
	inactivationErrors := []ccError{}
	totalInactivated := 0

	for _, cc := range cti {
		ccErr := checkInactivateError(cfg, cc)
		if ccErr.CCid != 0 {
			inactivationErrors = append(inactivationErrors, ccErr)
			continue
		}

		err := cfg.db.InactivateCageCard(context.Background(), int32(cc))
		if err != nil {
			tcce := ccError{
				CCid: cc,
				Err:  err.Error(),
			}
			inactivationErrors = append(inactivationErrors, tcce)
			continue
		}

		// nothing to verbose
		totalInactivated++
	}

	fmt.Printf("%v cards inactivated\n", totalInactivated)
	if len(inactivationErrors) > 0 {
		fmt.Println("Errors reactivating these cage cards:")
		for _, cce := range inactivationErrors {
			fmt.Printf("%v -- %s\n", cce.CCid, cce.Err)
		}
	}

	return nil
}

func checkInactivateError(cfg *Config, cc int) ccError {
	// check if not in db
	ad, err := cfg.db.GetActivationDate(context.Background(), int32(cc))
	if err != nil && err.Error() == "sql: no rows in result set" {
		// card not in db
		tcce := ccError{
			CCid: int(cc),
			Err:  "CC not added to database",
		}
		return tcce
	}

	// any other db error
	if err != nil {
		// any other error
		tcce := ccError{
			CCid: int(cc),
			Err:  err.Error(),
		}
		return tcce
	}

	if !ad.Valid {
		// not already active
		tcce := ccError{
			CCid: int(cc),
			Err:  "CC was not activated",
		}
		return tcce
	}

	// check if already deactivated
	dd, err := cfg.db.GetDeactivationDate(context.Background(), int32(cc))
	if err != nil {
		// any other error
		tcce := ccError{
			CCid: int(cc),
			Err:  err.Error(),
		}
		return tcce
	}

	// previously deactivated
	if dd.Valid {
		tcce := ccError{
			CCid: int(cc),
			Err:  "CC is deactivated",
		}
		return tcce
	}

	// everything ok
	return ccError{}
}
