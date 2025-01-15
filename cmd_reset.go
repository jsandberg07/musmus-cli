package main

/* possibly implemented later, seems excessive though compared to just truncating via sql
// really truly commented out to satisfy staticcheck
func getResetCmd() Command {
	// maybe add flags later to reset only 1 tables
	resetCmd := Command{
		name:        "reset",
		description: "Resets all tables (for testing purposes).",
		function:    resetCommand,
	}

	return resetCmd
}

func resetCommand(cfg *Config, args []Argument) error {
	err := cfg.db.ResetCageCards(context.Background())
	if err != nil {
		fmt.Println("Error resetting cage cards")
		return err
	}
	fmt.Println("Cage cards deleted from DB.")
	return nil
}
*/
