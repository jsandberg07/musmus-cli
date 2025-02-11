package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jsandberg07/clitest/internal/database"
)

func getAddCCCmd() Command {
	addCCFlags := make(map[string]Flag)
	addCCCmd := Command{
		name:        "add",
		description: "Used to add a range of unactivated cage cards",
		function:    addCCFunction,
		flags:       addCCFlags,
	}

	return addCCCmd
}

// [s]tart, [e]nd, [a]dd, [i]nvestigator, [p]rotocol, help, exit, save
// have to handle dupes as the input int is
func getAddCCFlags() map[string]Flag {
	addCCFlags := make(map[string]Flag)
	sFlag := Flag{
		symbol:      "s",
		description: "Sets start of range (inclusive)",
		takesValue:  true,
	}
	addCCFlags["-"+sFlag.symbol] = sFlag

	eFlag := Flag{
		symbol:      "e",
		description: "Sets end of range (inclusive)",
		takesValue:  true,
	}
	addCCFlags["-"+eFlag.symbol] = eFlag

	aFlag := Flag{
		symbol:      "a",
		description: "Adds range of cards to database without exiting",
		takesValue:  false,
	}
	addCCFlags["-"+aFlag.symbol] = aFlag

	iFlag := Flag{
		symbol:      "i",
		description: "Sets who the cards will be added under",
		takesValue:  true,
	}
	addCCFlags["-"+iFlag.symbol] = iFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Sets the protocol the cards will be added under",
		takesValue:  true,
	}
	addCCFlags["-"+pFlag.symbol] = pFlag

	// ect as needed or remove the "-"+ for longer ones

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	addCCFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits the current command",
		takesValue:  false,
	}
	addCCFlags[exitFlag.symbol] = exitFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Adds range of cage cards to database and exits",
		takesValue:  false,
	}
	addCCFlags[saveFlag.symbol] = saveFlag

	return addCCFlags

}

// look into removing the args thing, might have to stay
func addCCFunction(cfg *Config) error {
	var nilProtocol database.Protocol
	protocol, err := getStructPrompt(cfg, "Enter a protocol for the cards be used with", getProtocolStruct)
	if err != nil {
		return err
	}
	if protocol == nilProtocol {
		fmt.Println("Exiting...")
		return nil
	}
	investigator := *cfg.loggedInInvestigator

	// get flags
	flags := getAddCCFlags()

	// set defaults
	exit := false
	var start int
	var end int

	// the reader
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Adding cards under logged in user. Change using flags.")
	fmt.Println("Use the flags to set the range of cards to be added to the DB")
	// fmt.Println("Enter help to see available flags")
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

		// [s]tart, [e]nd, [a]dd, [i]nvestigator, [p]rotocol, help, exit, save
		for _, arg := range args {
			switch arg.flag {
			case "-s":
				num, err := getNumberFromFlag(arg.value)
				if err != nil {
					fmt.Println(err)
				}
				start = num

			case "-e":
				num, err := getNumberFromFlag(arg.value)
				if err != nil {
					fmt.Println(err)
				}
				end = num

			case "-i":
				inv, err := getInvestigatorByFlag2(cfg, arg.value)
				if err != nil {
					return err
				}
				investigator = inv
				fmt.Printf("Investigator set as %s\n", investigator.IName)

			case "-p":
				pro, err := getProtocolByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				protocol = pro

			case "save":
				fmt.Println("Saving and exiting...")
				exit = true
				fallthrough
			case "-a":
				err := addCCtoDB(cfg, start, end, investigator, protocol)
				if err != nil {
					return err
				}

			case "help":
				cmdHelp(flags)

			case "exit":
				fmt.Println("Exiting...")
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

// kind of a place holder. TODO: decide what to do with error
func getNumberFromFlag(input string) (int, error) {
	if input == "" {
		fmt.Println("No input found. Please try again.")
		return 0, nil
	}
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Could not read number from input")
		return 0, err
	}
	return num, nil
}

// possibly go routine this later, for now should be fast enough
func addCCtoDB(cfg *Config, start, end int, inv database.Investigator, pro database.Protocol) error {
	if end < start {
		fmt.Println("Start of range is larger than end. Please check inputs and try again.")
		return nil
	}

	// maybe remove
	fmt.Printf("Adding %v cards\n", (end-start)+1)
	added := 0
	// inclusive. Not interating an array. The <= is intended.
	for i := start; i <= end; i++ {
		accParams := database.AddCageCardParams{
			CcID:           int32(i),
			ProtocolID:     pro.ID,
			InvestigatorID: inv.ID,
		}
		cc, err := cfg.db.AddCageCard(context.Background(), accParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("%v has already been added to the DB\n", i)
				continue
			} else {
				fmt.Println("Error adding cage card to DB")
				fmt.Println(err)
				continue
			}
		}
		if verbose {
			fmt.Println(cc)
		}
		added++
	}
	fmt.Printf("%v cards added!\n", added)
	return nil
}

// literally rewrote this by accident so here's a smaller version, will decide on what to use later
func getInvestigatorByFlag2(cfg *Config, input string) (database.Investigator, error) {
	if input == "" {
		fmt.Println("No input found, please try again")
		return database.Investigator{}, nil
	}
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
	if err != nil {
		fmt.Println("Error getting investigator from db")
		return database.Investigator{}, err
	}
	if len(investigators) == 0 {
		fmt.Println("No investigator by that name found. Nicknames also work as well.")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Vague investigator name. Please try again.")
		return database.Investigator{}, nil
	}
	return investigators[0], nil
}

func getProtocolStruct(cfg *Config, input string) (database.Protocol, error) {
	protocol, err := cfg.db.GetProtocolByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Protocol{}, errors.New("protocol not found. please try again")
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting strain from DB.")
		return database.Protocol{}, err
	}

	return protocol, nil
}
