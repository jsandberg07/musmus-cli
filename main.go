package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func exitCommand(args []Argument) error {
	fmt.Println("exiting...")
	os.Exit(0)
	return nil
}

// can assume all fake flags have been taken care of, all arguments are able to be operated on
// operate on general flags and set bools
// operate on value flags and set values
func printCommand(args []Argument) error {
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

func helpCommand(args []Argument) error {
	cmdMap := getMap()
	for _, key := range cmdMap {
		fmt.Println(key.name)
		fmt.Println(key.description)
		fmt.Println(key.flags)
	}
	return nil
}

// each function takes a []string that is passed in earlier
// check with acceptable args
// returns an error

// need a classic input getter
// that you type shit into
// and it cleans the inputs
// creates the args and flags
// and returns it (or errors out)

// need a new reader
// use bufio newReader getstring or something

func main() {
	// start program
	fmt.Println("Hello borld")

	// command map
	commandMap := getMap()

	// enter a command
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a command and flags")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("oops input error")
		fmt.Println(err)
		os.Exit(1)
	}

	// parse input
	// get the command from the command name
	cmdName, parameters, err := readCommandName(text)
	if err != nil {
		fmt.Println("oops error getting command name")
		fmt.Println(err)
		os.Exit(1)
	}

	command, ok := commandMap[cmdName]
	if !ok {
		fmt.Println("Invalid command")
		os.Exit(1)
	}
	// check to see if the flags are available, and if they take values, return flags and args
	arguments, err := parseArguments(&command, parameters)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// pass the args into the commands function, then run it
	err = command.function(arguments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
