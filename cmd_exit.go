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

func exitCommand(cfg *Config, args []Argument) error {
	fmt.Println("exiting...")
	os.Exit(0)
	return nil
}
