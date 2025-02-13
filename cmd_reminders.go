package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// just save, exit, help, and prompt everything else
func getAddReminderFlags() map[string]Flag {
	addReminderFlags := make(map[string]Flag)

	// ect as needed or remove the "-"+ for longer ones

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

// look into removing the args thing, might have to stay
// prompt stuff, then print or save it. that's about it.
func addReminderFunction(cfg *Config) error {
	// permission check
	err := checkPermission(cfg.loggedInPosition, PermissionReminders)
	if err != nil {
		return err
	}
	// get flags
	flags := getAddReminderFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// date, cage card, investigator, note
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
	if err != nil {
		return err
	}
	nilCC := database.CageCard{}
	if cc == nilCC {
		fmt.Println("Exiting...")
		return nil
	}

	investigator, err := getStructPrompt(cfg, "Enter investigator who will recieve the reminder", getInvestigatorStruct)
	if err != nil {
		return err
	}
	nilInv := database.Investigator{}
	if investigator == nilInv {
		fmt.Println("Exiting...")
		return nil
	}

	note, err := getStringInput("Enter a note for the reminder")
	if err != nil {
		return err
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

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// just save, exit, help, and prompt everything else
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

func getDatePrompt(prompt string) (time.Time, error) {
	fmt.Println(prompt + " or exit to cancel")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			return time.Time{}, nil
		}

		// then have check if unique or check if not unique after
		output, err := parseDate(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return output, nil

	}
}

func getCageCardStructActive(cfg *Config, input string) (database.CageCard, error) {

	ccid, err := strconv.Atoi(input)
	if err != nil {
		return database.CageCard{}, err
	}
	cc, err := cfg.db.GetCageCardByID(context.Background(), int32(ccid))
	if err != nil {
		return database.CageCard{}, err
	}
	if !cc.ActivatedOn.Valid {
		return database.CageCard{}, errors.New("cage card is not active")
	}
	if cc.DeactivatedOn.Valid {
		return database.CageCard{}, errors.New("cage card is deactivated")
	}

	return cc, nil
}

// do i add the exit on cancel to this one? dunno
func getStringInput(prompt string) (string, error) {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}

		return input, nil
	}

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

// get all the reminders for a certain date.
// list them in whatever order, then enter 0 to delete or # to delete that one
func deleteReminderFunction(cfg *Config) error {
	// permission check
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

	/* removed because this let anybody delete anybody's reminders
	reminders, err = cfg.db.GetAllTodayReminders(context.Background(), date)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}
	*/
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

	printRemindersList(&reminders)
	count := len(reminders)

	num, err := getIntPrompt("Enter a number to delete the corresponding reminder")
	// handle out of bounds
	if num < 1 || num > count {
		fmt.Println("Invalid entry. Exiting...")
		return nil
	}

	if err != nil {
		return err
	}
	if num == -1 {
		fmt.Println("exiting...")
		return nil
	}

	err = cfg.db.DeleteReminder(context.Background(), reminders[num-1].ID)
	if err != nil {
		return err
	}
	fmt.Println("Reminder deleted")
	return nil

}

func printRemindersList(reminders *[]database.Reminder) {
	for i, r := range *reminders {
		fmt.Printf("* %v: %v -- %s\n", i+1, r.RCcID, r.Note)
	}
}

// print reminders for currently logged in user on startup
func getTodaysReminders(cfg *Config) error {
	gurParams := database.GetUserTodayRemindersParams{
		RDate:          normalizeDate(time.Now()),
		InvestigatorID: cfg.loggedInInvestigator.ID,
	}
	// TODO: remember that 0 results wont throw an sql error
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

// adds a reminder X days after the
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
