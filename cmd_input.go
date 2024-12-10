package main

import (
	"errors"
	"fmt"
	"strings"
)

// NOT expecting a command name
// if string len = 1, just try to run that, otherwise send to be parsed
func readSubcommandInput(input string) ([]string, error) {
	if input == "" {
		fmt.Println("No input found")
		return []string{}, nil
	}

	splitArgs := strings.Split(input, " ")
	for i, arg := range splitArgs {
		splitArgs[i] = strings.TrimSpace(arg)
	}

	return splitArgs, nil

}

// maps is already reference type
// used in commands
// DONT HANDLE QUOTES
// JUST REPLACE _ WITH SPACES
// BINGOBANGO
func parseArguments(flags map[string]Flag, parameters []string) ([]Argument, error) {
	// used when in a subcommand, not expecting a command name. just give it the subcommand map.
	// flags -p, command like things do not so figure out how to do that
	// flags should also take a value for now so exploit that
	if len(parameters) == 0 {
		return nil, errors.New("Nothing entered. Please try again.")
	}

	var arguments []Argument

	for i := 0; i < len(parameters); i++ {
		// REWRITE
		// just do the two: flags with values, commands without
		// why is it like this? who the fuck know
		flag, ok := flags[parameters[i]]
		if !ok {
			err := fmt.Sprintf("%s is not a flag allowed for this command", parameters[i])
			return nil, errors.New(err)
		}

		tArg := Argument{}
		if flag.takesValue {
			tArg.flag = parameters[i]
			if i+1 == len(parameters) || strings.Contains(parameters[i+1], "-") {
				err := fmt.Sprintf("%s is a flag that takes a value", parameters[i])
				return nil, errors.New(err)
			}
			i++
			tArg.value = parameters[i]
		} else {
			tArg.flag = parameters[i]
		}
		arguments = append(arguments, tArg)

	}

	return arguments, nil
}
