package main

func getProtocolMap() map[string]Command {
	// put your protocol Commands here
	cmds := []Command{getAddProtocolCmd(), getEditProtocolCmd(), getAddInvestigatorToProtocolCmd()}
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
