package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// TODO: this will load the basic settings, checking if it was set up already
// as a place holder, just say "first time set up complete"
// set the values as you want so you don't have to repeatedly
// not a stored value, grab bools as needed from the db
// ALSO TODO: there's a first time setup complete row that doesnt do anything yet
// and i'll add functionality to that later (asking what settings to use, or load test data)
// but for now this is FINE
func (cfg *Config) checkFirstTimeSetup() error {
	setting, err := cfg.db.GetSettings(context.Background())
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = cfg.runFirstTimeSetUp()
		if err != nil {
			return err
		}
		// first time set up completed
		return nil
	}

	if !setting.SettingsComplete {
		fmt.Println("First time setup was not completed successfully. Resetting DB, please try again")
		err = cfg.db.ResetDatabase(context.Background())
		if err != nil {
			return err
		}
		os.Exit(0)
	}

	// settings loaded, test data present
	if setting.TestDataLoaded {
		// ask if they want to reset for regular use
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Test data has been previously loaded. Enter 'reset' if you'd like to reset the DB, or anything else to continue")
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "reset" {
			fmt.Println("Resetting DB...")
			err := cfg.db.ResetDatabase(context.Background())
			if err != nil {
				return err
			}
			os.Exit(0)
		} else {

		}
	}

	// settings loaded, not test data
	return nil

}

func (cfg *Config) runFirstTimeSetUp() error {
	fmt.Println("First time setting up database")

	reader := bufio.NewReader(os.Stdin)
	err := cfg.db.SetUpSettings(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Enter 'test' if you'd like to load test data, or anything else to continue regular set up")
	fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input string: %s", err)
		os.Exit(1)
	}
	input := strings.TrimSpace(text)
	if input == "test" {
		err = cfg.testData()
		if err != nil {
			return err
		}
	} else {
		// just create an admin so you can log in and set things up yourself
		err := cfg.createAdmin()
		if err != nil {
			return err
		}
	}

	err = cfg.db.FirstTimeSetupComplete(context.Background())
	if err != nil {
		return err
	}

	return nil

}
