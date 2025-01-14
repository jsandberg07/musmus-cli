package main

import (
	"context"
	"errors"
	"fmt"
)

// TODO: this will load the basic settings, checking if it was set up already
// as a place holder, just say "first time set up complete"
// set the values as you want so you don't have to repeatedly
// not a stored value, grab bools as needed from the db
// ALSO TODO: there's a first time setup complete row that doesnt do anything yet
// and i'll add functionality to that later (asking what settings to use, or load test data)
// but for now this is FINE
func (cfg *Config) loadSettings() error {
	settings, err := cfg.db.GetSettings(context.Background())
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}

	if len(settings) == 1 {
		return nil
	}

	if len(settings) == 0 {
		fmt.Println("Creating settings file...")
		err := cfg.db.SetUpSettings(context.Background())
		if err != nil {
			return err
		}
		return nil
	}

	if len(settings) > 1 {
		return errors.New("too many rows of settings found. Should only ever be 1")
	}

	return errors.New("shouldn't see this error while loading settings")
}
