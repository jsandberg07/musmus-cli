package main

// common functions available in every state
// help, exit, set state

func getCommonCmds() []Command {
	commonCmds := []Command{getSetStateCmd(), getHelpCmd(), getExitCmd()}
	return commonCmds
}
