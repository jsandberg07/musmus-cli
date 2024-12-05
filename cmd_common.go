package main

// common functions available in every state
// help, exit, back

func getCommonCmds() []Command {
	commonCmds := []Command{getBackCmd(), getStateHelpCmd(), getExitCmd()}
	return commonCmds
}
