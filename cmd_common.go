package main

import (
	"fmt"
	"os"
)

// common functions available in every state
// help, exit, back

func getCommonCmds() []Command {
	commonCmds := []Command{getBackCmd(), getStateHelpCmd(), getExitCmd()}
	return commonCmds
}

// just set the state as main
func getBackCmd() Command {
	backCmd := Command{
		name:        "back",
		description: "goes back to main menu",
		function:    backCommand,
	}

	return backCmd
}
func backCommand(cfg *Config) error {
	cfg.nextState = getMainState()
	return nil
}

func getExitCmd() Command {
	exitCmd := Command{
		name:        "exit",
		description: "Exits program.",
		function:    exitCommand,
	}

	return exitCmd
}

// because you os exit this way, you never hit the end of the program to clean things
// if you want to reset anything, put it here
func exitCommand(cfg *Config) error {
	fmt.Println("exiting...")
	os.Exit(0)
	return nil
}

// for the main menu
// prints all processes that are available for the state
func getStateHelpCmd() Command {
	helpCmd := Command{
		name:        "help",
		description: "Prints descriptions of all available functions.",
		function:    stateHelpCommand,
	}

	return helpCmd
}

func stateHelpCommand(cfg *Config) error {
	cmdMap := cfg.currentState.currentCommands
	for _, key := range cmdMap {
		fmt.Printf("* %s\n", key.name)
		fmt.Println(key.description)
		for _, key := range key.flags {
			fmt.Printf("%s - %s", key.symbol, key.description)
			if key.takesValue {
				fmt.Print(" Requires value.")
			}
			fmt.Println()
		}

	}
	return nil
}
