package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getEditProtocolCmd() Command {
	editProtocolFlags := make(map[string]Flag)
	editProtocolCmd := Command{
		name:        "edit",
		description: "Used for editing an existing protocol",
		function:    editProtocolFunction,
		flags:       editProtocolFlags,
	}

	return editProtocolCmd
}

// [t]itle, [P]I, [N]umber, [A]llocated, [B]alance, [E]xpiration
// save, exit, help
func getEditProtocolFlags() map[string]Flag {
	editProtocolFlags := make(map[string]Flag)
	tFlag := Flag{
		symbol:      "t",
		description: "Changes protocol title",
		takesValue:  true,
	}
	editProtocolFlags["-"+tFlag.symbol] = tFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Changed protocol's PI",
		takesValue:  true,
	}
	editProtocolFlags["-"+pFlag.symbol] = pFlag

	aFlag := Flag{
		symbol:      "a",
		description: "Sets allocated animals",
		takesValue:  true,
	}
	editProtocolFlags["-"+aFlag.symbol] = aFlag

	bFlag := Flag{
		symbol:      "b",
		description: "Changes protocol balance",
		takesValue:  true,
	}
	editProtocolFlags["-"+bFlag.symbol] = bFlag

	eFlag := Flag{
		symbol:      "e",
		description: "Changes expiration date",
		takesValue:  true,
	}
	editProtocolFlags["-"+eFlag.symbol] = eFlag

	// ect as needed or remove the "-"+ for longer ones

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving changes",
		takesValue:  false,
	}
	editProtocolFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags and their uses",
		takesValue:  false,
	}
	editProtocolFlags[helpFlag.symbol] = helpFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves changes made and exits",
		takesValue:  false,
	}
	editProtocolFlags[saveFlag.symbol] = saveFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current changes for review",
		takesValue:  false,
	}
	editProtocolFlags[printFlag.symbol] = printFlag

	return editProtocolFlags

}

// look into removing the args thing, might have to stay
func editProtocolFunction(cfg *Config, args []Argument) error {
	protocol, err := getProtocolByNumber(cfg)
	if err != nil {
		return err
	}
	blankProtocol := database.Protocol{ID: uuid.Nil}
	if protocol == blankProtocol {
		fmt.Println("Exiting...")
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

	// get flags
	flags := getEditProtocolFlags()

	// set defaults
	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Use flags to change protocol parameters. Enter 'help' to see all available flags")
	fmt.Println("When entering values with a space, replace it with an underscore")

	// da loop
	// [t]itle, [P]I, [N]umber, [A]llocated, [B]alance, [E]xpiration
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
		if reviewed.ChangesMade == true {
			reviewed.Printed = false
		}

		// but normal loop now
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
				// updating the pi outside of the block, but will it mess with the error check? glass coffin
				pi, err = getInvestigatorByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				upParams.PrimaryInvestigator = pi.ID
				reviewed.ChangesMade = true
			case "-n":
				if arg.value == protocol.PNumber {
					// same as the OG, so it'll throw an error if trying to change back
					upParams.PNumber = arg.value
					break
				}
				number, err := checkIfProtocolNumberUnique(cfg, arg.value)
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
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
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

func getProtocolByNumber(cfg *Config) (database.Protocol, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the number of the protocol you'd like to edit, or exit to cancel")
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
			return database.Protocol{ID: uuid.Nil}, nil
		}

		protocol, err := cfg.db.GetProtocolByID(context.Background(), input)
		if err != nil && err.Error() != "sql: no rows in result set" {
			// error that isnt related to no rows returned
			fmt.Println("Error checking DB for protocol")
			return database.Protocol{ID: uuid.Nil}, err
		}
		if err != nil && err.Error() == "sql: no rows in result set" {
			fmt.Println("No protocol with that number found. Please try again.")
			continue
		}

		return protocol, nil
	}
}

// [t]itle, [P]I, [N]umber, [A]llocated, [B]alance, [E]xpiration, is [A]ctive
func getInvestigatorByFlag(cfg *Config, i string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), i)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting investigators from DB")
		return database.Investigator{}, err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		fmt.Println("Investigator not found. Please try again")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Vague investigator name. Please try again")
		return database.Investigator{}, nil
	}
	if len(investigators) == 0 {
		fmt.Println("Investigator not found. Please try again")
		return database.Investigator{}, nil
	}
	return investigators[0], nil
}

func checkIfProtocolNumberUnique(cfg *Config, n string) (string, error) {
	_, err := cfg.db.GetProtocolByID(context.Background(), n)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting protocols from DB")
		return "", err
	}
	if err == nil {
		// protocol found
		fmt.Println("Protocol with that number already exists. Please try again")
		return "", err
	}

	// if nothing found, then unique and ok
	return n, nil

}

// look into having it accept one of two generic types as uses same values to print as createProtParam
func printEditProtocol(up *database.UpdateProtocolParams, pi *database.Investigator) {
	fmt.Printf("PI: %s\n", pi.IName)
	fmt.Printf("Number: %s\n", up.PNumber)
	fmt.Printf("Title: %s\n", up.Title)
	fmt.Printf("Allocated: %v\n", up.Allocated)
	fmt.Printf("Expiration Date: %v\n", up.ExpirationDate)
}
