package main

import (
	"context"
	"fmt"
)

func (cfg *Config) login() error {
	// TODO: actually logging in. Just get the PI who can do everything.
	investigator, err := cfg.getInvestigator("Johnny Boi")
	if err != nil {
		fmt.Println("Error logging user in")
		return err
	}
	position, err := cfg.db.GetUserPosition(context.Background(), investigator.Position)
	if err != nil {
		fmt.Println("Error getting logged in user position")
		return err
	}

	cfg.loggedInInvestigator = &investigator
	cfg.loggedInPosition = &position

	return nil

}

func (cfg *Config) printLogin() {
	fmt.Printf("Logged in as %s -- %s\n", cfg.loggedInInvestigator.IName, cfg.loggedInPosition.Title)
}
