package main

import (
	"fmt"
	"os"
)

func getExitCmd() Command {
	exitCmd := Command{
		name:        "exit",
		description: "Exits program.",
		function:    exitCommand,
	}

	return exitCmd
}

// because you os exit this way, you never hit the end of the program to clean things
// if you want to reset anything, put it here
func exitCommand(cfg *Config, args []Argument) error {
	fmt.Println("exiting...")
	os.Exit(0)
	return nil
}
