package main

func getMainMap() map[string]Command {
	exitCmd := getExitCmd()
	printCmd := getPrintCmd()
	helpCmd := getHelpCmd()

	cmdSlice := []Command{exitCmd, printCmd, helpCmd}
	commandMap := make(map[string]Command)
	for _, cmd := range cmdSlice {
		commandMap[cmd.name] = cmd
	}

	return commandMap
}

func getMainState() *State {
	mainMap := getMainMap()

	mainState := State{
		currentCommands: mainMap,
		cliMessage:      "main",
	}

	return &mainState
}
