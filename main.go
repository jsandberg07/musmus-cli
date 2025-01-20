package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsandberg07/clitest/internal/database"

	_ "github.com/lib/pq"
)

// CURRENTLY:
// remove args param in command functions lemoo
// writing tests for functions
// read these:
// https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
// https://circleci.com/blog/unit-testing-vs-integration-testing/
// use cmp, or reflect. I guess cmp is better. you can get the difference between tho values
// go has built in conversion for identical structs! use it instead of normalizing!
// x := t1; y := t2(x); works!
// see what things you can actually break down and parse there might not be a lot currently (bad)
// will end up being reliant on spinning up a test db for the bulk of the testing feels like

// dont forget to run your tests before merging!
// Next:
// Cc activation is weak, make it a go routine that just does it as it goes. it's fast enough im sure and not http based. you're literally typing by hand. it's fake remember?

// After NEXT:
// adding reminders, orders, and tests for those

// reminders have a CC#, can be set to automatically add to CC activation, or an order. do E something or other, or a reminder for ~21 days from now, or whatever
// reminders show up for a person or for everybody
// see reminders with dates for today, next week, export. no past stuff. once it's done, have it be done (dont delete it anyway)

// AFTER THAT:
// the great polishing
// making maps print sortedly
// adding permissions that work! you already have rolls. and logins
// AFTER THAT:
// the great readme-en-ing. write a readme and update the github page
// AFTER THAT:
// the great adding the project to my portfolioening
// lots of test data! for fun!
// AFTER THAT:
// the great job applyening, apply for jobs
// DURING THAT:
// the adding automated DB testing as well! cause that's a thing!
// SOMEWHERE DURING THAT
// i changed my mind and was cage cards to have go routines instead, so closer to cayuse process as you go
// it's fast cause it's not html based
// ADDITIONALLY:
// a RESTful server version! throw out all your cli code and keep the sql! no more parsing! just json!

// change if you want a million things printed or not
const verbose bool = false

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load env file: %s\n", err)
		os.Exit(1)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Could not open DB: %s\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	fmt.Println("Hello borld")
	cfg := Config{
		currentState:         nil,
		nextState:            getMainState(),
		db:                   dbQueries,
		loggedInInvestigator: nil,
		loggedInPosition:     nil,
	}

	err = cfg.db.ResetDatabase(context.Background())
	if err != nil {
		fmt.Printf("Error resetting DB: %s", err)
		os.Exit(1)
	}

	err = cfg.loadSettings()
	if err != nil {
		fmt.Printf("Error checking settings from DB: %s", err)
		os.Exit(1)
	}

	err = cfg.testData()
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)

	err = cfg.login()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfg.printLogin()
	fmt.Println("\n* Welcome to Musmus!")

	for {
		// check if new state
		if cfg.nextState != nil {
			cfg.currentState = cfg.nextState
			cfg.nextState = nil
		}

		fmt.Printf(">%s - ", cfg.currentState.cliMessage)

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		cmdName, parameters, err := readCommandName(text)
		if err != nil {
			fmt.Println("oops error getting command name")
			fmt.Println(err)
			os.Exit(1)
		}

		command, ok := cfg.currentState.currentCommands[cmdName]
		if !ok {
			fmt.Println("Invalid command")
			continue
		}
		// check to see if the flags are available, and if they take values, return flags and args
		arguments, err := parseCommandArguments(&command, parameters)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// pass the args into the commands function, then run it
		err = command.function(&cfg, arguments)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// spacing :^3
		fmt.Println()
	}

	// cool facts: this part of the code is never reached beacuse exit uses os dot Exit(0)
}
