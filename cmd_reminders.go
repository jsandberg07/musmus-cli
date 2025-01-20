package main

import (
	"bufio"
	"fmt"
	"os"
)

func getAddReminderCmd() Command {
	addReminderFlags := make(map[string]Flag)
	addReminderCmd := Command{
		name:        "add",
		description: "Used for adding reminders",
		function:    addReminderFunction,
		flags:       addReminderFlags,
	}

	return addReminderCmd
}

func getAddReminderFlags() map[string]Flag {
	addReminderFlags := make(map[string]Flag)
	XFlag := Flag{
		symbol:      "X",
		description: "Sets X",
		takesValue:  false,
	}
	addReminderFlags["-"+XFlag.symbol] = XFlag

	// ect as needed or remove the "-"+ for longer ones

	return addReminderFlags

}

// look into removing the args thing, might have to stay
func addReminderFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getAddReminderFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
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

		// do weird behavior here

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-X":
				exit = true
			default:
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}
