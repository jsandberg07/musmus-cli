package main

// common functions available in every state
// help, exit, set state

func getCommonCmds() []Command {
	commonCmds := []Command{getBackCmd(), getStateHelpCmd(), getExitCmd()}
	return commonCmds
}
