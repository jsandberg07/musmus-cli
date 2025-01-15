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
// PI (skip)
// protocol
// investigator
// pi + investigator (skip)
// protocol + investigator (skip this as well for now)

// then set start and end dates

// bonus ones:
// all active
// ALL cage cards

// write sql for each type of query
// write the go for setting the tags
// write the normalizer functions (have to because psql will probably reset them)
// the data is all the same but the name is different
// so return that, string that, export that
// pass a function first class, normalize response, return, stringify, export
// why do i have a dozen "true activate" "real activate" "honest activate" -_-
// i have literally spend tons of time looking through each step BEFORE this
// gotta plan the whole thing next time lol

// FACT: i didn't make it so i can get them by PI because it created vague returns in postgres so skip it for now lmao
// TODO: get the PI as well in the export
// TODO: get the note in the export
// TODO: reset the protoocl or investigator with like x or whatever (hope nobody is named X lmao)
// TODO: make sure the dates get the correct data, since they're set at like midnigt or something

// [s]tart date, [e]nd date, [pr]otocol, [in]vestigator, print, active, all, help, exit, query
func getCCQueriesFlags() map[string]Flag {
	ccQueriesFlags := make(map[string]Flag)

	sFlag := Flag{
		symbol:      "s",
		description: "Sets start date for query.",
		takesValue:  true,
	}
	ccQueriesFlags["-"+sFlag.symbol] = sFlag

	eFlag := Flag{
		symbol:      "e",
		description: "Sets end date for query.",
		takesValue:  true,
	}
	ccQueriesFlags["-"+eFlag.symbol] = eFlag

	prFlag := Flag{
		symbol:      "pr",
		description: "Gets cards under set protocol. Can either have investigator or protocol",
		takesValue:  true,
	}
	ccQueriesFlags["-"+prFlag.symbol] = prFlag

	inFlag := Flag{
		symbol:      "in",
		description: "Gets cards under set investigator. Can either have investigator or protocol",
		takesValue:  true,
	}
	ccQueriesFlags["-"+inFlag.symbol] = inFlag

	// ect as needed or remove the "-"+ for longer ones

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current settings for query",
		takesValue:  false,
	}
	ccQueriesFlags[printFlag.symbol] = printFlag

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

	queryFlag := Flag{
		symbol:      "query",
		description: "Runs query with current settings",
		takesValue:  false,
	}
	ccQueriesFlags[queryFlag.symbol] = queryFlag

	return ccQueriesFlags

}

// you were working on adding tags, writing a few functions for getting the data
// and a struct normalizer because i'll do that instead of changing the return values

