package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
)

// Commands that are available for each state / menu
func getCommonCmds() []Command {
	commonCmds := []Command{getBackCmd(), getStateHelpCmd(), getExitCmd()}
	return commonCmds
}

func getBackCmd() Command {
	backCmd := Command{
		name:        "back",
		description: "goes back to main menu",
		function:    backCommand,
		printOrder:  99,
	}

	return backCmd
}
func backCommand(cfg *Config) error {
	cfg.nextState = getMainState()
	return nil
}

func getExitCmd() Command {
	exitCmd := Command{
		name:        "exit",
		description: "Exits program.",
		function:    exitCommand,
		printOrder:  100,
	}

	return exitCmd
}

// uses os.Exit instead of going back to the loop in main. Any clean up upon exit needs to be done here.
func exitCommand(cfg *Config) error {
	fmt.Println("exiting...")
	os.Exit(0)
	return nil
}

// for the state / menu. prints all processes that are available for the state
func getStateHelpCmd() Command {
	helpCmd := Command{
		name:        "help",
		description: "Prints descriptions of all available functions.",
		function:    stateHelpCommand,
		printOrder:  100,
	}

	return helpCmd
}

func stateHelpCommand(cfg *Config) error {
	cmds := slices.Collect(maps.Values(cfg.currentState.currentCommands))
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].printOrder < cmds[j].printOrder
	})
	for _, cmd := range cmds {
		fmt.Printf("* %s\n", cmd.name)
		if cmd.description != "" {
			fmt.Print(cmd.description)
		}
		fmt.Println()
	}
	return nil
}
