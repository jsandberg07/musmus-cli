package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getCCQueriesCmd() Command {
	CCQueriesFlags := make(map[string]Flag)
	CCQueriesCmd := Command{
		name:        "cc",
		description: "Run queries on cage cards",
		function:    CCQueriesFunction,
		flags:       CCQueriesFlags,
	}

	return CCQueriesCmd
}

// dag how do i wanna do this
// thinking is king innit?
// maybe prompts
// like quick ones >all active cards
// or limit it like just by PI or ivnestigator. its a csv
// easier to short than write and parse a bunch of unique queries
// problem is there are so many not nulls that need to have a match
// it would be easy if protocol could be null
// or investigator could be null, then it could be optional that way
// param structs
// where its like investigator, protocol, activated on, deactivated on
// then have a bunch of cases for which query to run
// no no
// investigator, protocol, PI are optional. you set and case those.
// check if that protocol is under that PI and eliminate one query
// dont do that for investigator, in case they are removed from that protocol
// then you can do activated on, deactivated as optional like null or not null
// depending on if you want active, all, deactivated
// how to do for particular dates? fuck that, that's all narrow it down later on your own time
// fine this is ok

// can you parse greater than it dates? probably
// uhh for now get ALL cage cards, then export to CSV in a folder with an UUID
// then you can worry about parameters

func getCCQueriesFlags() map[string]Flag {
	ccQueriesFlags := make(map[string]Flag)

	// ect as needed or remove the "-"+ for longer ones

	activeFlag := Flag{
		symbol:      "active",
		description: "Exports all currently active cage cards and exits",
		takesValue:  false,
	}
	ccQueriesFlags[activeFlag.symbol] = activeFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
	}
	ccQueriesFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without exporting",
		takesValue:  false,
	}
	ccQueriesFlags[exitFlag.symbol] = exitFlag

	return ccQueriesFlags

}

// look into removing the args thing, might have to stay
func CCQueriesFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getCCQueriesFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Just enter active or exit")

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

			case "help":
				cmdHelp(flags)

			case "active":
				fmt.Println("Getting active cages")
				err := exportActive(cfg)
				if err != nil {
					return err
				}
				exit = true

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

func exportActive(cfg *Config) error {
	// JOINS
	// you want the cc#, protocol, activated and deactivated dates, investigator, strain, notes
	// fk to protocol for number, fk to investigator for name, fk protocol to pi for name,
	// need to readd PI name probably, is ambiguious currently
	// or dont! just know. it's unlikely a lab would have more than 1 anyway.
	activeCages, err := cfg.db.GetActiveTestCards(context.Background())
	if err != nil {
		fmt.Println("Error getting active cages")
		return err
	}

	if len(activeCages) == 0 {
		fmt.Println("Oops no active cages found!")
		return nil
	}

	fmt.Printf("// expected lines: %v\n", len(activeCages))

	count, err := exportData(&activeCages)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

// all exported data is more or less the same, so this can probably be made generic
func exportData(cages *[]database.GetActiveTestCardsRow) (int, error) {
	filename := getExportFileName()
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating csv file")
		return 0, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	count := 0

	// add top row
	topRow := []string{"CC", "Investigator", "Protocol", "Strain", "Activated On", "Deactivated On"}
	err = writer.Write(topRow)
	if err != nil {
		fmt.Println("Error writing top row to csv")
		return 0, err
	}

	for _, cage := range *cages {
		err := writer.Write(stringifyCage(&cage))
		if err != nil {
			fmt.Printf("Error writing to csv: %s", err)
			continue
		}
		count++
	}

	return count, nil
}

// TODO: format the dates so theyre just like a day and not a stinkin millisecond
func stringifyCage(c *database.GetActiveTestCardsRow) []string {

	output := make([]string, 6)
	output[0] = strconv.Itoa(int(c.CcID))

	output[1] = c.IName

	output[2] = c.PNumber

	if c.SName.Valid {
		output[3] = c.SName.String
	}

	if c.ActivatedOn.Valid {
		output[4] = c.ActivatedOn.Time.String()
	}
	if c.DeactivatedOn.Valid {
		output[5] = c.DeactivatedOn.Time.String()
	}

	/* going step by step to make sure the sql works so adding a table at a time
	if c.Investigator.Valid {
		output[1] = c.Investigator.String
	}
	if c.PNumber.Valid {
		output[2] = c.PNumber.String
	}
	if c.SName.Valid {
		output[3] = c.SName.String
	}
	if c.ActivatedOn.Valid {
		output[4] = c.ActivatedOn.Time.String()
	}
	if c.DeactivatedOn.Valid {
		output[5] = c.DeactivatedOn.Time.String()
	}
	*/

	return output

}

func getExportFileName() string {
	uuid := uuid.New().String()
	return "zzz_" + uuid[0:8] + ".csv"
}
