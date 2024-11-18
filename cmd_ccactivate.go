package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// uhh how do i want this to work
// you either add a number just with a number
// or you set the date
// try to parse to int and if that fails read the other stuff
// for now have some confirmation
// no that doesnt work cause its in the main loop FUCK
// print some stuff and come back to this
// have a CC loop instead, then you CALL activate or deact with params
// new loop that just parses #s
// then you type process or something and it runs them all
// then returns to CC loop
// yeah

func getActivateCmd() Command {
	activateFlags := make(map[string]Flag)

	dFlag := Flag{
		symbol:      "d",
		description: "Sets date. Defaults to today.",
		takesValue:  true,
	}
	activateFlags["-"+dFlag.symbol] = dFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Sets person. Defaults to 'Mouse.'",
		takesValue:  true,
	}
	activateFlags["-"+pFlag.symbol] = pFlag

	activateCmd := Command{
		name:        "activate",
		description: "Used for activating cards.",
		function:    activateCommand,
		flags:       activateFlags,
	}

	return activateCmd
}

// process itself
// TODO: Add a delete that just pops one of the cards off
func activateCommand(cfg *Config, args []Argument) error {
	var date time.Time
	var person string
	var err error // because blocks

	for _, argument := range args {
		switch argument.flag {
		case "-d":
			date, err = parseDate(argument.value)
			if err != nil {
				return errors.New("Couldn't parse date. Try this format: month/day/year.")
			}
		case "-p":
			// when db, add to make sure somebody is in the db
			person = argument.value
		}
	}

	if date.IsZero() {
		year, month, day := time.Now().Date()
		date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}

	if person == "" {
		person = "Mouse"
	}

	fmt.Printf("Enter cage cards to activate. Date: %v - Person: %s\n", date, person)
	fmt.Println("Enter 'process' to process the cards and exit, or exit to leave without processing.")
	reader := bufio.NewReader(os.Stdin)
	var cardsToProcess []CageCard
	exit := false

	for true {
		fmt.Print(">")
		text, err := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "process" {
			exit = false
			break
		}
		if text == "exit" {
			exit = true
			break
		}

		if err != nil {
			fmt.Println("Error reading input string")
			return err
		}
		ccid, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid input. Expecting a number or 'process'.")
			continue
		}

		tcc := CageCard{
			CCid:   ccid,
			Date:   date,
			Person: person,
		}

		cardsToProcess = append(cardsToProcess, tcc)

	}

	if exit == true {
		fmt.Println("Exiting without processing cards.")
		return nil
	}

	var wg sync.WaitGroup

	if len(cardsToProcess) != 0 {
		// fmt.Println("Here's where cards would be processed :^3")
		for _, cc := range cardsToProcess {
			wg.Add(1)
			go func() {
				defer wg.Done()
				processCard(&cc)
			}()
		}
	}

	wg.Wait()

	return nil

}

// because who knows what dates people are going to enter
// TODO: make sure there isnt fuckery like "this card technically wasn't this day beacuase the time was off"
// set everything to be active at like midnight, queries at midnight and see if it works otherwise +1 second lmao

func parseDate(input string) (time.Time, error) {
	// create an array of the formats (with 0s, without, 4 digit year, 2 digit year)
	// go through parse works and then return
	var date time.Time
	var err error
	timeFormats := []string{"1/2/06", "1/2/2006", "01/02/06", "01/02/2006"}
	for _, format := range timeFormats {
		date, err = time.Parse(format, input)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("Error parsing date.")
		return time.Time{}, err
	}

	return date, nil
}

// go routine for just printing cards, with a DB it'll be a sql thing
func processCard(cc *CageCard) error {
	fmt.Printf("# - %v Date - %v Person - %v\n", cc.CCid, cc.Date.Format("DateOnly"), cc.Person)
	return nil
}
