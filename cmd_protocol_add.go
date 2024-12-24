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

func getAddProtocolCmd() Command {
	addProtocolFlags := make(map[string]Flag)
	addProtocolCmd := Command{
		name:        "add",
		description: "Used for adding a new protocol",
		function:    addProtocolFunction,
		flags:       addProtocolFlags,
	}

	return addProtocolCmd
}

// TODO: add editing it before saving it, goes hand in hand with editing but that can wait
// just save and exit in the mean time
func getAddProtocolFlags() map[string]Flag {
	addProtocolFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the new protocol",
		takesValue:  false,
	}
	addProtocolFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	addProtocolFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
	}
	addProtocolFlags[helpFlag.symbol] = helpFlag

	return addProtocolFlags

}

// look into removing the args thing, might have to stay
func addProtocolFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getAddProtocolFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// TODO: not all of these even have errs that can be returned, remove as needed
	pi, err := getNewProtocolPI(cfg)
	if err != nil {
		return err
	}
	title, err := getNewProtocolTitle()
	if err != nil {
		return err
	}
	number, err := getNewProtocolNumber(cfg)
	if err != nil {
		return err
	}
	allocated, err := getNewProtocolAlocated()
	if err != nil {
		return err
	}
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

	// save and exit for now
	// da loop
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

		// do weird behavior here

		// but normal loop now
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
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

// TODO: CURRENTLY 0 restriction on what positions can be a PI
// get everything else going first
// requires a DB change, add a field for like name_on_protocol or is_pi_protocol for position
// just another prompt and tag that should be easy, then test data included

// TODO: enter previous protocol FIRST, copies a lot of work, updates cards
// increasing todo list buddy

// TODO: add exit for when you dont want to make a whole new protocol
func getNewProtocolPI(cfg *Config) (database.Investigator, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the name of the PI overseeing the protocol")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again")
			continue
		}

		investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
		if err != nil && err.Error() != "sql: no rows in result set" {
			// error that isnt related to no rows returned
			fmt.Println("Error checking db for investigator")
			return database.Investigator{ID: uuid.Nil}, err
		}
		if len(investigators) == 0 {
			fmt.Println("No investigator by that name or nickname found. Please try again.")
			continue
		}
		if len(investigators) > 1 {
			fmt.Println("Vague investigator name. Consider entering a nickname instead.")
			continue
		}
		return investigators[0], nil

	}
}

// TODO: add check if protocol by that title already in the DB (number more important really)
func getNewProtocolTitle() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the title of the protocol")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again")
			continue
		}

		return input, nil
	}

}

func getNewProtocolNumber(cfg *Config) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the protocol number")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again")
			continue
		}

		// TODO: this is actually protocol by number, not ID, needs fixing
		protocol, err := cfg.db.GetProtocolByNumber(context.Background(), input)
		if err != nil && err.Error() != "sql: no rows in result set" {
			// any other error than no rows
			fmt.Println("Error getting protocol from DB")
			return "", err
		}
		if err != nil && err.Error() == "sql: no rows in result set" {
			// no rows, which is ideal
			return input, nil
		}
		if err == nil {
			fmt.Println("A protocol with that number already exists. Please try again.")
			fmt.Printf("Protocol with same number: %s\n", protocol.Title)
		}
	}
}

func getNewProtocolAlocated() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the numbers of animals allocated to the protocol")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again")
			continue
		}

		allocated, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error getting number. Be sure to enter an integer.")
			continue
		}

		return allocated, nil
	}
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
