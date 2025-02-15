package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
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
		printOrder:  3,
	}

	return addCCCmd
}

func getAddCCFlags() map[string]Flag {
	addCCFlags := make(map[string]Flag)
	sFlag := Flag{
		symbol:      "-s",
		description: "Sets start of range (inclusive)",
		takesValue:  true,
		printOrder:  1,
	}
	addCCFlags[sFlag.symbol] = sFlag

	eFlag := Flag{
		symbol:      "-e",
		description: "Sets end of range (inclusive)",
		takesValue:  true,
		printOrder:  2,
	}
	addCCFlags[eFlag.symbol] = eFlag

	aFlag := Flag{
		symbol:      "-a",
		description: "Adds range of cards to database without exiting",
		takesValue:  false,
		printOrder:  3,
	}
	addCCFlags[aFlag.symbol] = aFlag

	iFlag := Flag{
		symbol:      "-i",
		description: "Sets who the cards will be added under",
		takesValue:  true,
		printOrder:  4,
	}
	addCCFlags[iFlag.symbol] = iFlag

	pFlag := Flag{
		symbol:      "-p",
		description: "Sets the protocol the cards will be added under",
		takesValue:  true,
		printOrder:  5,
	}
	addCCFlags[pFlag.symbol] = pFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
		printOrder:  101,
	}
	addCCFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits the current command",
		takesValue:  false,
		printOrder:  102,
	}
	addCCFlags[exitFlag.symbol] = exitFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Adds range of cage cards to database and exits",
		takesValue:  false,
		printOrder:  100,
	}
	addCCFlags[saveFlag.symbol] = saveFlag

	return addCCFlags

}

func addCCFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionActivateInactivate)
	if err != nil {
		return err
	}
	protocol, err := getStructPrompt(cfg, "Enter a protocol for the cards be used with", getProtocolStruct)
	if err != nil && err.Error() != CancelError {
		return err
	}
	if err != nil && err.Error() == CancelError {
		fmt.Println(CancelMsg)
		return nil
	}
	investigator := *cfg.loggedInInvestigator

	err = investigatorProtocolCheck(cfg, &investigator, &protocol)
	if err != nil {
		return err
	}

	flags := getAddCCFlags()

	exit := false
	var start int
	var end int

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Adding cards under logged in user. Change using flags.")
	fmt.Println("Use the flags to set the range of cards to be added to the DB")
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
			case "-s":
				num, err := getNumberFromFlag(arg.value)
				if err != nil {
					fmt.Println(err)
				}
				if num <= 0 {
					fmt.Println("Cage card IDs cannot be negative")
					continue
				}
				start = num

			case "-e":
				num, err := getNumberFromFlag(arg.value)
				if err != nil {
					fmt.Println(err)
				}
				if num <= 0 {
					fmt.Println("Cage card IDs cannot be negative")
					continue
				}
				end = num

			case "-i":
				inv, err := getInvestigatorByFlag2(cfg, arg.value)
				if err != nil {
					return err
				}
				nilInvestigator := database.Investigator{}
				if inv == nilInvestigator {
					break
				}

				err = investigatorProtocolCheck(cfg, &inv, &protocol)
				if err != nil {
					fmt.Println(err)
					continue
				}

				investigator = inv
				fmt.Printf("Investigator set as %s\n", investigator.IName)

			case "-p":
				pro, err := getProtocolByFlag(cfg, arg.value)
				if err != nil {
					return err
				}

				err = investigatorProtocolCheck(cfg, &investigator, &pro)
				if err != nil {
					fmt.Println(err)
					continue
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
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func addCCtoDB(cfg *Config, start, end int, inv database.Investigator, pro database.Protocol) error {
	if end < start {
		fmt.Println("Start of range is larger than end. Please check inputs and try again.")
		return nil
	}

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
