package main

import (
	"bufio"
	"fmt"
	"os"
)

// NEXT:
// make a different state
// swap to it
// have some dumb cage card function that just prints shit (with go routines!)

func main() {
	fmt.Println("Hello borld")
	cfg := Config{
		currentState: nil,
		nextState:    getMainState(),
	}

	reader := bufio.NewReader(os.Stdin)

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
		err = command.function(arguments)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}
