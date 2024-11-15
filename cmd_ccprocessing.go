package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		description: "sets date. Defaults to today",
		takesValue:  true,
	}
	activateFlags["-"+dFlag.symbol] = dFlag

	pFlag := Flag{
		symbol:      "p",
		description: "sets person. Defaults to 'mouse'",
		takesValue:  true,
	}
	activateFlags["-"+pFlag.symbol] = pFlag

	activateCmd := Command{
		name:        "activate",
		description: "used for activating cards",
		function:    activateCommand,
		flags:       activateFlags,
	}

	return activateCmd
}

// process itself
func activateCommand(cfg *Config, args []Argument) error {
	date := time.Now()
	person := "Mouse"
	var err error // because blocks

	for _, argument := range args {
		switch argument.flag {
		case "-d":
			date, err = parseDate(argument.value)
			if err != nil {
				return errors.New("Couldnt parse date. Try this format: day/month/year")
			}
		case "-p":
			// when db, add to make sure somebody is in the db
			person = argument.value
		}
	}

	fmt.Printf("Enter cage cards to activate. Date: %v - Person: %s", date.Format("DateOnly"), person)
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

	if len(cardsToProcess) != 0 {
		fmt.Println("Here's where cards would be processed :^3")
	}

	return nil

}

// because who knows what dates people are going to enter
func parseDate(string) (time.Time, error) {
	// TODO: make this literally not just return yesterday
	return time.Now().AddDate(0, 0, -1), nil
}

// go routine for just printing cards, with a DB it'll be a sql thing
func processCard(cc *CageCard) error {
	fmt.Printf("# - %v Date - %v Person - %v", cc.CCid, cc.Date.Format("DateOnly"), cc.Person)
	return nil
}
