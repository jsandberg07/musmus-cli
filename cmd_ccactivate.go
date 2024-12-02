package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

// TODO:
// new activation
// you go to the activate menu with params if youd like
// then in that menu you can also set the params like change person, date, # animals
// sum # of animals and then ding it later so you dont have a gorillion threads affecting that.
// im gonna fucking hate all the UUIDs i bet
// and i'll have plenty of anger with the other functions first
// this is literally the last step but i have an idea of how i want to do it

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
		if text == "pop" || text == "delete" {
			length := len(cardsToProcess)
			if length == 0 {
				fmt.Println("No cards have been entered.")
				continue
			}
			cardsToProcess = cardsToProcess[:length-1]
			continue
		}

		if err != nil {
			fmt.Println("Error reading input string")
			return err
		}
		ccid, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid input. Expecting a number or valid command.")
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
				processCard2(cfg, &cc)
			}()
		}
	}

	wg.Wait()

	return nil

}

// go routine for just printing cards, with a DB it'll be a sql thing
func processCard(cc *CageCard) error {
	fmt.Printf("# - %v Date - %v Person - %v\n", cc.CCid, cc.Date.Format("DateOnly"), cc.Person)
	return nil
}

// TODO: processCard3 :^3
func processCard2(cfg *Config, cc *CageCard) error {
	ccParams := database.ActivateCageCardParams{
		CcID:           int32(cc.CCid),
		ActivatedOn:    sql.NullTime{Time: cc.Date, Valid: true},
		InvestigatorID: uuid.New(),
	}
	createdCC, err := cfg.db.ActivateCageCard(context.Background(), ccParams)
	if err != nil {
		fmt.Println(createdCC)
	}

	return nil
}
