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

func getCCDeactivationCmd() Command {
	deactivationFlags := make(map[string]Flag)
	deactivationCmd := Command{
		name:        "deactivate",
		description: "Used for deactivating cage cards",
		function:    deactivateFunction,
		flags:       deactivationFlags,
		printOrder:  2,
	}

	return deactivationCmd
}

// enter number
// try to deactivate
// get error if need be
// delete reminder
// its easier trust me

// cc, pop, process, exit, list of errors (both previously deact and not activated)
// TODO: change this to work like activation (ie linear, check reminders)
func getDeactivationFlags() map[string]Flag {
	deactivationFlags := make(map[string]Flag)
	ccFlag := Flag{
		symbol:      "-cc",
		description: "Allows entering multiple cage cards in one pass",
		takesValue:  true,
		printOrder:  1,
	}
	deactivationFlags[ccFlag.symbol] = ccFlag

	dFlag := Flag{
		symbol:      "-d",
		description: "Sets the date the cage card will be deactivated",
		takesValue:  true,
		printOrder:  2,
	}
	deactivationFlags[dFlag.symbol] = dFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints the current settings for card deactivation",
		takesValue:  false,
		printOrder:  3,
	}
	deactivationFlags[printFlag.symbol] = printFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without processing cage cards",
		takesValue:  false,
		printOrder:  100,
	}
	deactivationFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags for current command",
		takesValue:  false,
		printOrder:  100,
	}
	deactivationFlags[helpFlag.symbol] = helpFlag

	return deactivationFlags

}

func deactivateFunction(cfg *Config) error {
	// permission check
	err := checkPermission(cfg.loggedInPosition, PermissionDeactivateReactivate)
	if err != nil {
		return err
	}
	// get flags
	flags := getDeactivationFlags()

	// set defaults
	exit := false
	date := normalizeDate(time.Now())

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
				err := deactivateCC(cfg, cc, date)
				if err != nil {
					fmt.Println(err)
				}

				// don't need to visit the switch, one input is assumed to be a cc#
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

			case "-d":
				newDate, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}

				// cant be after today ie cant deactivate on a day that hasnt happened
				if newDate.After(time.Now()) {
					fmt.Println("Deactivation date can't be set in the future")
					break
				}

				date = normalizeDate(newDate)
				fmt.Printf("Date set: %v\n", date)

			// TODO: -cc isnt working for some reason but im not testing it atm
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
					err := deactivateCC(cfg, cc, date)
					if err != nil {
						fmt.Println(err)
					}
				}

			case "print":
				printCurrentDeactivationParams(&date)

			case "help":
				cmdHelp(flags)

				/* removed because cards are no longer kept in a queue
				case "pop":
					length := len(cardsToDeactivate)
					if length == 0 {
						fmt.Println("No cards have been entered")
						break
					}
					popped := cardsToDeactivate[length-1]
					fmt.Printf("Popped %v\n", popped.CcID)
					cardsToDeactivate = cardsToDeactivate[0 : length-1]
				*/

				/* removed because cards are no longer kept in a queue
				case "process":
					fmt.Println("Processing...")
					err := deactivateCageCards(cfg, cardsToDeactivate)
					if err != nil {
						fmt.Println(err)
					}
					exit = true
				*/

			case "exit":
				fmt.Println("Exiting without saving...")
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

// yeah, just the date. Keep the 'deactivated_by' hidden
func printCurrentDeactivationParams(date *time.Time) {
	fmt.Println("Current settings for cards being added to deactivation queue:")
	fmt.Printf("Date: %v\n", *date)
}

