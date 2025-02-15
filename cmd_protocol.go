package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddProtocolCmd() Command {
	addProtocolFlags := make(map[string]Flag)
	addProtocolCmd := Command{
		name:        "add",
		description: "Used for adding a new protocol",
		function:    addProtocolFunction,
		flags:       addProtocolFlags,
		printOrder:  1,
	}

	return addProtocolCmd
}

func getAddProtocolFlags() map[string]Flag {
	addProtocolFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the new protocol",
		takesValue:  false,
		printOrder:  100,
	}
	addProtocolFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
		printOrder:  100,
	}
	addProtocolFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
		printOrder:  100,
	}
	addProtocolFlags[helpFlag.symbol] = helpFlag

	return addProtocolFlags

}

func addProtocolFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionProtocol)
	if err != nil {
		return err
	}
	flags := getAddProtocolFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)

	pi, err := getStructPrompt(cfg, "Enter the name of the PI overseeing the protocol", getPIStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	title, err := getStringPrompt(cfg, "Enter title of the new protocol", checkFuncNil)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	number, err := getStringPrompt(cfg, "Enter number of new protocol", checkProtocolUnique)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	allocated, err := getIntPrompt("Enter the numbers of animals allocated to the protocol")
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	// don't use prompts for this one, defaults to three years from today
	expiration, err := getNewProtocolExpiration()
	if err != nil {
		return err
	}

	cpParam := database.CreateProtocolParams{
		PNumber:             number,
		PrimaryInvestigator: pi.ID,
		Title:               title,
		Allocated:           int32(allocated),
		Balance:             int32(0),
		ExpirationDate:      expiration,
	}

	fmt.Println("Current info:")
	printAddProtocol(cpParam, pi)
	fmt.Println("Enter 'save' to save, 'exit' to exit without saving")
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
			case "help":
				cmdHelp(flags)
			case "save":
				fmt.Println("Saving...")
				protocol, err := cfg.db.CreateProtocol(context.Background(), cpParam)
				if err != nil {
					fmt.Println("Error saving protocol")
					return err
				}
				exit = true
				if verbose {
					fmt.Println(protocol)
				}
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

func getNewProtocolExpiration() (time.Time, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter expiration date of the protocol")
	fmt.Println("Enter nothing to set it to 3 years from today")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			today := time.Now()
			then := today.AddDate(3, 0, 0)
			return then, nil
		}

		expirationDate, err := parseDate(input)
		if err != nil {
			continue
		}

		return expirationDate, nil
	}

}

func printAddProtocol(cp database.CreateProtocolParams, pi database.Investigator) {
	fmt.Printf("PI: %s\n", pi.IName)
	fmt.Printf("Number: %s\n", cp.PNumber)
	fmt.Printf("Title: %s\n", cp.Title)
	fmt.Printf("Allocated: %v\n", cp.Allocated)
	fmt.Printf("Expiration Date: %v\n", cp.ExpirationDate)
}

func getEditProtocolCmd() Command {
	editProtocolFlags := make(map[string]Flag)
	editProtocolCmd := Command{
		name:        "edit",
		description: "Used for editing an existing protocol",
		function:    editProtocolFunction,
		flags:       editProtocolFlags,
		printOrder:  3,
	}

	return editProtocolCmd
}

func getEditProtocolFlags() map[string]Flag {
	editProtocolFlags := make(map[string]Flag)
	tFlag := Flag{
		symbol:      "-t",
		description: "Changes protocol title",
		takesValue:  true,
		printOrder:  1,
	}
	editProtocolFlags[tFlag.symbol] = tFlag

	pFlag := Flag{
		symbol:      "-p",
		description: "Changed protocol's PI",
		takesValue:  true,
		printOrder:  2,
	}
	editProtocolFlags[pFlag.symbol] = pFlag

	aFlag := Flag{
		symbol:      "-a",
		description: "Sets allocated animals",
		takesValue:  true,
		printOrder:  3,
	}
	editProtocolFlags[aFlag.symbol] = aFlag

	bFlag := Flag{
		symbol:      "-b",
		description: "Changes protocol balance",
		takesValue:  true,
		printOrder:  4,
	}
	editProtocolFlags[bFlag.symbol] = bFlag

	eFlag := Flag{
		symbol:      "-e",
		description: "Changes expiration date",
		takesValue:  true,
		printOrder:  5,
	}
	editProtocolFlags[eFlag.symbol] = eFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving changes",
		takesValue:  false,
		printOrder:  100,
	}
	editProtocolFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags and their uses",
		takesValue:  false,
		printOrder:  100,
	}
	editProtocolFlags[helpFlag.symbol] = helpFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves changes made and exits",
		takesValue:  false,
		printOrder:  99,
	}
	editProtocolFlags[saveFlag.symbol] = saveFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current changes for review",
		takesValue:  false,
		printOrder:  99,
	}
	editProtocolFlags[printFlag.symbol] = printFlag

	return editProtocolFlags

}

func editProtocolFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionProtocol)
	if err != nil {
		return err
	}
	protocol, err := getStructPrompt(cfg, "Enter number of protocol to edit", checkProtocolExists)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}

	upParams := database.UpdateProtocolParams{
		ID:                  protocol.ID,
		PNumber:             protocol.PNumber,
		PrimaryInvestigator: protocol.PrimaryInvestigator,
		Title:               protocol.Title,
		Allocated:           protocol.Allocated,
		Balance:             protocol.Balance,
		ExpirationDate:      protocol.ExpirationDate,
	}
	pi, err := cfg.db.GetInvestigatorByID(context.Background(), protocol.PrimaryInvestigator)
	if err != nil {
		fmt.Println("Error getting PI for protocol")
		return err
	}

	flags := getEditProtocolFlags()

	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Use flags to change protocol parameters. Enter 'help' to see all available flags")
	fmt.Println("When entering values with a space, replace it with an underscore")

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

		if reviewed.ChangesMade {
			reviewed.Printed = false
		}

		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-t":
				upParams.Title = arg.value
				reviewed.ChangesMade = true

			case "-p":
				pi, err = getInvestigatorByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				upParams.PrimaryInvestigator = pi.ID
				reviewed.ChangesMade = true

			case "-n":
				if arg.value == protocol.PNumber {
					upParams.PNumber = arg.value
					break
				}
				number, err := getUniqueProtocolFromFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				upParams.PNumber = number
				reviewed.ChangesMade = true

			case "-a":
				allocated, err := strconv.Atoi(arg.value)
				if err != nil {
					fmt.Printf("Error updating allocated animals: %s\n", err)
					break
				}
				upParams.Allocated = int32(allocated)
				reviewed.ChangesMade = true

			case "-b":
				balance, err := strconv.Atoi(arg.value)
				if err != nil {
					fmt.Printf("Error updating protocol balance: %s\n", err)
					break
				}
				upParams.Balance = int32(balance)
				reviewed.ChangesMade = true

			case "-e":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				if time.Now().After(date) {
					fmt.Println("New date is after today, meaning protocol is expired")
					fmt.Println("Change will be made but please double check input")
				}
				upParams.ExpirationDate = date
				reviewed.ChangesMade = true

			case "help":
				cmdHelp(flags)

			case "print":
				printEditProtocol(&upParams, &pi)
				reviewed.ChangesMade = false
				reviewed.Printed = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true

			case "save":
				fmt.Println("Saving...")
				err := cfg.db.UpdateProtocol(context.Background(), upParams)
				if err != nil {
					fmt.Println("Error updating protocol")
					return err
				}
				exit = true
			default:
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if upParams.Balance > upParams.Allocated {
			fmt.Println("Protocol balance exceeds allocated animals.")
			fmt.Println("Please double check these numbers, as this would mean the protocol is in compliance")
		}

		if exit {
			break
		}

	}

	return nil
}

func printEditProtocol(up *database.UpdateProtocolParams, pi *database.Investigator) {
	fmt.Printf("PI: %s\n", pi.IName)
	fmt.Printf("Number: %s\n", up.PNumber)
	fmt.Printf("Title: %s\n", up.Title)
	fmt.Printf("Allocated: %v\n", up.Allocated)
	fmt.Printf("Expiration Date: %v\n", up.ExpirationDate)
}

func addBalanceToProtocol(cfg *Config, input int, cc *database.CageCard) error {
	abParams := database.AddBalanceParams{
		ID:      cc.ProtocolID,
		Balance: int32(input),
	}
	err := cfg.db.AddBalance(context.Background(), abParams)
	if err != nil {
		return err
	}

	protocol, err := cfg.db.GetProtocolByID(context.Background(), cc.ProtocolID)
	if err != nil {
		return err
	}
	if protocol.Allocated < protocol.Balance {
		fmt.Println("Animals used on protocol exceeds allotment")
		fmt.Printf("Allocated - %v \nBalance - %v\n", protocol.Allocated, protocol.Balance)
	}

	return nil
}
