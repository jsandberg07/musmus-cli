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
		name:        "2activate",
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

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints help messages and flags for commands available",
		takesValue:  false,
	}
	activateFlags[helpFlag.symbol] = helpFlag

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
	fmt.Println("Cage card activation.")
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
		fmt.Printf("** %v\n", inputs)

		// try to run as a number, and add it to the list of cards to activate using the current values
		if len(inputs) == 1 {
			cc, err := strconv.Atoi(inputs[0])
			if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
				// an error occured and it was not from passing a word in to atoi
				fmt.Println("Error convering input to cage card number")
				fmt.Println(err)
				continue
			}
			if cc != 0 {
				tAccp := database.ActivateCageCardParams{
					CcID:           int32(cc),
					ActivatedOn:    sql.NullTime{Valid: true, Time: date},
					InvestigatorID: cfg.loggedInInvestigator.ID,
				}
				fmt.Println("** card added!")
				cardsToProcess = append(cardsToProcess, tAccp)
				continue
			}
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
				newDate, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				date = newDate
				fmt.Printf("Date set: %v\n", date)
			case "-a":
				fmt.Println("TODO: add allotments to the protocols")
				// set the allotment, just parsing an int how hard could it be
				// make sure you see if it's like above a gorillion or not
			case "process":
				err := processCageCards(cardsToProcess)
				if err != nil {
					fmt.Println(err)
				}
				exit = true
			case "pop":
				length := len(cardsToProcess)
				if length == 0 {
					fmt.Println("No cards have been entered")
					break
				}
				cardsToProcess = cardsToProcess[0 : length-1]
			case "help":
				err := scmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}
			case "exit":
				fmt.Println("Exiting without processing")
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

// because who knows what dates people are going to enter
// TODO: make sure there isnt fuckery like "this card technically wasn't this day beacuase the time was off"
// set everything to be active at like midnight, queries at midnight and see if it works otherwise +1 second lmao

func parseDate(input string) (time.Time, error) {
	// create an array of the formats (with 0s, without, 4 digit year, 2 digit year)
	// go through parse works and then return
	var date time.Time
	var err error
	timeFormats := []string{"1/2/06", "1/2/2006", "01/02/06", "01/02/2006"}
	for _, format := range timeFormats {
		date, err = time.Parse(format, input)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("Error parsing date.")
		return time.Time{}, err
	}

	return date, nil
}
