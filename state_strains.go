package main

func getStrainsMap() map[string]Command {
	cmds := []Command{getAddStrainCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getStrainsState() *State {
	strainsMap := getStrainsMap()

	strainState := State{
		currentCommands: strainsMap,
		cliMessage:      "strains",
	}

	return &strainState
}
