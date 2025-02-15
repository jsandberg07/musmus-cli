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

// Next:
// some polishing while i have the time, then i'll go through set up on my laptop and write the read me with that
// cover up passwords - done
// reorganize files lmao
// returning ConstReturnString as an error message and using that to return
// and make that different from exiting. or closing or whatever. get a theseaurus.
// More test data! doesn't matter nobody uses it!
// more like "x entered, y updated" in edit functions. Can i do like a func printSet(type string, s {}interface) what would that do / look like. Yes it works, but watch setting structs
// check how many lines of code you've added lmao
// check for other TODOs and put there here:
//
// delete the massive ammount of comments that are everywhere

// After that:
// work on the read me cause it'll be different

// AFTER THAT:
// the great readme-en-ing. write a readme and update the github page
// add some consts for the strings you use for consistency + fancy + easier changes
// like exiting program vs exiting command
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

	// removed because i always start my programs with this but it's time to move on ToT
	// fmt.Println("Hello borld")

	// so we need to
	// 1. load settings into cfg
	// 2. check to see if first time set up has been run
	// 3. if not, prompt to see if they want to load test data or regular data
	// 4. if test data, load test data and mark test data in settings
	// 5. if test data is loaded, on each boot ask if they want to reset it
	// 6. reset reset's everything including settings so it's first time all over again
	// 7. have a reset option only the admin account can perform

	cfg := Config{
		currentState:         nil,
		nextState:            getMainState(),
		db:                   dbQueries,
		loggedInInvestigator: nil,
		loggedInPosition:     nil,
	}

	/*
		err = cfg.db.ResetDatabase(context.Background())
		if err != nil {
			fmt.Printf("Error resetting DB: %s", err)
			os.Exit(1)
		}
	*/

	err = cfg.checkFirstTimeSetup()
	if err != nil {
		fmt.Printf("Error checking settings from DB: %s\n", err)
		os.Exit(1)
	}

	/* removed becuse it's called in first time set up
	err = cfg.testData()
	if err != nil {
		fmt.Println(err)
	}
	*/

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
