package main

import "fmt"

func getHelpCmd() Command {
	helpCmd := Command{
		name:        "help",
		description: "prints descriptions of all available functions",
		function:    helpCommand,
	}

	return helpCmd
}

func helpCommand(cfg *Config, args []Argument) error {
	cmdMap := cfg.currentState.currentCommands
	for _, key := range cmdMap {
		fmt.Println(key.name)
		fmt.Println(key.description)
		fmt.Println(key.flags)
	}
	return nil
}
