package main

// add a way to swap between states
// linear states a->b->c or web (go from a to c directly)
// map just for states?
// function that gets COMMON maps that are used often
func getProcessingMap() map[string]Command {
	activateCmd := getActivateCmd()
	commonCmds := getCommonCmds()
	cmdSlice := []Command{activateCmd}
	cmdSlice = append(cmdSlice, commonCmds...)
	commandMap := make(map[string]Command)
	for _, cmd := range cmdSlice {
		commandMap[cmd.name] = cmd
	}

	return commandMap
}

func getProcessingState() *State {
	processingMap := getProcessingMap()
	processingState := State{
		currentCommands: processingMap,
		cliMessage:      "processing",
	}

	return &processingState
}
