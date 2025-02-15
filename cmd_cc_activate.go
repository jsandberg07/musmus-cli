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

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getCCActivationCmd() Command {
	activateFlags := make(map[string]Flag)
	ccActivationCmd := Command{
		name:        "activate",
		description: "Used for activating cage cards",
		function:    activateFunction,
		flags:       activateFlags,
		printOrder:  1,
	}

	return ccActivationCmd
}

func getActivationFlags() map[string]Flag {
	activateFlags := make(map[string]Flag)

	ccFlag := Flag{
		symbol:      "-cc",
		description: "Add multiple cage cards at once. Will be activated in order it is entered (including other flags)",
		takesValue:  true,
		printOrder:  1,
	}
	activateFlags[ccFlag.symbol] = ccFlag

	dFlag := Flag{
		symbol:      "-d",
		description: "Sets Date. Use format MM/DD/YYYY",
		takesValue:  true,
		printOrder:  2,
	}
	activateFlags[dFlag.symbol] = dFlag

	aFlag := Flag{
		symbol:      "-a",
		description: "Sets number of animals added to protocol on activation",
		takesValue:  true,
		printOrder:  3,
	}
	activateFlags[aFlag.symbol] = aFlag

	nFlag := Flag{
		symbol:      "-n",
		description: "Sets the note for only the next card to be added. Enter 'x' to clear\n Use underscores in place of spaces",
		takesValue:  true,
		printOrder:  4,
	}
	activateFlags[nFlag.symbol] = nFlag

	NFlag := Flag{
		symbol:      "-N",
		description: "Sets the note for all cage cards added until changes. Enter 'x' to clear\n Use underscores in place of spaces",
		takesValue:  true,
		printOrder:  5,
	}
	activateFlags[NFlag.symbol] = NFlag

	sFlag := Flag{
		symbol:      "s-",
		description: "Sets the strain for only the next card to be added. Enter 'x' to clear",
		takesValue:  true,
		printOrder:  6,
	}
	activateFlags[sFlag.symbol] = sFlag

	SFlag := Flag{
		symbol:      "-S",
		description: "Sets the strain for all cage cards added until changed. Enter 'x' to clear",
		takesValue:  true,
		printOrder:  7,
	}
	activateFlags[SFlag.symbol] = SFlag

	rFlag := Flag{
		symbol:      "-r",
		description: "Will add a reminder input days after activation date. \nRequires a note for the reminder. Enter 'x' to clear",
		takesValue:  true,
		printOrder:  8,
	}
	activateFlags[rFlag.symbol] = rFlag

	RFlag := Flag{
		symbol:      "-R",
		description: "Will add a reminder input days after activation date to all cages until changes. \nRequires a note for the reminder. Enter 'x' to clear",
		takesValue:  true,
		printOrder:  9,
	}
	activateFlags[RFlag.symbol] = RFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints the settings that will be applied to the next card added to the queue",
		takesValue:  false,
		printOrder:  100,
	}
	activateFlags[printFlag.symbol] = printFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints help messages and flags for commands available",
		takesValue:  false,
		printOrder:  100,
	}
	activateFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing cards",
		takesValue:  false,
		printOrder:  100,
	}
	activateFlags[exitFlag.symbol] = exitFlag

	return activateFlags
}

func activateFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionActivateInactivate)
	if err != nil {
		return err
	}

	flags := getActivationFlags()

	exit := false
	ccParams := CageCardActivationParams{}
	ccParams.init()

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
			ccID, err := strconv.Atoi(inputs[0])
			if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
				// an error occured and it was not from passing a word in to atoi
				fmt.Println("Error convering input to cage card number")
				fmt.Println(err)
				continue
			}

			// a misread on cc means the value 0 init
			if ccID != 0 {
				ccParams.ccID = ccID
				err := activationWrapper(cfg, &ccParams)
				if err != nil {
					fmt.Println(err)
				}
				// don't need to visit the switch, one input is assumed to be a cc#
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
			case "-d":

				newDate, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				ccParams.date = normalizeDate(newDate)
				fmt.Printf("Date set: %v\n", ccParams.date)

			case "-cc":
				ccID, err := strconv.Atoi(arg.value)
				if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
					// an error occured and it was not from passing a word in to atoi
					fmt.Println("Error convering input to cage card number")
					fmt.Println(err)
					continue
				}
				if err != nil {
					fmt.Println("Invalid entry. Please enter an integer for cage card")
					continue
				}
				ccParams.ccID = ccID
				err = activationWrapper(cfg, &ccParams)
				if err != nil {
					fmt.Println(err)
				}

			case "-a":
				num, err := strconv.Atoi(arg.value)
				if err != nil && !strings.Contains(err.Error(), "invalid syntax") {
					// an error occured and it was not from passing a word in to atoi
					fmt.Println("Error convering input to cage card number")
					fmt.Println(err)
					continue
				}
				if err != nil {
					fmt.Println("Invalid entry. Please enter an integer for allotment")
					continue
				}
				if num < 0 {
					ccParams.allotment = 0
				} else {
					ccParams.allotment = num
				}

			case "-s":
				fallthrough
			case "-S":
				if arg.value == "x" || arg.value == "X" {
					fmt.Println("Strain reset")
					ccParams.strain = database.Strain{ID: uuid.Nil}
					ccParams.keepStrain = false
					continue
				}
				s, err := getStrainByFlag(cfg, arg.value)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Strain set")
				ccParams.strain = s
				if arg.flag == "-s" {
					ccParams.keepStrain = false
				} else {
					ccParams.keepStrain = true
				}

			case "-n":
				fallthrough
			case "-N":
				if arg.value == "x" || arg.value == "X" {
					fmt.Println("Note reset")
					ccParams.note = ""
					ccParams.keepNote = false
					continue
				} else {
					fmt.Println("Note set")
					ccParams.note = arg.value
				}
				if arg.flag == "-n" {
					ccParams.keepNote = false
				} else {
					ccParams.keepNote = true
				}

			case "-r":
				fallthrough
			case "-R":
				if arg.value == "x" || arg.value == "X" {
					fmt.Println("Reminder cleared")
					ccParams.daysReminder = 0
					ccParams.keepReminder = false
					continue
				}
				num, err := strconv.Atoi(arg.value)
				if err != nil {
					fmt.Println("Error getting days for reminder. Please try again")
					fmt.Println(err)
				}
				ccParams.daysReminder = num
				fmt.Printf("Reminder will be set %v days from activation\n", ccParams.daysReminder)

				if arg.flag == "-r" {
					ccParams.keepReminder = false
				} else {
					ccParams.keepReminder = true
				}

			case "help":
				fmt.Println("Notes and strains can be added for individual cards, or set for many")
				fmt.Println("Then you can either add only cage cards, or mark a cage card for activation with -cc")
				cmdHelp(flags)

			case "print":
				err := printCurrentActivationParams(cfg, &ccParams)
				if err != nil {
					fmt.Println(err)
				}

			case "exit":
				fmt.Println("Exiting...")
				exit = true

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

