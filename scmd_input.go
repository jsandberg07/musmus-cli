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

	var arguments []string

	if len(splitArgs) != 0 {
		arguments = splitArgs[1:]
	}

	return arguments, nil

}

// maps is already reference type
func parseSubcommand(flags map[string]Flag, parameters []string) ([]Argument, error) {
	// used when in a subcommand, not expecting a command name. just give it the subcommand map.
	// flags -p, command like things do not so figure out how to do that
	// flags should also take a value for now so exploit that
	if len(parameters) == 0 {
		return nil, errors.New("Nothing entered. Please try again.")
	}

	var arguments []Argument

	for i := 0; i < len(parameters); i++ {
		if !strings.Contains(parameters[i], "-") {
			// - not included in flag
			err := fmt.Sprintf("%s isn't formatted as a flag, or a value without a flag", parameters[i])
			return nil, errors.New(err)
		}
		flag, ok := flags[parameters[i]]
		if !ok {
			// flag now allowed for this command
			err := fmt.Sprintf("%s is not a flag allowed for this command", parameters[i])
			return nil, errors.New(err)
		}

		tArg := Argument{}
		// check to see if the flag exists (not indexing out of bounds) && isn't also a flag
		// TODO: if the next param contains a - (like a name with a hyphen) it'll throw so make sure it just checks the first character
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
