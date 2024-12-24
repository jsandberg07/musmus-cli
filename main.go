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
// queries! now that you have data, let's actually be able to see it
// after the success of like adding all these functions (kinda)
// i know where the top of the mountain is. I'm above the tree line.
// there's still much to do and that's scary but it won't take as long
// 5k lines of code added want to work more consistently
// FOR NOW literally just write to a file
// mostly about getting cage cards
// do active by date, active for investigator, active by protocol
// and a custom query
// to get a particular param, deactivated on a particular day, activated on a particular day
// how do to that? i dont want 16 functions with each option on or not
// gotta remember joins

// NEXT:
// the great refactoring
// using like 'get cc from flag' instead of copy and pasting code
// using check functions and prompts for what i actually want
// greatly neatening up clode
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
