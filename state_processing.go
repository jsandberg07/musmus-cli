package main

func getProcessingMap() map[string]Command {
	cmds := []Command{getCCActivationCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getProcessingState() *State {
	processingMap := getProcessingMap()
	processingState := State{
		currentCommands: processingMap,
		cliMessage:      "cc processing",
	}

	return &processingState
}
