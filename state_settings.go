package main

func getSettingsMap() map[string]Command {
	// put your settings Commands here
	cmds := []Command{}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}

func getSettingsState() *State {
	settingsMap := getSettingsMap()

	settingsState := State{
		currentCommands: settingsMap,
		cliMessage:      "settings",
	}

	return &settingsState
}
