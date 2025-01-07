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

// change if you want a million things printed or not
const verbose bool = false

// CURRENTLY:
// the great refactoring
// using like 'get cc from flag' instead of copy and pasting code
// using check functions and prompts for what i actually want
// greatly neatening up clode

// do i start putting things in other packages? would literally be impossible lmao
// but we can reduce the net number of files by menu, all the states are the same so start there
// then start in the same order with the refactoring

// position by prompt
// position by flag

// investigator by prompt
// investigator by flag

// protocol by prompt
// protocol by flag

// strain by prompt and flag i think i have

// NEXT:

// AFTER THAT:
// the EASY CICD testing
// AFTER THAT:
// adding reminders, orders, and tests for those
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

func main() {
	godotenv.Load()
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

	for true {
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
