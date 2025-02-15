package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

// Parses int when passed in via flag (as opposed to prompt)
func getNumberFromFlag(input string) (int, error) {
	if input == "" {
		fmt.Println("No input found. Please try again.")
		return 0, nil
	}
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Could not read number from input")
		return 0, err
	}
	return num, nil
}

// TODO: decide on which version is better. Note: sql doesnt throw an error when a []struct query returns nothing
// but will for single values
func getInvestigatorByFlag2(cfg *Config, input string) (database.Investigator, error) {
	if input == "" {
		fmt.Println("No input found, please try again")
		return database.Investigator{}, nil
	}
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
	if err != nil {
		fmt.Println("Error getting investigator from db")
		return database.Investigator{}, err
	}
	if len(investigators) == 0 {
		fmt.Println("No investigator by that name found. Nicknames also work as well.")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Vague investigator name. Please try again.")
		return database.Investigator{}, nil
	}
	return investigators[0], nil
}

func getInvestigatorByFlag(cfg *Config, i string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), i)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting investigators from DB")
		return database.Investigator{}, err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		fmt.Println("Investigator not found. Please try again")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Vague investigator name. Please try again")
		return database.Investigator{}, nil
	}
	if len(investigators) == 0 {
		fmt.Println("Investigator not found. Please try again")
		return database.Investigator{}, nil
	}
	return investigators[0], nil
}

func getPositionByFlag(cfg *Config, title string) (database.Position, error) {
	position, err := cfg.db.GetPositionByTitle(context.Background(), title)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("Error getting position from DB")
		return database.Position{}, err
	}
	if err.Error() == "sql: no rows in result set" {
		fmt.Println("No position by that title found")
		return database.Position{}, err
	}
	return position, nil
}

func getOrderByFlag(cfg *Config, input string) (database.Order, error) {
	order, err := cfg.db.GetOrderByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// no order found
		fmt.Println("No order by that number found. Please try again.")
		return database.Order{}, nil
	}
	if err != nil {
		// any other error
		fmt.Println("Error getting order from DB.")
		return database.Order{}, err
	}

	// found and ok
	return order, nil

}

func getProtocolByFlag(cfg *Config, n string) (database.Protocol, error) {
	protocol, err := cfg.db.GetProtocolByNumber(context.Background(), n)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting protocol from DB")
		return database.Protocol{}, err

	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		// no results
		fmt.Println("Protocol by that number not found. Please try again")
		return database.Protocol{}, nil
	}

	return protocol, nil
}

func getStrainByFlag(cfg *Config, input string) (database.Strain, error) {
	strain, err := cfg.db.GetStrainByName(context.Background(), input)

	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error with DB
		fmt.Println("Error getting strain from DB")
		return database.Strain{ID: uuid.Nil}, err
	}

	if err == nil {
		// strain found by name
		return strain, nil
	}

	strain, err = cfg.db.GetStrainByCode(context.Background(), input)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error with DB
		fmt.Println("Error getting strain from DB")
		return database.Strain{ID: uuid.Nil}, err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		fmt.Println("Strain not found by name or number. Please try again.")
		return database.Strain{ID: uuid.Nil}, nil
	}

	// strain found by code
	return strain, nil
}
