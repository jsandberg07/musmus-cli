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

// NEXT:
// work on the CLI
// more states for protocols, positions, investigators, ect
// each has an add, check, delete
// activating cards changes so that you change the date mid thing
// like -d to change the date, -a to change allotment
// start there really
// or just add the numbers
// checking permissions

// TODO:
// DRY up get investigator by name because you always have to check the length of the returned array
// return errors.new

// ALSO:
// set up CI testing
// you probably wont use it often but it's nice ^_^

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

	err = cfg.checkSettings()
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
		arguments, err := parseArguments(&command, parameters)
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

	resetCommand(&cfg, nil)
}