// look into removing the args thing, might have to stay
func CCQueriesFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getCCQueriesFlags()

	// set defaults
	exit := false
	investigator := database.Investigator{}
	protocol := database.Protocol{}

	// might be a problem since data is stored at midnight, might have to do some rounding
	startDate := time.Now()
	endDate := time.Now()

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a protocol or investigator, date range to get all cards active during that time frame")
	fmt.Println("Or 'active' to get all currently active cage cards")

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

		// [s]tart date, [e]nd date, [pr]otocol, [in]vestigator, print, active, all, help, exit, query
		for _, arg := range args {
			switch arg.flag {

			case "-s":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				startDate = date

			case "-e":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				endDate = date

			case "-pr":
				pr, err := getProtocolByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				protocol = pr
				investigator = database.Investigator{}

			case "-in":
				inv, err := getInvestigatorByFlag2(cfg, arg.value)
				if err != nil {
					return err
				}
				investigator = inv
				protocol = database.Protocol{}

			case "help":
				cmdHelp(flags)

			case "active":
				fmt.Println("Getting active cages")
				err := CCQueryActive(cfg)
				if err != nil {
					return err
				}
				exit = true

			case "all":
				fmt.Println("Getting all cages")
				err := CCQueryAll(cfg)
				if err != nil {
					return err
				}
				exit = true

			case "query":
				nilInvestigator := database.Investigator{}
				nilProtocol := database.Protocol{}
				if investigator == nilInvestigator && protocol == nilProtocol {
					fmt.Println("Getting cards active during date range")
					err := CCQueryDateRange(cfg, startDate, endDate)
					if err != nil {
						return err
					}
					break
				}

				if investigator != nilInvestigator {
					fmt.Println("Getting cards active during date range for investigator")
					err := CCQueryInvestigator(cfg, startDate, endDate, &investigator)
					if err != nil {
						return err
					}
					break
				}

				if protocol != nilProtocol {
					fmt.Println("Getting cards active during date range for protocol")
					err := CCQueryProtocol(cfg, startDate, endDate, &protocol)
					if err != nil {
						return err
					}
					break
				}

				fmt.Println("No query run")

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

// keep it simple
// different function for option
// pass in the params
// run the query
// normalize the result
// then pass it to export
// so write a function for each query
// normalize for each struct
// one export

func CCQueryActive(cfg *Config) error {
	ccs, err := cfg.db.GetCageCardsActive(context.Background())
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeActive(&ccs)

	count, err := exportData(&exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

func NormalizeActive(ccs *[]database.GetCageCardsActiveRow) []CageCardExport {
	xps := []CageCardExport{}
	for _, cc := range *ccs {
		txp := CageCardExport{
			CcID:          cc.CcID,
			IName:         cc.IName,
			PNumber:       cc.PNumber,
			SName:         cc.SName,
			ActivatedOn:   cc.ActivatedOn,
			DeactivatedOn: cc.DeactivatedOn,
		}

		xps = append(xps, txp)
	}
	return xps
}

// normalize
// psql generates structs with different names, even if theyre the same
// rather than mess with generated files that might reset each time theyre run
// just have normalize functions that change the returned stucts into a regular type
// also cant be a member function (is that the term if its not oop) since it works on an array

// TODO: candidate for refactoring into one big hideous super function that just branches

func CCQueryAll(cfg *Config) error {
	ccs, err := cfg.db.GetCageCardsAll(context.Background())
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeAll(&ccs)

	count, err := exportData(&exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil
}

func NormalizeAll(ccs *[]database.GetCageCardsAllRow) []CageCardExport {
	xps := []CageCardExport{}
	for _, cc := range *ccs {
		txp := CageCardExport{
			CcID:          cc.CcID,
			IName:         cc.IName,
			PNumber:       cc.PNumber,
			SName:         cc.SName,
			ActivatedOn:   cc.ActivatedOn,
			DeactivatedOn: cc.DeactivatedOn,
		}

		xps = append(xps, txp)
	}
	return xps
}

func CCQueryDateRange(cfg *Config, start, end time.Time) error {
	gcdrParam := database.GetCardsDateRangeParams{
		ActivatedOn:   sql.NullTime{Valid: true, Time: start},
		DeactivatedOn: sql.NullTime{Valid: true, Time: end},
	}
	ccs, err := cfg.db.GetCardsDateRange(context.Background(), gcdrParam)
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeDateRange(&ccs)

	count, err := exportData(&exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil
}

func NormalizeDateRange(ccs *[]database.GetCardsDateRangeRow) []CageCardExport {
	xps := []CageCardExport{}
	for _, cc := range *ccs {
		txp := CageCardExport{
			CcID:          cc.CcID,
			IName:         cc.IName,
			PNumber:       cc.PNumber,
			SName:         cc.SName,
			ActivatedOn:   cc.ActivatedOn,
			DeactivatedOn: cc.DeactivatedOn,
		}

		xps = append(xps, txp)
	}
	return xps
}

func CCQueryInvestigator(cfg *Config, start, end time.Time, inv *database.Investigator) error {
	gciParam := database.GetCageCardsInvestigatorParams{
		ActivatedOn:    sql.NullTime{Valid: true, Time: start},
		DeactivatedOn:  sql.NullTime{Valid: true, Time: end},
		InvestigatorID: inv.ID,
	}

	ccs, err := cfg.db.GetCageCardsInvestigator(context.Background(), gciParam)
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeInvestigator(&ccs)

	count, err := exportData(&exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

func NormalizeInvestigator(ccs *[]database.GetCageCardsInvestigatorRow) []CageCardExport {
	xps := []CageCardExport{}
	for _, cc := range *ccs {
		txp := CageCardExport{
			CcID:          cc.CcID,
			IName:         cc.IName,
			PNumber:       cc.PNumber,
			SName:         cc.SName,
			ActivatedOn:   cc.ActivatedOn,
			DeactivatedOn: cc.DeactivatedOn,
		}

		xps = append(xps, txp)
	}
	return xps
}

func CCQueryProtocol(cfg *Config, start, end time.Time, pro *database.Protocol) error {
	gcpParam := database.GetCageCardsProtocolParams{
		ActivatedOn:   sql.NullTime{Valid: true, Time: start},
		DeactivatedOn: sql.NullTime{Valid: true, Time: end},
		ProtocolID:    pro.ID,
	}

	ccs, err := cfg.db.GetCageCardsProtocol(context.Background(), gcpParam)
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeProtocol(&ccs)

	count, err := exportData(&exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

func NormalizeProtocol(ccs *[]database.GetCageCardsProtocolRow) []CageCardExport {
	xps := []CageCardExport{}
	for _, cc := range *ccs {
		txp := CageCardExport{
			CcID:          cc.CcID,
			IName:         cc.IName,
			PNumber:       cc.PNumber,
			SName:         cc.SName,
			ActivatedOn:   cc.ActivatedOn,
			DeactivatedOn: cc.DeactivatedOn,
		}

		xps = append(xps, txp)
	}
	return xps
}

func exportData(cages *[]CageCardExport) (int, error) {
	err := createExportDirectory()
	if err != nil {
		return 0, err
	}

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
func stringifyCage(c *CageCardExport) []string {

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

// make a constant or make it changable to have the file name be consistent / alterable
func getExportFileName() string {
	uuid := uuid.New().String()
	return "exports/" + "zzz_" + uuid[0:8] + ".csv"
}

func createExportDirectory() error {
	// err := os.Mkdir("exports", os.ModePerm)
	err := os.Mkdir("exports", 0750)
	if err != nil && err.Error() == "mkdir exports: file exists" {
		// it already exists, just skip
		return nil

	}
	if err != nil {
		// any other error
		fmt.Println("Error creating directory")
		return err
	}

	return nil
}

/* this block
is a reminder
of your huberis
and all work
below it is
real but also
fake */

/* because these are all test functions and ive copied what i needed

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
	err := createExportDirectory()
	if err != nil {
		return 0, err
	}
	fmt.Println("Do you see a directory?")

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

*/
