package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddReminderCmd() Command {
	addReminderFlags := make(map[string]Flag)
	addReminderCmd := Command{
		name:        "add",
		description: "Used for adding reminders",
		function:    addReminderFunction,
		flags:       addReminderFlags,
		printOrder:  1,
	}

	return addReminderCmd
}

func getAddReminderFlags() map[string]Flag {
	addReminderFlags := make(map[string]Flag)

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the entered reminder",
		takesValue:  false,
		printOrder:  99,
	}
	addReminderFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	addReminderFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags for command",
		takesValue:  false,
		printOrder:  100,
	}
	addReminderFlags[helpFlag.symbol] = helpFlag

	return addReminderFlags

}

func addReminderFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionReminders)
	if err != nil {
		return err
	}

	flags := getAddReminderFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)

	date, err := getDatePrompt("Enter date for the reminder")
	if err != nil {
		return err
	}
	nilDate := time.Time{}
	if date == nilDate {
		fmt.Println("Exiting...")
		return nil
	}

	cc, err := getStructPrompt(cfg, "Enter cage card id for reminder. Cage card must be active", getCageCardStructActive)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	investigator, err := getStructPrompt(cfg, "Enter investigator who will recieve the reminder", getInvestigatorStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	note, err := getStringPrompt(cfg, "Enter a note for the reminder", checkFuncNil)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	arParams := database.AddReminderParams{
		RDate:          date,
		RCcID:          cc.CcID,
		InvestigatorID: investigator.ID,
		Note:           note,
	}

	fmt.Println("Reminder will be created with the following info:")
	printAddReminder(&arParams, &investigator)
	fmt.Println("Enter 'save' to keep or 'exit' to discard")

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

		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "save":
				fmt.Println("Saving...")
				reminder, err := cfg.db.AddReminder(context.Background(), arParams)
				if err != nil {
					fmt.Println("Error creating reminder")
					return err
				}
				if verbose {
					fmt.Println(reminder)
				}

				exit = true

			case "exit":
				fmt.Println("Exiting...")
				exit = true

			case "help":
				cmdHelp(getAddReminderFlags())

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

func printAddReminder(r *database.AddReminderParams, i *database.Investigator) {
	fmt.Printf("Date: %v\n", r.RDate)
	fmt.Printf("Cage Card: %v\n", r.RCcID)
	fmt.Printf("Investigator: %v\n", i.IName)
	fmt.Printf("Note: %v\n", r.Note)
}

func getDeleteReminderCmd() Command {
	deleteReminderFlags := make(map[string]Flag)
	deleteReminderCmd := Command{
		name:        "delete",
		description: "Used for deleting reminders",
		function:    deleteReminderFunction,
		flags:       deleteReminderFlags,
		printOrder:  2,
	}

	return deleteReminderCmd
}

func deleteReminderFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionReminders)
	if err != nil {
		return err
	}

	var reminders []database.Reminder

	date, err := getDatePrompt("Enter date of reminder to delete")
	if err != nil {
		return err
	}
	nilDate := time.Time{}
	if date == nilDate {
		fmt.Println("exiting...")
		return nil
	}
	date = normalizeDate(date)

	udrp := database.GetUserDayReminderParams{
		InvestigatorID: cfg.loggedInInvestigator.ID,
		RDate:          date,
	}
	reminders, err = cfg.db.GetUserDayReminder(context.Background(), udrp)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}

	if len(reminders) == 0 {
		fmt.Println("No reminders found for that date. Exiting...")
		return nil
	}

	printRemindersList(reminders)
	count := len(reminders)

	num, err := getIntPrompt("Enter a number to delete the corresponding reminder")
	// handle out of bounds
	if num < 1 || num > count {
		fmt.Println("Invalid entry. Exiting...")
		return nil
	}
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	err = cfg.db.DeleteReminder(context.Background(), reminders[num-1].ID)
	if err != nil {
		return err
	}
	fmt.Println("Reminder deleted")
	return nil

}

func printRemindersList(reminders []database.Reminder) {
	for i, r := range reminders {
		fmt.Printf("* %v: %v -- %s\n", i+1, r.RCcID, r.Note)
	}
}

// prints reminders for currently logged in user on startup
func getTodaysReminders(cfg *Config) error {
	gurParams := database.GetUserTodayRemindersParams{
		RDate:          normalizeDate(time.Now()),
		InvestigatorID: cfg.loggedInInvestigator.ID,
	}

	reminders, err := cfg.db.GetUserTodayReminders(context.Background(), gurParams)
	if err != nil {
		return err
	}

	if len(reminders) == 0 {
		fmt.Println("No reminders found!")
		return nil
	}

	fmt.Println("Reminders for today: ")
	for i, reminder := range reminders {
		fmt.Printf("* %v CC: %v -- %s\n", i+1, reminder.RCcID, reminder.Note)
	}

	return nil
}

// adds a reminder t days after the activation date of the cc (for whoever the cc is under the name of, not who activated it)
func ccActivationReminder(cfg *Config, t int, cc *database.CageCard) error {
	// can't create a reminder without a note
	if !cc.Notes.Valid {
		fmt.Println("Can't create a reminder without a note")
		return nil
	}
	// ugly, has to be a better way
	date := cc.ActivatedOn.Time.Add(time.Hour * (24 * time.Duration(t)))
	arParams := database.AddReminderParams{
		RDate:          date,
		RCcID:          cc.CcID,
		InvestigatorID: cc.InvestigatorID,
		Note:           cc.Notes.String,
	}
	reminder, err := cfg.db.AddReminder(context.Background(), arParams)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(reminder)
	}

	return nil
}
