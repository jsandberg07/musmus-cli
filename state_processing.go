package main

func getProcessingMap() map[string]Command {
	activateCmd := getCCActivationCmd()
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
		cliMessage:      "cc processing",
	}

	return &processingState
}
