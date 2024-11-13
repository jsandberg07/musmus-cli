package main

func getMap() map[string]Command {
	exitCmd := Command{
		name:        "exit",
		description: "Exits program",
		function:    exitCommand,
	}

	printCmd := getPrintCmd()

	helpCmd := Command{
		name:        "help",
		description: "prints descriptions of all available functions",
		function:    helpCommand,
	}

	cmdSlice := []Command{exitCmd, printCmd, helpCmd}
	commandMap := make(map[string]Command)
	for _, cmd := range cmdSlice {
		commandMap[cmd.name] = cmd
	}

	return commandMap
}
