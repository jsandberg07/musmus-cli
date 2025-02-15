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
		printOrder:  1,
	}

	return CCQueriesCmd
}
func getCCQueriesFlags() map[string]Flag {
	ccQueriesFlags := make(map[string]Flag)

	sFlag := Flag{
		symbol:      "-s",
		description: "Sets start date for query.",
		takesValue:  true,
		printOrder:  2,
	}
	ccQueriesFlags[sFlag.symbol] = sFlag

	eFlag := Flag{
		symbol:      "-e",
		description: "Sets end date for query.",
		takesValue:  true,
		printOrder:  3,
	}
	ccQueriesFlags[eFlag.symbol] = eFlag

	prFlag := Flag{
		symbol:      "-pr",
		description: "Gets cards under set protocol.",
		takesValue:  true,
		printOrder:  4,
	}
	ccQueriesFlags[prFlag.symbol] = prFlag

	inFlag := Flag{
		symbol:      "-in",
		description: "Gets cards under set investigator.",
		takesValue:  true,
		printOrder:  5,
	}
	ccQueriesFlags[inFlag.symbol] = inFlag

	orFlag := Flag{
		symbol:      "-or",
		description: "Gets cards that were added under set order",
		takesValue:  true,
		printOrder:  6,
	}
	ccQueriesFlags[orFlag.symbol] = orFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current settings for query",
		takesValue:  false,
		printOrder:  7,
	}
	ccQueriesFlags[printFlag.symbol] = printFlag

	activeFlag := Flag{
		symbol:      "active",
		description: "Exports all currently active cage cards and exits",
		takesValue:  false,
		printOrder:  0,
	}
	ccQueriesFlags[activeFlag.symbol] = activeFlag

	allFlag := Flag{
		symbol:      "all",
		description: "Exports all cage cards and exits",
		takesValue:  false,
		printOrder:  1,
	}
	ccQueriesFlags[allFlag.symbol] = allFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints all available flags",
		takesValue:  false,
		printOrder:  100,
	}
	ccQueriesFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without exporting",
		takesValue:  false,
		printOrder:  100,
	}
	ccQueriesFlags[exitFlag.symbol] = exitFlag

	queryFlag := Flag{
		symbol:      "query",
		description: "Runs query with current settings",
		takesValue:  false,
		printOrder:  99,
	}
	ccQueriesFlags[queryFlag.symbol] = queryFlag

	return ccQueriesFlags

}

func CCQueriesFunction(cfg *Config) error {
	err := checkPermission(cfg.loggedInPosition, PermissionRunQueries)
	if err != nil {
		return err
	}

	flags := getCCQueriesFlags()

	exit := false
	investigator := database.Investigator{}
	protocol := database.Protocol{}
	order := database.Order{}

	startDate := normalizeDate(time.Now())
	endDate := normalizeDate(time.Now())

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter a protocol or investigator, date range to get all cards active during that time frame")
	fmt.Println("Or 'active' to get all currently active cage cards")
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
				nilProtocol := database.Protocol{}
				if pr == nilProtocol {
					break
				}

				protocol = pr
				investigator = database.Investigator{}
				order = database.Order{}

			case "-in":
				inv, err := getInvestigatorByFlag2(cfg, arg.value)
				if err != nil {
					return err
				}
				nilInvestigator := database.Investigator{}
				if inv == nilInvestigator {
					break
				}

				investigator = inv
				protocol = database.Protocol{}
				order = database.Order{}

			case "-or":
				or, err := getOrderByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				nilOrder := database.Order{}
				if or == nilOrder {
					break
				}

				if !or.Received {
					fmt.Println("Order has not been recieved so no cage cards are associated with it")
					break
				}

				investigator = database.Investigator{}
				protocol = database.Protocol{}
				order = or

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
				nilOrder := database.Order{}
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

				if order != nilOrder {
					fmt.Println("Getting cage cards for order")
					err := CCQueryOrder(cfg, &order)
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
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func CCQueryActive(cfg *Config) error {
	ccs, err := cfg.db.GetCageCardsActive(context.Background())
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

func CCQueryAll(cfg *Config) error {
	ccs, err := cfg.db.GetCageCardsAll(context.Background())
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil
}

func CCQueryDateRange(cfg *Config, start, end time.Time) error {
	gcdrParam := database.GetCardsDateRangeParams{
		Overlaps:   sql.NullTime{Valid: true, Time: start},
		Overlaps_2: sql.NullTime{Valid: true, Time: end},
	}
	ccs, err := cfg.db.GetCardsDateRange(context.Background(), gcdrParam)
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil
}

func CCQueryInvestigator(cfg *Config, start, end time.Time, inv *database.Investigator) error {
	gciParam := database.GetCageCardsInvestigatorParams{
		Overlaps:       sql.NullTime{Valid: true, Time: start},
		Overlaps_2:     sql.NullTime{Valid: true, Time: end},
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

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

func CCQueryOrder(cfg *Config, o *database.Order) error {
	ccs, err := cfg.db.GetCageCardsOrder(context.Background(), uuid.NullUUID{Valid: true, UUID: o.ID})
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil
}

func CCQueryProtocol(cfg *Config, start, end time.Time, pro *database.Protocol) error {
	gcpParam := database.GetCageCardsProtocolParams{
		Overlaps:   sql.NullTime{Valid: true, Time: start},
		Overlaps_2: sql.NullTime{Valid: true, Time: end},
		ProtocolID: pro.ID,
	}

	ccs, err := cfg.db.GetCageCardsProtocol(context.Background(), gcpParam)
	if err != nil {
		return err
	}

	if len(ccs) == 0 {
		fmt.Println("No cage cards found!")
		return nil
	}

	exp := NormalizeCCExport(ccs)

	count, err := exportData(exp)
	if err != nil {
		return err
	}

	fmt.Printf("// Exported lines: %v\n", count)
	return nil

}

// Output rows are consistent between types, but sqlc generates new struct for each query.
// Changes them into a format that can be turned into a string to be exported.
func NormalizeCCExport[T database.GetCageCardsInvestigatorRow |
	database.GetCageCardsAllRow |
	database.GetCageCardsActiveRow |
	database.GetCardsDateRangeRow |
	database.GetCageCardsProtocolRow |
	database.GetCageCardsOrderRow](ccs []T) []CageCardExport {
	if len(ccs) == 0 {
		return []CageCardExport{}
	}
	output := make([]CageCardExport, len(ccs))
	for i, cc := range ccs {
		ts := CageCardExport(cc)
		output[i] = ts
	}
	return output
}

func exportData(cages []CageCardExport) (int, error) {
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
	topRow := []string{"CC", "Investigator", "Protocol", "Strain", "Activated On", "Deactivated On", "Order Number"}
	err = writer.Write(topRow)
	if err != nil {
		fmt.Println("Error writing top row to csv")
		return 0, err
	}

	for _, cage := range cages {
		err := writer.Write(stringifyCage(&cage))
		if err != nil {
			fmt.Printf("Error writing to csv: %s", err)
			continue
		}
		count++
	}

	return count, nil
}

func stringifyCage(c *CageCardExport) []string {

	output := make([]string, 7)
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
	if c.OrderNumber.Valid {
		output[6] = c.OrderNumber.String
	}

	return output

}

func getExportFileName() string {
	uuid := uuid.New().String()
	return "exports/" + uuid[0:8] + ".csv"
}

func createExportDirectory() error {
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
