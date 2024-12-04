package main

func getInvestigatorsMap() map[string]Command {
	// addInvestigatorCmd := getAddInvestigatorCmd()
	commonCmds := getCommonCmds()
	cmdSlice := commonCmds
	commandMap := make(map[string]Command)
	for _, cmd := range cmdSlice {
		commandMap[cmd.name] = cmd
	}

	return commandMap
}

func getInvesitatorsState() *State {
	investigatorsMap := getInvestigatorsMap()
	investigatorState := State{
		currentCommands: investigatorsMap,
		cliMessage:      "investigator",
	}

	return &investigatorState
}
