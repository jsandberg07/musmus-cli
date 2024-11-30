package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

// cage card processing to actually use the DB
// having a login first is best actually lets do that hella fast

func fart() {
	fmt.Println("Fart")
}

func getCCActivationCmd() Command {
	// subcommand that starts its own loop
	activateFlags := make(map[string]Flag)
	ccActivationCmd := Command{
		name:        "activate",
		description: "Used for activating cage cards",
		function:    activateSubcommand,
		flags:       activateFlags,
	}

	return ccActivationCmd
}

func getActivationFlags() map[string]Flag {

	activateFlags := make(map[string]Flag)
	dFlag := Flag{
		symbol:      "d",
		description: "Sets Date. Use format MM/DD/YYYY",
		takesValue:  true,
	}
	activateFlags["-"+dFlag.symbol] = dFlag

	aFlag := Flag{
		symbol:      "a",
		description: "Sets number of animals added to protocol on activation",
		takesValue:  true,
	}
	activateFlags["-"+aFlag.symbol] = aFlag

	processFlag := Flag{
		symbol:      "process",
		description: "Processes cage cards that have been entered then exits",
		takesValue:  false,
	}
	activateFlags[processFlag.symbol] = processFlag

	popFlag := Flag{
		symbol:      "pop",
		description: "Deletes the most recently scanned cage card",
		takesValue:  false,
	}
	activateFlags[popFlag.symbol] = popFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing cards",
		takesValue:  false,
	}
	activateFlags[exitFlag.symbol] = exitFlag

	return activateFlags
}

// add everything to a list with the details
// literally make an array of dates when entering them to save space because FUCK YEAH
// then put them in params, make an array of params
// then for NOW acticate them individually
// do batch later
// return and print errors
func activateSubcommand(cfg *Config, args []Argument) error {
	// start another loop
	// parse subflags and set from there
	// dont mix commands and cards and flags
	flags := getActivationFlags()

	// set defaults for the command
	exit := false
	cardsToProcess := []database.ActivateCageCardParams{}
	date := time.Now()

	reader := bufio.NewReader(os.Stdin)
	for {

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

		}

		// otherwise set values based on what was passed in, or process things
		args, err := parseSubcommand(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-d":
				// parse the time, dont fuck it up

			case "-a":
				// set the allotment, just parsing an int how hard could it be
				// make sure you see if it's like above a gorillion or not
			case "process":
				err := processCageCards(cardsToProcess)
				if err != nil {
					fmt.Println(err)
				}
				exit = true
			case "pop":
				// delete one
			case "exit":
				exit = true
			default:
				fmt.Printf("Oops a fake flag snuck in: %s", arg.flag)
			}
		}

		if exit {
			break
		}
	}

	return nil

}

func processCageCards(cctp []database.ActivateCageCardParams) error {
	if len(cctp) == 0 {
		return errors.New("Oops! No cards!")
	}
	for i := 0; i < len(cctp); i++ {
		fmt.Printf("Processing card %v... ;^3", i)
	}

	return nil
}
