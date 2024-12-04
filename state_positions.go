package main

func getInvestigatorMap() map[string]Command {
	// put that commands related to investigators you want here
	cmds := []Command{}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getInvestigatorState() *State {
	investigatorMap := getInvestigatorMap()
	processingState := State{
		currentCommands: investigatorMap,
		cliMessage:      "investigator",
	}

	return &processingState
}
