package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/jsandberg07/clitest/internal/database"
)

func getChangeSettingsCmd() Command {
	settingsFlags := make(map[string]Flag)
	settingsCmd := Command{
		name:        "change",
		description: "Used for reviewing and changing settings",
		function:    changeSettingsFunction,
		flags:       settingsFlags,
		printOrder:  1,
	}

	return settingsCmd
}

func getChangeSettingsFlags() map[string]Flag {
	settingsFlags := make(map[string]Flag)
	aFlag := Flag{
		symbol:      "-a",
		description: "Toggles the 'only activate self' setting.\nTrue means investigators can activate cards that aren't in their own name.",
		takesValue:  false,
		printOrder:  1,
	}
	settingsFlags[aFlag.symbol] = aFlag

	rFlag := Flag{
		symbol:      "-r",
		description: "Review settings.\nDisplays the current settings BEFORE any changes are made.",
		takesValue:  false,
		printOrder:  2,
	}
	settingsFlags[rFlag.symbol] = rFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the current settings, prints them, then exits.",
		takesValue:  false,
		printOrder:  99,
	}
	settingsFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving.",
		takesValue:  false,
		printOrder:  100,
	}
	settingsFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available commands.",
		takesValue:  false,
		printOrder:  100,
	}
	settingsFlags[helpFlag.symbol] = helpFlag

	return settingsFlags
}

func changeSettingsFunction(cfg *Config) error {
	flags := getChangeSettingsFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)

	currentSetting, err := cfg.db.GetSettings(context.Background())
	if err != nil {
		fmt.Println("Error getting current settings")
		return err
	}

	updatedSettings := currentSetting

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-a":
				updatedSettings.OnlyActivateSelf = !updatedSettings.OnlyActivateSelf
				fmt.Printf("Only activate self set to %v\n", updatedSettings.OnlyActivateSelf)

			case "-r":
				printSettings(&currentSetting)

			case "save":
				if updatedSettings == currentSetting {
					fmt.Println("No changes were made. Exiting...")
				} else {
					fmt.Println("Saving...")
					// just an single bool for now so pass that in. Will become struct when more settings are added.
					usParams := updatedSettings.OnlyActivateSelf
					err := cfg.db.UpdateSettings(context.Background(), usParams)
					if err != nil {
						fmt.Println("Error saving settings.")
						return err
					}
				}
				exit = true

			case "exit":
				fmt.Println("Exiting without saving...")
				exit = true

			case "help":
				cmdHelp(flags)

			default:
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

// Only setting for now! Ignore first time set up
func printSettings(s *database.Setting) {
	fmt.Printf("* Only activate self: %v", s.OnlyActivateSelf)
}
