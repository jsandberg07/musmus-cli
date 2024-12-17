package main

import (
	"bufio"
	"context"
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

// IDEA: struct that contains an array of dates, and creates a cage card activation param with dates with a reference to that date
// then a true param is created with the value, pass to a go routine and activated
// a memory saving measure, help with a LOT of cards but really a few dates and IDs aren't anything major but it would be cool
// would also work for investigator

// TODO: set reminders with a flag. either by like E17 for a date or just a date

// TODO: update the sql for activation, need more fields
// actually making the activate cards uhh contact the db
// handling cards that are already active, previously deactivated and other misc errors
// changing date and adding to activation
// adding strain to activation
// adding notes to activation
// add a flag for cards, so you can add a note at the same time
// a print for what the current settings are
// capital S and N for keep notes and strains for multiple cages
// automatically add the "activated by"
// handling allotment / updating protocol (be fast and have structs with the prot uudis + total)

// then more of the same for deactivation ect
// cards that haven't been activated

func fart() {
	fmt.Println("Fart")
}

func getCCActivationCmd() Command {
	activateFlags := make(map[string]Flag)
	ccActivationCmd := Command{
		name:        "activate",
		description: "Used for activating cage cards",
		function:    activateFunction,
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

	nFlag := Flag{
		symbol:      "n",
		description: "Sets the note for only the next card to be added. Enter 'x' to clear\n Use underscores in place of spaces",
		takesValue:  true,
	}
	activateFlags["-"+nFlag.symbol] = nFlag

	NFlag := Flag{
		symbol:      "N",
		description: "Sets the note for all cage cards added until changes. Enter 'x' to clear\n Use underscores in place of spaces",
		takesValue:  true,
	}
	activateFlags["-"+NFlag.symbol] = NFlag

	sFlag := Flag{
		symbol:      "s",
		description: "Sets the strain for only the next card to be added. Enter 'x' to clear",
		takesValue:  true,
	}
	activateFlags["-"+sFlag.symbol] = sFlag

	SFlag := Flag{
		symbol:      "S",
		description: "Sets the strain for all cage cards added until changes. Enter 'x' to clear",
		takesValue:  true,
	}
	activateFlags["-"+SFlag.symbol] = SFlag

	ccFlag := Flag{
		symbol:      "cc",
		description: "Adds a cage card to the queue to be activated",
		takesValue:  true,
	}
	activateFlags["-"+ccFlag.symbol] = ccFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints the settings that will be applied to the next card added to the queue",
		takesValue:  false,
	}
	activateFlags[printFlag.symbol] = printFlag

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

func activateFunction(cfg *Config, args []Argument) error {

	flags := getActivationFlags()

	// set defaults for the command
	exit := false
	cardsToProcess := []database.TrueActivateCageCardParams{}
	date := time.Now()
	allotment := 0
	strain := database.Strain{ID: uuid.Nil}
	keepStrain := false
	notes := ""
	keepNote := false

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
				tAccp := getCCToAdd(cc, &date, &strain, cfg.loggedInInvestigator, &notes)
				cardsToProcess = append(cardsToProcess, tAccp)
				fmt.Printf("%v card added\n", cc)

				if keepNote == false {
					notes = ""
				}
				if keepStrain == false {
					strain.ID = uuid.Nil
				}
				continue
			}

		}

		// otherwise set values based on what was passed in, or process things
		args, err := parseArguments(flags, inputs)
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

			case "-cc":

			case "-a":
				fmt.Println("TODO: add allotments to the protocols")
				// set the allotment, just parsing an int how hard could it be
				// make sure you see if it's like above a gorillion or not

			case "-s":
				s, err := getStrainByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				strain = s
				keepStrain = false

			case "-S":
				s, err := getStrainByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				strain = s
				keepStrain = true

			case "-n":
				if arg.value == "x" || arg.value == "X" {
					notes = ""
				} else {
					notes = arg.value
				}
				keepNote = false

			case "-N":
				if arg.value == "x" || arg.value == "X" {
					notes = ""
				} else {
					notes = arg.value
				}
				keepNote = true

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
				popped := cardsToProcess[length-1]
				fmt.Printf("Popped %v\n", popped.CcID)
				cardsToProcess = cardsToProcess[0 : length-1]

			case "help":
				fmt.Println("Notes and strains can be added for individual cards, or set for many")
				fmt.Println("Then you can either add only cage cards, or mark a cage card for activation with -cc")
				err := cmdHelp(flags)
				if err != nil {
					fmt.Println(err)
				}

			case "print":
				printCurrentActivationParams(&date, &allotment, &strain, &notes)

			case "exit":
				fmt.Println("Exiting without processing")
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

func processCageCards(cctp []database.TrueActivateCageCardParams) error {
	if len(cctp) == 0 {
		return errors.New("Oops! No cards!")
	}
	for i := 0; i < len(cctp); i++ {
		fmt.Printf("Processing card %v... ;^3\n", i)
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

// at what point do you start passing single digit ints by reference?
func printCurrentActivationParams(date *time.Time, allotment *int, strain *database.Strain, note *string) {
	fmt.Println("Current settings for cards being added to activation queue:")
	fmt.Printf("Date: %v\n", date)
	fmt.Printf("Number of animals: %v\n", allotment)
	if strain.ID != uuid.Nil {
		fmt.Printf("Strain: %v\n", strain.SName)
	}
	if *note != "" {
		fmt.Printf("Notes: %s\n", *note)
	}
}

// works with both code and name
func getStrainByFlag(cfg *Config, input string) (database.Strain, error) {
	if input == "x" || input == "X" {
		return database.Strain{ID: uuid.Nil}, nil
	}
	strain, err := cfg.db.GetStrainByName(context.Background(), input)

	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error with DB
		fmt.Println("Error getting strain from DB")
		return database.Strain{ID: uuid.Nil}, err
	}

	if err == nil {
		// strain found by name
		return strain, nil
	}

	// look for it by code
	strain, err = cfg.db.GetStrainByCode(context.Background(), input)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error with DB
		fmt.Println("Error getting strain from DB")
		return database.Strain{ID: uuid.Nil}, err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		fmt.Println("Strain not found by name or number. Please try again.")
		return database.Strain{ID: uuid.Nil}, nil
	}

	// strain found by code
	return strain, nil
}

// probably a candidate for using channels and a go routine to feed this into another function
func getCCToAdd(cc int,
	date *time.Time,
	strain *database.Strain,
	activatedBy *database.Investigator,
	notes *string) database.TrueActivateCageCardParams {

	tdate := sql.NullTime{Valid: true, Time: *date}

	var tstrain uuid.NullUUID
	if strain.ID == uuid.Nil {
		tstrain.Valid = false
	} else {
		tstrain.Valid = true
		tstrain.UUID = strain.ID
	}

	var tnote sql.NullString
	if *notes == "" {
		tnote.Valid = false
	} else {
		tnote.Valid = true
		tnote.String = *notes
	}

	tactivatedBy := uuid.NullUUID{Valid: true, UUID: activatedBy.ID}

	taccp := database.TrueActivateCageCardParams{
		CcID:        int32(cc),
		ActivatedOn: tdate,
		Strain:      tstrain,
		ActivatedBy: tactivatedBy,
		Notes:       tnote,
	}
	return taccp
}
