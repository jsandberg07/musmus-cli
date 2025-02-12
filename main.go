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

// COME BACK TO IT:
// writing tests. spin up a fake db and do tests there. everything is too tight
// parsing works and that's like the one thing that isn't tied to a db
// https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
// https://circleci.com/blog/unit-testing-vs-integration-testing/
// use cmp, or reflect. I guess cmp is better. you can get the difference between tho values

// australia
// china #9
// china #11
// do it

// CURRENTLY:
// add permission checks to everything. You literally have the roles already you just need to check with logged in position stored in cfg
// this wil literally be an easy one

// After that:
// work on the read me cause it'll be different

// Next:
// add permissions that work (you can delete anybodys reminders currently)
// get the # of care days calculated for expenses
// adding people to protocols does nothing too
// can only add cage cards to allowed protocols check

// AFTER THAT:
// the great polishing
// logins (crpyto, storing that, creating new people, the admin tier account)
// AFTER THAT:
// the great readme-en-ing. write a readme and update the github page
// add some consts for the strings you use for consistency + fancy + easier changes
// like exiting program vs exiting command
// continue to use SIGNED ints but double check them!!! this is a deliberate design decision!
// AFTER THAT:
// the great adding the project to my portfolioening
// lots of test data! for fun!
// AFTER THAT:
// the great job applyening, apply for jobs
// DURING THAT:
// the adding automated DB testing as well! cause that's a thing!
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

	err = getTodaysReminders(&cfg)
	if err != nil {
		fmt.Println("Error getting today's reminders")
		fmt.Println(err)
	}

	err = getTodaysOrders(&cfg)
	if err != nil {
		fmt.Println("Error getting today's orders")
		fmt.Println(err)
	}

	// spacing :^3
	fmt.Println("")

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

		cmdName, err := readCommandName(text)
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
		/* removed because arguments are no longer passed to commands (they were just never used)
		// check to see if the flags are available, and if they take values, return flags and args
		arguments, err := parseCommandArguments(&command, parameters)
		if err != nil {
			fmt.Println(err)
			continue
		}
		*/

		// pass the args into the commands function, then run it
		err = command.function(&cfg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// spacing :^3
		fmt.Println()
	}

	// cool facts: this part of the code is never reached beacuse exit uses os dot Exit(0)
}
