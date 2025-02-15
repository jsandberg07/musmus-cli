package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

// Prints prompt (+ instructions on how to cancel), takes an input from the user, then runs it through the check function for uniqueness or
// if valid entry in the database. Will repeat if checkFunc returns an error. Returns an Error() const that can be checked
// for when the user decides to cancel and should be checked

// functions that require a CheckFunc to parse inputs. An attempt to DRY up code using 1st order functions.
// Did not want 10 very similar functions to return 5 different pieces of data

// prints prompt, takes input, runs through a check (like if a name is unique for investigators or strains) and
// returns the string. Will repeat asking for input if an error occurs in checkFunc.
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

// For when you really want just a number, and to not be able to continue unless you have it
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

func getDatePrompt(prompt string) (time.Time, error) {
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
			return time.Time{}, nil
		}

		output, err := parseDate(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return output, nil

	}
}

// functions that are used as CheckFuncs. Return an error when another attempt at an input is expected.

// used as a CheckFunc. because sometimes you just have two people who have the same name, this triggers a message
// to add a nickname which can be used for signing in instead.
func checkIfInvestigatorNameUnique(cfg *Config, name string) error {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		fmt.Println("Error getting name from DB")
		return err
	}
	if len(investigators) != 0 {
		fmt.Println("Investigator name is not unique. Please consider adding a nickname to both investigators.")
	}
	return nil
}

func checkIfPositionTitleUnique(cfg *Config, input string) error {
	_, err := cfg.db.GetPositionByTitle(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// no position by that name was found
		return nil
	}
	if err != nil {
		// any other DB error so exit
		fmt.Printf("Error checking database for title: %s\n", err)
		return err
	}
	return errors.New("position titles must be unique. Please try again")
}

func checkIfOrderNumberUnique(cfg *Config, input string) error {
	_, err := cfg.db.GetOrderByNumber(context.Background(), input)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}

	if err == nil {
		// is not unique
		return errors.New("order number is not unique. Please try again")
	}

	// is unique
	return nil

}

func checkProtocolExists(cfg *Config, input string) (database.Protocol, error) {
	protocol, err := cfg.db.GetProtocolByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// not found
		return database.Protocol{}, errors.New("no protocol with that number found. Please try again")
	}
	if err != nil {
		// any other error
		return database.Protocol{}, err
	}

	return protocol, nil
}

func checkProtocolUnique(cfg *Config, input string) error {
	protocol, err := cfg.db.GetProtocolByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// input is unique
		return nil
	}
	if err != nil {
		// any other error
		fmt.Println("Error checking db for protocol.")
		return err
	}

	// not unique
	fmt.Printf("Protocol with same number: %s\n", protocol.Title)
	return errors.New("a protocol with that number already exists. Please try again")

}

// does double duty, checks if string is unique protocol # from flag
func getUniqueProtocolFromFlag(cfg *Config, n string) (string, error) {
	_, err := cfg.db.GetProtocolByNumber(context.Background(), n)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting protocols from DB")
		return "", err
	}
	if err == nil {
		// protocol found
		fmt.Println("Protocol with that number already exists. Please try again")
		return "", err
	}

	// if nothing found, then unique and ok
	return n, nil

}

func checkIfStrainCodeUnique(cfg *Config, s string) error {
	_, err := cfg.db.GetStrainByCode(context.Background(), s)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error retrieving data from the DB")
		return err
	}
	if err == nil {
		// strain found, meaning input is not unique
		return errors.New("strain by that ID already exists. Please try again")
	}

	// strain is unique
	return nil

}

// pass into getPrompt functions when no checks need to be done
func checkFuncNil(cfg *Config, s string) error {
	// look into optional 1st order func params
	return nil
}

// prints prompt, takes an input from the user, then runs passed in func to fetch a struct from the db.
// getFuncs are named a variation of getXXXStruct().
// Returns an constant CancelErr if "exit" or "cancel" entered that can be checked for.
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

// functions that are used as getFuncs (ie return a struct)

// will repeat if investigator name is vague
func getInvestigatorStruct(cfg *Config, input string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	// TODO: does returning many even throw the "no rows in result set" error?
	if len(investigators) == 0 {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	if len(investigators) > 1 {
		return database.Investigator{}, errors.New("vague investigator name. Please try again")
	}
	if err != nil {
		// any other error
		return database.Investigator{}, err
	}

	return investigators[0], nil
}

func getProtocolStruct(cfg *Config, input string) (database.Protocol, error) {
	protocol, err := cfg.db.GetProtocolByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Protocol{}, errors.New("protocol not found. please try again")
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		fmt.Println("Error getting strain from DB.")
		return database.Protocol{}, err
	}

	return protocol, nil
}

// currently identical to getInvestivatorStruct, TODO: add PI / can oversee protocol restriction
func getPIStruct(cfg *Config, input string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	if len(investigators) == 0 {
		return database.Investigator{}, errors.New("investigator not found. Please try again")
	}
	if len(investigators) > 1 {
		return database.Investigator{}, errors.New("vague investigator name. Please try again")
	}
	if err != nil {
		// any other error
		return database.Investigator{}, err
	}

	return investigators[0], err
}

func getCageCardStructActive(cfg *Config, input string) (database.CageCard, error) {

	ccid, err := strconv.Atoi(input)
	if err != nil {
		return database.CageCard{}, err
	}
	cc, err := cfg.db.GetCageCardByID(context.Background(), int32(ccid))
	if err != nil {
		return database.CageCard{}, err
	}
	if !cc.ActivatedOn.Valid {
		return database.CageCard{}, errors.New("cage card is not active")
	}
	if cc.DeactivatedOn.Valid {
		return database.CageCard{}, errors.New("cage card is deactivated")
	}

	return cc, nil
}
