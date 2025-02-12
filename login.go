package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jsandberg07/clitest/internal/database"
)

func (cfg *Config) login() error {
	// TODO: actually logging in. Just get the PI who can do everything.
	// dan for testing permissions NOT being allowed
	// investigator, err := cfg.getInvestigator("Dan")
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

func (cfg *Config) getInvestigator(name string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), name)
	if err != nil {
		fmt.Println("Error getting investigator")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Error getting investigator")
		return database.Investigator{}, errors.New("vague investigator name")
	}

	return investigators[0], nil
}
