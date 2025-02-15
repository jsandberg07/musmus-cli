package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// temp file name, part of the refactor

// idea for how i should have created more reusable functions for all the other data types
// more generic get string with a prompt, instead of separate functions for everything
// just pass in "get investigator name to edit" instead of new function just to say "to edit"
// write the same program 3 times and you'll realize what you want you want to refactor

// Prints prompt (+ instructions on how to cancel), takes an input from the user, then runs it through the check function for uniqueness or
// if valid entry in the database. Will repeat if checkFunc returns an error. Returns an Error() const that can be checked
// for when the user decides to cancel and should be checked
func getStringPrompt(cfg *Config, prompt string, checkFunc func(*Config, string) error) (string, error) {
	fmt.Println(prompt + " or exit to cancel")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			// return an error string instead of a blank version, then check that
			return "", errors.New(CancelError)
		}

		// then have check if unique or check if not unique after
		err = checkFunc(cfg, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return input, nil

	}
}

// quite possibly used one place but an alternative to making the prompts generic ie insaner than necessary.
// returns -1 on exit.
// TODO: return a specific error that informs the function to exit
func getIntPrompt(prompt string) (int, error) {
	fmt.Println(prompt + " or exit to cancel")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			// -1 instead of the 0 value for an int because checking to exit
			return 0, errors.New(CancelError)
		}

		output, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return output, nil
	}
}

// pass into getPrompt functions when no checks need to be done
func checkFuncNil(cfg *Config, s string) error {
	// look into optional 1st order func params
	return nil
}

// prints prompt, takes an input from the user, then runs check function to fetch a struct from the db.
// get func is something like getXXXStruct(). Returns an empty struct of T if "exit" or "cancel" entered
// Uses generics! It's cool! I did not want 10 different functions to return 5 different structs!
// TODO: return a specific error that informs the function to exit instead of checking for nil
func getStructPrompt[T any](cfg *Config, prompt string, getFunc func(*Config, string) (T, error)) (T, error) {
	fmt.Println(prompt + " or exit to cancel")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}
		input := strings.TrimSpace(text)
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			var nilT T
			return nilT, errors.New(CancelError)
		}

		// then have check if unique or check if not unique after
		output, err := getFunc(cfg, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return output, nil

	}
}
