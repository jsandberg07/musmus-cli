package main

import "fmt"

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

func stateHelpCommand(cfg *Config, args []Argument) error {
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
