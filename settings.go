package main

import (
	"context"
	"fmt"
	"strings"
)

// TODO: this will load the basic settings, checking if it was set up already
// as a place holder, just say "first time set up complete"
// set the values as you want so you don't have to repeatedly
// not a stored value, grab bools as needed from the db
func (cfg *Config) checkSettings() error {
	settings, err := cfg.db.GetSettings(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			fmt.Println("Creating settings file...")
			cfg.db.SetUpSettings(context.Background())
		} else {
			return err
		}
	}

	if settings.SettingsComplete == false {
		cfg.db.UpdateActivateSelf(context.Background(), true)
		cfg.db.FirstTimeSetupComplete(context.Background())
		fmt.Println("First time set up complete!")
	}

	return nil
}
