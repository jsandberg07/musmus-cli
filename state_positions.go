package main

func getPositionMap() map[string]Command {
	// put that commands related to positions you want here
	cmds := []Command{getAddPositionCmd(), getEditPositionCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getPositionState() *State {
	positionsMap := getPositionMap()
	positionState := State{
		currentCommands: positionsMap,
		cliMessage:      "position",
	}

	return &positionState
}
