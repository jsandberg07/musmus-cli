package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jsandberg07/clitest/internal/database"
)

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
