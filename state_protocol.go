package main

func getProtocolMap() map[string]Command {
	// put your protocol Commands here
	cmds := []Command{}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getProtocolState() *State {
	protocolMap := getProtocolMap()

	protocolState := State{
		currentCommands: protocolMap,
		cliMessage:      "protocol",
	}

	return &protocolState
}
