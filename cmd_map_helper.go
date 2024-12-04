package main

func cmdMapHelper(cmds []Command) map[string]Command {
	commonCmds := getCommonCmds()
	cmdSlice := []Command{}
	cmdSlice = append(cmdSlice, cmds...)
	cmdSlice = append(cmdSlice, commonCmds...)
	commandMap := make(map[string]Command)
	for _, cmd := range cmdSlice {
		commandMap[cmd.name] = cmd
	}

	return commandMap

}
