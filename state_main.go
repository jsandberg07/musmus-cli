package main

func getMainMap() map[string]Command {
	setStateCmd := getSetStateCmd()
	commonCmds := getCommonCmds()
	cmdSlice := []Command{setStateCmd}
	cmdSlice = append(cmdSlice, commonCmds...)
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
