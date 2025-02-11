package main

import (
	"errors"
	"fmt"
	"strings"
)

// just clean up the inputs
// then we can have another function check and create the list of flags+args from the command
func readCommandName(input string) (string, error) {
	if input == "" {
		return "", errors.New("no input found")
	}

	splitArgs := strings.Split(input, " ")
	for i, arg := range splitArgs {
		splitArgs[i] = strings.TrimSpace(arg)
	}

	if len(splitArgs) > 1 {
		fmt.Println("Too many command names entered. Only 1 is needed.")
	}

	cmdName := splitArgs[0]

	/* formerly returned an array of strings that were never used
	var arguments []string

	if len(splitArgs) != 0 {
		arguments = splitArgs[1:]
	}
	return cmdName, arguments, nil
	*/

	return cmdName, nil
}

/* removed because excessively complicated way to say "goto state." Removed with simplification of commands (but i still think it's clever)
// for now used outside of commands but is a lot for what is essentially "goto this state"
func parseCommandArguments(cmd *Command, parameters []string) ([]Argument, error) {
	// no params passed in
	if len(parameters) == 0 {
		return nil, nil
	}

	var arguments []Argument

	for i := 0; i < len(parameters); i++ {
		if !strings.Contains(parameters[i], "-") {
			// - not included in flag
			err := fmt.Sprintf("%s isn't formatted as a flag, or a value without a flag", parameters[i])
			return nil, errors.New(err)
		}

		flag, ok := cmd.flags[parameters[i]]
		if !ok {
			// flag now allowed for this command
			err := fmt.Sprintf("%s is not a flag allowed for this command", parameters[i])
			return nil, errors.New(err)
		}

		tArg := Argument{}
		// check to see if the flag exists (not indexing out of bounds) && isn't also a flag
		// TODO: if the next param contains a - (like a name with a hyphen) it'll throw so make sure it just checks the first character
		if flag.takesValue {
			tArg.flag = parameters[i]
			if i+1 == len(parameters) || strings.Contains(parameters[i+1], "-") {
				err := fmt.Sprintf("%s is a flag that takes a value", parameters[i])
				return nil, errors.New(err)
			}
			i++
			tArg.value = parameters[i]
		} else {
			tArg.flag = parameters[i]
		}
		arguments = append(arguments, tArg)
	}

	return arguments, nil
}
*/

// takes commands, adds common commands, then makes a map using its name as the key.
// used in every state.
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

func getMainMap() map[string]Command {
	cmds := []Command{getGotoCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getMainState() *State {
	mainMap := getMainMap()

	mainState := State{
		currentCommands: mainMap,
		cliMessage:      "main",
	}

	return &mainState
}

// TODO: gotta be some way to condense this into one function. Theyre literally copy and paste
// input a string, use a switch to get a map, set the cli message, return &state
func getInvestigatorsMap() map[string]Command {
	cmds := []Command{getAddInvestigatorCmd(), getEditInvestigatorCmd()}
	commandsMap := cmdMapHelper(cmds)

	return commandsMap
}
func getInvesitatorsState() *State {
	investigatorsMap := getInvestigatorsMap()
	investigatorState := State{
		currentCommands: investigatorsMap,
		cliMessage:      "investigator",
	}

	return &investigatorState
}

func getPositionMap() map[string]Command {
	// put that commands related to positions you want here
	cmds := []Command{getAddPositionCmd(), getEditPositionCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getPositionState() *State {
	positionsMap := getPositionMap()
	positionState := State{
		currentCommands: positionsMap,
		cliMessage:      "position",
	}

	return &positionState
}

func getProcessingMap() map[string]Command {
	cmds := []Command{getCCActivationCmd(),
		getAddCCCmd(),
		getCCDeactivationCmd(),
		getCCReactivateCmd(),
		getCCInactivateCmd()}
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

func getQueriesMap() map[string]Command {
	// put your query Commands here
	cmds := []Command{getCCQueriesCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getQueriesState() *State {
	queriesMap := getQueriesMap()
	queriesState := State{
		currentCommands: queriesMap,
		cliMessage:      "queries",
	}

	return &queriesState
}

func getSettingsMap() map[string]Command {
	// put your settings Commands here
	cmds := []Command{getChangeSettingsCmd()}
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

func getStrainsMap() map[string]Command {
	cmds := []Command{getAddStrainCmd(), getEditStrainCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getStrainsState() *State {
	strainsMap := getStrainsMap()

	strainState := State{
		currentCommands: strainsMap,
		cliMessage:      "strains",
	}

	return &strainState
}

func getRemindersMap() map[string]Command {
	cmds := []Command{getAddReminderCmd(), getDeleteReminderCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getRemindersState() *State {
	remindersMap := getRemindersMap()

	remindersState := State{
		currentCommands: remindersMap,
		cliMessage:      "reminders",
	}

	return &remindersState
}

func getOrdersMap() map[string]Command {
	cmds := []Command{getAddOrderCmd(), getEditOrderCmd(), getReceiveOrderCmd()}
	commandMap := cmdMapHelper(cmds)

	return commandMap
}
func getOrdersState() *State {
	ordersMap := getOrdersMap()

	ordersState := State{
		currentCommands: ordersMap,
		cliMessage:      "orders",
	}

	return &ordersState
}
