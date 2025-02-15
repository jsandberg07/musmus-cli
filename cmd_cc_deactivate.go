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
	err := checkPermission(cfg.loggedInPosition, PermissionDeactivateReactivate)
	if err != nil {
		return err
	}
	flags := getDeactivationFlags()

	exit := false
	date := normalizeDate(time.Now())

	fmt.Println("Enter cards to deactivate.")
	reader := bufio.NewReader(os.Stdin)
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

				if newDate.After(normalizeDate(time.Now())) {
					fmt.Println("Deactivation date can't be set in the future")
					break
				}

				date = normalizeDate(newDate)
				fmt.Printf("Date set: %v\n", date)

			case "-cc":
				cc, err := strconv.Atoi(arg.value)
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

			case "print":
				printCurrentDeactivationParams(&date)

			case "help":
				cmdHelp(flags)

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

func printCurrentDeactivationParams(date *time.Time) {
	fmt.Println("Current settings for cards being added to deactivation queue:")
	fmt.Printf("Date: %v\n", *date)
}

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
