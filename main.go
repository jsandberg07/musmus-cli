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

// wasted an hour to realize i can't remove []Arguments from functions. so maybe not a waste. It's tied to goto -cc ect.
// "Problem" is reuse of types, but it's a small problem. Goto is probably gonna be the most used function (tied with back and exit)
// addendum: remove the goto function and replace it with regular functions. it's the only one that uses the []Arguments ie why you cant remove it

// australia
// china #9
// china #11
// do it

// CURRENTLY:
// reminders, orders
// run tests before merging!
// reminders have a CC#, can be set to automatically add to CC activation, or an order. do E something or other, or a reminder for ~21 days from now, or whatever
// reminders show up for a person or for everybody
// see reminders with dates for today, next week, export. no past stuff. once it's done, have it be done (dont delete it anyway)

// fix CC activation
// add the ability to automatically create reminders when entering breeding for next day or E something for E days + 1 or just reminder days then note something like that

// Next:
// Cc activation is weak, make it a go routine that just does it as it goes. it's fast enough im sure and not http based. you're literally typing by hand. it's fake remember?

// AFTER THAT:
// the great polishing
// making maps print sortedly
// adding permissions that work! you already have rolls. and logins
// DRY up the state function, you don't need a separate function for each one. use a string and switch to return a state
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
