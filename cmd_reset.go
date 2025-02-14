package main

import (
	"context"
	"fmt"
	"os"
)

func getResetCmd() Command {
	resetCmd := Command{
		name:        "reset",
		description: "Resets database",
		function:    resetCommand,
		printOrder:  2,
	}

	return resetCmd
}

func resetCommand(cfg *Config) error {
	if !cfg.loggedInPosition.IsAdmin {
		fmt.Println("Reset can only be performed by admin position")
	}

	err := cfg.db.ResetDatabase(context.Background())
	if err != nil {
		fmt.Println("Error resetting resetting DB")
		return err
	}
	fmt.Println("Database has been reset. Exiting...")
	os.Exit(0)
	return nil
}
