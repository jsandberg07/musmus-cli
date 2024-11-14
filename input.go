package main

import (
	"errors"
	"fmt"
	"strings"
)

// just clean up the inputs
// then we can have another function check and create the list of flags+args from the command
func readCommandName(input string) (string, []string, error) {
	if input == "" {
		return "", nil, errors.New("no input found")
	}

	splitArgs := strings.Split(input, " ")
	for i, arg := range splitArgs {
		splitArgs[i] = strings.TrimSpace(arg)
	}

	cmdName := splitArgs[0]
	var arguments []string

	if len(splitArgs) != 0 {
		arguments = splitArgs[1:]
	}

	return cmdName, arguments, nil

}

func parseArguments(cmd *Command, parameters []string) ([]Argument, error) {
	// no params passed in
	if len(parameters) == 0 {
		return nil, nil
	}

	var arguments []Argument

	for i := 0; i < len(parameters); i++ {
		if !strings.Contains(parameters[i], "-") {
			// - not included in flag
			err := fmt.Sprintf("%s isn't formatted as a flag, or a value without a flag", parameters[i])
			return nil, errors.New(err)
		}

		flag, ok := cmd.flags[parameters[i]]
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
