package main

func getMainMap() map[string]Command {
	cmds := []Command{getSetStateCmd()}
	commandMap := cmdMapHelper(cmds)

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
