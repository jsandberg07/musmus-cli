package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getAddInvestigatorToProtocolCmd() Command {
	addInvestToProtFlags := make(map[string]Flag)
	addInvProCmd := Command{
		name:        "investigator",
		description: "Used for adding investigators to a protocol",
		function:    addInvestigatorToProtocolFunction,
		flags:       addInvestToProtFlags,
		printOrder:  2,
	}

	return addInvProCmd
}

func getAddInvestToProtFlags() map[string]Flag {
	addInvestToProtFlags := make(map[string]Flag)
	aFlag := Flag{
		symbol:      "-a",
		description: "Use when adding the investigator to the protocol",
		takesValue:  false,
		printOrder:  1,
	}
	addInvestToProtFlags[aFlag.symbol] = aFlag

	rFlag := Flag{
		symbol:      "-r",
		description: "Use when removing the investigator from the protocol",
		takesValue:  false,
		printOrder:  2,
	}
	addInvestToProtFlags[rFlag.symbol] = rFlag

	iFlag := Flag{
		symbol:      "-i",
		description: "Set what investigator to add or remove",
		takesValue:  true,
		printOrder:  3,
	}
	addInvestToProtFlags[iFlag.symbol] = iFlag

	pFlag := Flag{
		symbol:      "-p",
		description: "Set what protocol investigators will be added or removed from",
		takesValue:  true,
		printOrder:  4,
	}
	addInvestToProtFlags[pFlag.symbol] = pFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "exits the current menu",
		takesValue:  false,
		printOrder:  100,
	}
	addInvestToProtFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags for command",
		takesValue:  false,
		printOrder:  100,
	}
	addInvestToProtFlags[helpFlag.symbol] = helpFlag

	return addInvestToProtFlags

}

func addInvestigatorToProtocolFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionProtocol)
	if err != nil {
		return err
	}
	flags := getAddInvestToProtFlags()

	exit := false
	investigator := database.Investigator{ID: uuid.Nil}
	protocol := database.Protocol{ID: uuid.Nil}

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

		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-a":
				if investigator.ID == uuid.Nil {
					fmt.Println("Missing investigator to add")
					break
				}
				if protocol.ID == uuid.Nil {
					fmt.Println("Missing protocol to add investigator ")
					break
				}
				aipParams := database.AddInvestigatorToProtocolParams{
					InvestigatorID: investigator.ID,
					ProtocolID:     protocol.ID,
				}
				addedInvest, err := cfg.db.AddInvestigatorToProtocol(context.Background(), aipParams)

				if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					fmt.Println("Error adding investigator to protocol")
					return err
				}

				if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					fmt.Printf("%s is already added to %s\n", investigator.IName, protocol.PNumber)
					break
				}
				if verbose {
					fmt.Println(addedInvest)
				}
				fmt.Printf("%s added to %s\n", investigator.IName, protocol.PNumber)

			case "-r":
				if investigator.ID == uuid.Nil {
					fmt.Println("Missing investigator to remove")
					break
				}
				if protocol.ID == uuid.Nil {
					fmt.Println("Missing protocol to remove investigator from")
					break
				}
				ripParams := database.RemoveInvestigatorFromProtocolParams{
					InvestigatorID: investigator.ID,
					ProtocolID:     protocol.ID,
				}
				err := cfg.db.RemoveInvestigatorFromProtocol(context.Background(), ripParams)
				if err != nil {
					fmt.Println("Error removing investigator from protocol")
					return err
				}
				// throws no errors if you delete somebody who isn't on something
				fmt.Printf("%s removed from %s\n", investigator.IName, protocol.PNumber)

			case "-i":
				nilInvest := database.Investigator{}
				newInvestigator, err := getInvestigatorByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				if newInvestigator == nilInvest {
					break
				}
				investigator = newInvestigator
				fmt.Printf("Investigator set - %s\n", investigator.IName)

			case "-p":
				nilProtocol := database.Protocol{}
				newProtocol, err := getProtocolByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				if newProtocol == nilProtocol {
					break
				}
				protocol = newProtocol
				fmt.Printf("Protocol set - %s\n", protocol.PNumber)

			case "exit":
				fmt.Println("Exiting...")
				exit = true

			case "help":
				cmdHelp(flags)

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

// investigators need to be added to a protocol before working on it. Ie adding cage cards, activating cage cards, adding orders
func investigatorProtocolCheck(cfg *Config, i *database.Investigator, p *database.Protocol) error {
	cip := database.CheckInvestigatorProtocolParams{
		InvestigatorID: i.ID,
		ProtocolID:     p.ID,
	}
	check, err := cfg.db.CheckInvestigatorProtocol(context.Background(), cip)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// not found
		return errors.New("investigator not on protocol")
	}
	if err != nil {
		// any other error
		return err
	}

	if verbose {
		fmt.Println(check)
	}
	return nil
}
