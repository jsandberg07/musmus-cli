package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

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

// 5 functions
// PI
// protocol
// investigator
// pi + investigator
// protocol + investigator

// then set start and end dates

// bonus ones:
// all active
// ALL cage cards
// these need the same function cause generics are tough

func getCCQueriesFlags() map[string]Flag {
	ccQueriesFlags := make(map[string]Flag)

	// ect as needed or remove the "-"+ for longer ones

	activeFlag := Flag{
		symbol:      "active",
		description: "Exports all currently active cage cards and exits",
		takesValue:  false,
	}
	ccQueriesFlags[activeFlag.symbol] = activeFlag

	allFlag := Flag{
		symbol:      "all",
		description: "Exports all cage cards and exits",
		takesValue:  false,
	}
	ccQueriesFlags[allFlag.symbol] = allFlag

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
				err := exportQuickCC(cfg)
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

func exportQuickCC(cfg *Config) error {
	start := time.Now()
	end := time.Now()

	dates := database.GetCardsDateRangeParams{
		ActivatedOn:   sql.NullTime{Valid: true, Time: start},
		DeactivatedOn: sql.NullTime{Valid: true, Time: end},
	}

	cages, err := cfg.db.GetCardsDateRange(context.Background(), dates)
	if err != nil {
		fmt.Println("Error getting active cages")
		return err
	}

	if len(cages) == 0 {
		fmt.Println("Oops no active cages found!")
		return nil
	}

	fmt.Printf("// expected lines: %v\n", len(cages))

	count, err := exportQuickData(&cages)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

// rename this function, it is no longer a test
func exportActive(cfg *Config) error {
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

func exportQuickData(cages *[]database.GetCardsDateRangeRow) (int, error) {
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
		err := writer.Write(stringifyQuickCage(&cage))
		if err != nil {
			fmt.Printf("Error writing to csv: %s", err)
			continue
		}
		count++
	}

	return count, nil
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

func stringifyQuickCage(c *database.GetCardsDateRangeRow) []string {

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

	return output
}

func getExportFileName() string {
	uuid := uuid.New().String()
	return "zzz_" + uuid[0:8] + ".csv"
}
