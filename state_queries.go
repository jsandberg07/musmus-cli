package main

func getQueriesMap() map[string]Command {
	// put your query Commands here
	cmds := []Command{}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getQueriesState() *State {
	queriesMap := getQueriesMap()
	queriesState := State{
		currentCommands: queriesMap,
		cliMessage:      "queries",
	}

	return &queriesState
}
