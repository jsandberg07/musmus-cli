package main

func getInvestigatorsMap() map[string]Command {
	cmds := []Command{getAddInvestigatorCmd(), getEditInvestigatorCmd()}
	commandsMap := cmdMapHelper(cmds)

	return commandsMap
}

func getInvesitatorsState() *State {
	investigatorsMap := getInvestigatorsMap()
	investigatorState := State{
		currentCommands: investigatorsMap,
		cliMessage:      "investigator",
	}

	return &investigatorState
}
