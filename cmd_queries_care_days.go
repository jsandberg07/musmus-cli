package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func getCareDaysCmd() Command {
	careDaysFlags := make(map[string]Flag)
	CareDaysCmd := Command{
		name:        "caredays",
		description: "Used for calculating the number of care days during a given period",
		function:    careDaysFunction,
		flags:       careDaysFlags,
		printOrder:  1,
	}

	return CareDaysCmd
}

// [s]tart period, [e]nd period, help, run, toggle printing by [i]nvestigator? [p] protocol too? YEAH DO IT it's like a pivot table
// it's 1 query and and a handful of calculations it's really truly honestly whatever BUT WE CAN ALSO WRITE TESTS HELL YEAH ok now im pumped
// [s]tart, [e]nd, [i]nvestigator, [p]rotocol, query, print, help, exit
func getCareDaysFlags() map[string]Flag {
	CareDaysFlags := make(map[string]Flag)
	sFlag := Flag{
		symbol:      "-s",
		description: "Sets start date of query period",
		takesValue:  true,
		printOrder:  1,
	}
	CareDaysFlags[sFlag.symbol] = sFlag

	eFlag := Flag{
		symbol:      "-e",
		description: "Sets end date of query period",
		takesValue:  true,
		printOrder:  1,
	}
	CareDaysFlags[eFlag.symbol] = eFlag

	/* TODO: add the ability to get # of care days by investigator or protocol
	iFlag := Flag{
		symbol:      "i",
		description: "Toggles reporting by investigator",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags["-"+iFlag.symbol] = iFlag


	pFlag := Flag{
		symbol:      "p",
		description: "Toggles reporting by protocol",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags["-"+pFlag.symbol] = pFlag
	*/

	queryFlag := Flag{
		symbol:      "query",
		description: "Runs query with current parameters",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags[queryFlag.symbol] = queryFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current query parameters",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags[printFlag.symbol] = printFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available flags",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits command",
		takesValue:  false,
		printOrder:  1,
	}
	CareDaysFlags[exitFlag.symbol] = exitFlag

	return CareDaysFlags

}

// TODO: needs a way to do something more than just give a number. Should be able to calc for each person or protocol or both. Probably
// just one or the other to start
func careDaysFunction(cfg *Config) error {
	// get flags
	flags := getCareDaysFlags()

	// set defaults
	exit := false
	start := normalizeDate(time.Now())
	end := normalizeDate(time.Now())
	/* TODO: add the ability to get # of care days by investigator or protocol
	investigatorReport := false
	protocolReport := false
	*/

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
	fmt.Println("Enter 'help' to see a list of parameters available to calculate care days")
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

		// [s]tart, [e]nd, [i]nvestigator, [p]rotocol, query, print, help, exit
		for _, arg := range args {
			switch arg.flag {
			case "-s":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					continue
				}
				start = normalizeDate(date)

				if end.Before(start) {
					fmt.Println("End date before start. Please double check before running query.")
				}

			case "-e":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					continue
				}
				end = normalizeDate(date)

				if end.Before(start) {
					fmt.Println("End date before start. Please double check before running query.")
				}

				/* TODO: add the ability to get # of care days by investigator or protocol
				case "-i":
					investigatorReport = !investigatorReport

				case "-p":
					protocolReport = !protocolReport
				*/

			case "query":
				/*
					gcdrp := database.GetCardsDateRangeParams{
						ActivatedOn:   sql.NullTime{Valid: true, Time: start},
						DeactivatedOn: sql.NullTime{Valid: true, Time: end},
					}
				*/
				gcdrp := database.GetCardsDateRangeParams{
					Overlaps:   sql.NullTime{Valid: true, Time: start},
					Overlaps_2: sql.NullTime{Valid: true, Time: end},
				}
				cageCards, err := cfg.db.GetCardsDateRange(context.Background(), gcdrp)

				if err != nil {
					fmt.Println("Error running query")
					return err
				}

				if len(cageCards) == 0 {
					fmt.Println("No cage cards found")
					continue
				}

				num := careDaysQuery(start, end, cageCards)
				fmt.Printf("Sum of care days - %v\n", num)

			case "print":
				fmt.Printf("Start date - %v\n", start)
				fmt.Printf("End  date  - %v\n", end)

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

// get the cage cards, then pass it into a function instead of being tied into once because testingu
// uhh its a test i guess print what you want until you get it. how do i run a test on one file?
// also consistent! sort by name or whatever
// how do we sort and print these
// sort and print, but return total care days regardless
func careDaysQuery(start, end time.Time, ccs []database.GetCardsDateRangeRow) int {
	var num time.Duration
	oneDay := 24 * time.Hour
	fmt.Printf("// Len of CCs: %v\n", len(ccs))
	for i, cc := range ccs {
		if cc.ActivatedOn.Time.Before(start) {
			cc.ActivatedOn.Time = start
		}
		if cc.DeactivatedOn.Time.After(end) {
			cc.DeactivatedOn.Time = end
		}
		if !cc.DeactivatedOn.Valid {
			cc.DeactivatedOn.Valid = true
			cc.DeactivatedOn.Time = end
		}
		ccs[i] = cc
		num += cc.DeactivatedOn.Time.Sub(cc.ActivatedOn.Time) + oneDay
	}
	total := (num.Hours() / 24)
	return int(total)
}
