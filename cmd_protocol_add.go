package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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
func getAddProtocolFlags() map[string]Flag {
	addProtocolFlags := make(map[string]Flag)
	XFlag := Flag{
		symbol:      "X",
		description: "Sets X",
		takesValue:  false,
	}
	addProtocolFlags["-"+XFlag.symbol] = XFlag

	// ect as needed or remove the "-"+ for longer ones

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

		for _, arg := range args {
			switch arg.flag {
			case "-X":
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
func getNewProtocolTitle(cfg *Config) (string, error) {
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
		fmt.Println("> ")
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
		protocol, err := cfg.db.GetProtocolByID(context.Background(), input)
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
			continue
		}

		fmt.Println("This message shouldn't appear")

	}

}

/*
func getNewProtocolAlocated(cfg *Config) (int, error) {

}

func getNewProtocolExpiration(cfg *Config) (time.Time, error) {

}

func getPreviousProtocol(cfg *Config) (uuid.UUID, error) {

}
*/
