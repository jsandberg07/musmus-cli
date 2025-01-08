package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// temp file name, part of the refactor

// idea for how i should have created more reusable functions for all the other data types
// more generic get string with a prompt, instead of separate functions for everything
// just pass in "get investigator name to edit" instead of new function just to say "to edit"
// write the same program 3 times and you'll realize what you want you want to refactor

// prints prompt, takes an input from the user, then runs it through the check function for uniqueness or
// if valid entry in the database. Will repeat if checkFunc returns an error. Can probably remove error return
func getStringPrompt(cfg *Config, prompt string, checkFunc func(*Config, string) error) (string, error) {
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
			return "", nil
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

// pass into getPrompt functions when no checks need to be done
func checkFuncNil(cfg *Config, s string) error {
	// look into optional 1st order func params
	return nil
}

// prints prompt, takes an input from the user, then runs check function to fetch a struct from the db.
// get func is something like getXXXStruct(). Returns an empty struct of T if "exit" or "cancel" entered
// Uses generics! It's cool! I did not want 10 different functions to return 5 different structs!
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
			return nilT, nil
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
