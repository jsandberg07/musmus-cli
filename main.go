package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsandberg07/clitest/internal/database"

	_ "github.com/lib/pq"
)

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

	cfg := Config{
		currentState:         nil,
		nextState:            getMainState(),
		db:                   dbQueries,
		loggedInInvestigator: nil,
		loggedInPosition:     nil,
	}

	err = cfg.checkFirstTimeSetup()
	if err != nil {
		fmt.Printf("Error checking settings from DB: %s\n", err)
		os.Exit(1)
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

	fmt.Println("")

	// TODO: cfg.updateState() would pretty much just be this
	for {
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

		err = command.function(&cfg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println()
	}

	// exit uses os.Exit(0) so this part of the code isn't reached
	// perform clean up in the exit cmd
}
