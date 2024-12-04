package main

import (
	"errors"
)

// just set the state as main
func getBackCmd() Command {
	backCmd := Command{
		name:        "back",
		description: "goes back to main menu",
		function:    backCommand,
	}

	return backCmd
}

func backCommand(cfg *Config, args []Argument) error {
	if len(args) != 0 {
		return errors.New("Back takes no params. Just takes you back to the main menu.")
	}

	cfg.nextState = getMainState()

	return nil
}