// TODO: naming things is hard. really just need the ccID huh
func deactivateCC(cfg *Config, ccID int, date time.Time) error {
	cc, err := cfg.db.GetCageCardByID(context.Background(), int32(ccID))
	if err != nil && err.Error() == "sql: no rows in result set" {
		// not added to DB
		return errors.New("cage card not found in DB")
	}
	if err != nil {
		// any other error
		return err
	}

	if !cc.ActivatedOn.Valid {
		return errors.New("cage card not currently active")
	}
	if date.Before(cc.ActivatedOn.Time) {
		// can't have a deactivation date set before activation
		msg := fmt.Sprintf("deactivation date can't be before activation -- %v", cc.ActivatedOn.Time)
		return errors.New(msg)
	}
	if cc.DeactivatedOn.Valid {
		msg := fmt.Sprintf("cage card previously deactivated -- %v", cc.DeactivatedOn.Time)
		return errors.New(msg)
	}

	dccParams := database.DeactivateCageCardParams{
		CcID:          int32(ccID),
		DeactivatedOn: sql.NullTime{Valid: true, Time: date},
		DeactivatedBy: uuid.NullUUID{Valid: true, UUID: cfg.loggedInInvestigator.ID},
	}
	deactivatedCC, err := cfg.db.DeactivateCageCard(context.Background(), dccParams)
	if err != nil {
		fmt.Println("Error deactivating cage card in DB")
		return err
	}
	fmt.Printf("%v deactivated!\n", deactivatedCC.CcID)
	if verbose {
		fmt.Println(deactivatedCC)
	}

	// check for reminders and delete those
	reminders, err := cfg.db.GetRemindersByCC(context.Background(), deactivatedCC.CcID)
	if err != nil {
		fmt.Println("Error getting reminders from DB")
		return err
	}
	if len(reminders) == 0 {
		// no reminders, just return
		return nil
	}
	fmt.Println("Reminders found for cage card:")
	for _, r := range reminders {
		fmt.Printf("%v -- %s\n", r.Note, r.RDate)
	}
	for i, r := range reminders {
		err := cfg.db.DeleteReminder(context.Background(), r.ID)
		if err != nil {
			fmt.Printf("Error deleteing reminder %v -- %v\n", i, r.Note)
		}
	}
	fmt.Println("Reminders have been deleted.")

	return nil
}

/* removed because no longer processing CCs via a queue
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

// candidate for DRYing up with a "process cc" function with an activate/deactivate enum
// behavior is just different enough to have them disentangled
// ie activating checks to see if it's already active, deact checks to see if it isnt active
func deactivateCageCards(cfg *Config, ctd []database.DeactivateCageCardParams) error {
	if len(ctd) == 0 {
		return errors.New("oops! No cards")
	}
	deactivationErrors := []ccError{}
	totalDeactivated := 0

	for _, cc := range ctd {
		ccErr := checkDeactivateError(cfg, &cc)
		// hacky way to see if a nil struct was returned, meaning no error
		if ccErr.CCid != 0 {
			deactivationErrors = append(deactivationErrors, ccErr)
			continue
		}

		dcc, err := cfg.db.DeactivateCageCard(context.Background(), cc)
		if err != nil {
			// any other error
			tcce := ccError{
				CCid: int(dcc.CcID),
				Err:  err.Error(),
			}
			deactivationErrors = append(deactivationErrors, tcce)
			continue
		}

		if verbose {
			fmt.Println(dcc)
		}

		totalDeactivated++

	}

	fmt.Printf("%v cards deactivated\n", totalDeactivated)
	if len(deactivationErrors) > 0 {
		fmt.Println("Errors deactivating these cage cards:")
		for _, cce := range deactivationErrors {
			fmt.Printf("%v -- %s\n", cce.CCid, cce.Err)
		}
	}
	return nil
}


// TODO: maybe add a check for if deactivation date is after today too
// like can only deactivate today or past, not future, to prevent errors of course
// no "this card will have had been deactivated"
// WORKING: seeing why stopping cards isn't working
// I FORGOT TO FINISH THE CASES LE MOO ALSO ADD A DATE SETTER AND CHECK IF ITS IN THE FUTURE
func checkDeactivateError(cfg *Config, cc *database.DeactivateCageCardParams) ccError {
	// check if already active
	td, err := cfg.db.GetActivationDate(context.Background(), cc.CcID)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// cc not added to db or not found
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  "CC not added to database",
		}

		return tcce
	}

	if !td.Valid {
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  "CC is not currently active",
		}
		return tcce
	}

	// check if deactivation date is before activation date
	if cc.DeactivatedOn.Time.Before(td.Time) {
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  "Deactivation date is before activation date",
		}
		return tcce
	}

	if err != nil {
		// any other error
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  err.Error(),
		}
		return tcce
	}

	// check if previously deactivated
	dd, err := cfg.db.GetDeactivationDate(context.Background(), cc.CcID)
	// dont need to check if not in db
	if dd.Valid {
		// card was previously deactivated
		errmsg := fmt.Sprintf("CC is already deactivated -- %s", dd.Time)
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  errmsg,
		}
		return tcce
	}

	if err != nil {
		// any other error
		tcce := ccError{
			CCid: int(cc.CcID),
			Err:  err.Error(),
		}
		return tcce
	}

	// everything ok
	return ccError{}
}
*/
