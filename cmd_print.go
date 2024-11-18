package main

import (
	"errors"
	"fmt"
	"strings"
)

func getPrintCmd() Command {

	printFlags := make(map[string]Flag)

	cFlag := Flag{
		symbol:      "c",
		description: "Allows custom text.",
		takesValue:  true,
	}
	printFlags["-"+cFlag.symbol] = cFlag

	bFlag := Flag{
		symbol:      "b",
		description: "Makes text uppercase.",
		takesValue:  false,
	}
	printFlags["-"+bFlag.symbol] = bFlag

	printCmd := Command{
		name:        "print",
		description: "Prints 'wow' or sometimes something other than 'wow.'",
		function:    printCommand,
		flags:       printFlags,
	}

	return printCmd
}

func printCommand(cfg *Config, args []Argument) error {
	fmt.Println("Printing...")
	output := "Wow!"
	uppercase := false

	for _, argument := range args {
		switch argument.flag {
		case "-b":
			uppercase = true
		case "-c":
			output = argument.value
		default:
			return errors.New("fake flag snuck into print")
		}
	}

	if uppercase {
		output = strings.ToUpper(output)
	}

	fmt.Println(output)

	return nil
}