func printCurrentActivationParams(cfg *Config, s *CageCardActivationParams) error {
	fmt.Println("Current settings for cage cards being activated:")
	fmt.Printf("* Date - %v\n", s.date)
	fmt.Printf("* Allotment - %v\n", s.allotment)

	if s.strain.ID != uuid.Nil {
		strain, err := cfg.db.GetStrainByID(context.Background(), s.strain.ID)
		if err != nil {
			fmt.Println("Could not get strain name from DB")
			return err
		}
		fmt.Printf("* Strain - %v -- Sticky -- %v\n", strain.SName, s.keepStrain)
	}

	if s.note != "" {
		fmt.Printf("* Note - %s -- Sticky -- %v\n", s.note, s.keepNote)
	}

	if s.daysReminder != 0 {
		fmt.Printf("* Reminder - %v days after activation -- Sticky -- %v\n", s.daysReminder, s.keepReminder)
	}

	return nil

}

// Handles all the processes of activating a cage card. Sets activation date, checks to see if optional activation params should be kept,
// adds balance to protocol, and checks for reminders.
func activationWrapper(cfg *Config, s *CageCardActivationParams) error {
	activateParams := getCCToActivate(s, cfg.loggedInInvestigator)
	activatedCC, err := activateCageCard(cfg, &activateParams)
	if err != nil {
		return err
	}
	fmt.Printf("%v activated!\n", activatedCC.CcID)
	if verbose {
		fmt.Println(activatedCC)
	}

	s.keepCheck()

	// add animals to allotment
	if s.allotment != 0 {
		err := addBalanceToProtocol(cfg, s.allotment, &activatedCC)
		if err != nil {
			fmt.Println("Could not add balance to protocol")
			fmt.Println(err)
		}
	}

	// check if reminder should be created
	if s.daysReminder != 0 {
		if !activatedCC.Notes.Valid {
			fmt.Println("Can't create a reminder without a note")
			return nil
		}
		err := ccActivationReminder(cfg, s.daysReminder, &activatedCC)
		if err != nil {
			fmt.Println("Could not create reminder")
			fmt.Println(err)
		}
	}

	return nil

}

// Helper that creates cage card activation struct
func getCCToActivate(s *CageCardActivationParams, i *database.Investigator) database.ActivateCageCardParams {

	tdate := sql.NullTime{
		Valid: true,
		Time:  s.date,
	}

	var tstrain uuid.NullUUID
	if s.strain.ID == uuid.Nil {
		tstrain.Valid = false
	} else {
		tstrain.Valid = true
		tstrain.UUID = s.strain.ID
	}

	var tnote sql.NullString
	if s.note == "" {
		tnote.Valid = false
	} else {
		tnote.Valid = true
		tnote.String = s.note
	}

	tactivatedBy := uuid.NullUUID{Valid: true, UUID: i.ID}

	taccp := database.ActivateCageCardParams{
		CcID:        int32(s.ccID),
		ActivatedOn: tdate,
		Strain:      tstrain,
		ActivatedBy: tactivatedBy,
		Notes:       tnote,
	}
	return taccp
}

func activateCageCard(cfg *Config, cc *database.ActivateCageCardParams) (database.CageCard, error) {
	tCard, err := cfg.db.GetCageCardByID(context.Background(), int32(cc.CcID))
	// check if added to db
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.CageCard{}, errors.New("cage card ID has not been added to the DB. Needs to be added first")
	}

	// db error
	if err != nil {
		fmt.Println("Error retrieving cage card details")
		return database.CageCard{}, err
	}

	// check if previously activated
	if tCard.ActivatedOn.Valid {
		// check if previously deactivated
		if tCard.DeactivatedOn.Valid {
			msg := fmt.Sprintf("cage card was previously deactivated %v", tCard.DeactivatedOn.Time)
			return database.CageCard{}, errors.New(msg)
		} else {
			msg := fmt.Sprintf("cage card was previously activated %v", tCard.ActivatedOn.Time)
			return database.CageCard{}, errors.New(msg)
		}
	}

	activatedCard, err := cfg.db.ActivateCageCard(context.Background(), *cc)
	if err != nil {
		fmt.Println("Error activating cage card")
		return database.CageCard{}, err
	}
	return activatedCard, nil
}
