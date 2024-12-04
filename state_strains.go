package main

func getStrainsMap() map[string]Command {
	cmds := []Command{getSetStateCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getStrainsState() *State {
	strainsMap := getStrainsMap()

	mainState := State{
		currentCommands: strainsMap,
		cliMessage:      "strains",
	}

	return &mainState
}
