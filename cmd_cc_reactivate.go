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

func getCCReactivateCmd() Command {
	reactivateFlags := make(map[string]Flag)
	reactivateCmd := Command{
		name:        "reactivate",
		description: "Reactivate cage cards, removing their deactivation date",
		function:    reactivateFunction,
		flags:       reactivateFlags,
		printOrder:  4,
	}

	return reactivateCmd
}

func getCCReactivateFlags() map[string]Flag {
	ReactivateFlags := make(map[string]Flag)

	ccFlag := Flag{
		symbol:      "-cc",
		description: "Adds CC to queue to be reactivated",
		takesValue:  true,
		printOrder:  1,
	}
	ReactivateFlags[ccFlag.symbol] = ccFlag

	processFlag := Flag{
		symbol:      "process",
		description: "Reactivates card in queue",
		takesValue:  false,
		printOrder:  2,
	}
	ReactivateFlags[processFlag.symbol] = processFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing queue",
		takesValue:  false,
		printOrder:  100,
	}
	ReactivateFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available flags for this command",
		takesValue:  false,
		printOrder:  100,
	}
	ReactivateFlags[helpFlag.symbol] = helpFlag

	return ReactivateFlags
}

func reactivateFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionDeactivateReactivate)
	if err != nil {
		return err
	}

	flags := getCCReactivateFlags()

	exit := false
	cardsToReactivate := []int{}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter cards to reactivate")

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
				cardsToReactivate = append(cardsToReactivate, cc)
				continue
			}
		}

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
				cardsToReactivate = append(cardsToReactivate, cc)
				fmt.Printf("%v card added\n", cc)

			case "process":
				err := reactivateCageCards(cfg, cardsToReactivate)
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

func reactivateCageCards(cfg *Config, ctr []int) error {
	if len(ctr) == 0 {
		return errors.New("oops! No cards")
	}
	reactivationErrors := []ccError{}
	totalReactivated := 0

	for _, cc := range ctr {
		ccErr := checkReactivateError(cfg, cc)
		if ccErr.CCid != 0 {
			reactivationErrors = append(reactivationErrors, ccErr)
			continue
		}

		err := cfg.db.ReactivateCageCard(context.Background(), int32(cc))
		if err != nil {
			tcce := ccError{
				CCid: cc,
				Err:  err.Error(),
			}
			reactivationErrors = append(reactivationErrors, tcce)
			continue
		}
		totalReactivated++
	}

	fmt.Printf("%v cards reactivated\n", totalReactivated)
	if len(reactivationErrors) > 0 {
		fmt.Println("Errors reactivating these cage cards:")
		for _, cce := range reactivationErrors {
			fmt.Printf("%v -- %s\n", cce.CCid, cce.Err)
		}
	}

	return nil

}

func checkReactivateError(cfg *Config, cc int) ccError {
	// check if not deactivated at all
	dd, err := cfg.db.GetDeactivationDate(context.Background(), int32(cc))
	if err != nil && err.Error() == "sql: no rows in result set" {
		// card not in db
		tcce := ccError{
			CCid: int(cc),
			Err:  "CC not added to database",
		}
		return tcce
	}

	if err != nil {
		// any other error
		tcce := ccError{
			CCid: int(cc),
			Err:  err.Error(),
		}
		return tcce
	}

	if !dd.Valid {
		// not deactivated anyway
		tcce := ccError{
			CCid: int(cc),
			Err:  "CC was not deactivated",
		}
		return tcce
	}

	// everything ok
	return ccError{}
}
