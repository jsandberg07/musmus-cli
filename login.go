package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/jsandberg07/clitest/internal/auth"
	"github.com/jsandberg07/clitest/internal/database"
	"golang.org/x/term"
)

func (cfg *Config) login() error {
	inv, err := getStructPrompt(cfg, "Enter name", getInvestigatorStruct)
	if err != nil {
		return err
	}
	nilInv := database.Investigator{}
	if inv == nilInv {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	if !inv.HashedPassword.Valid {
		password, err := getNewPassword()
		if err != nil {
			return err
		}
		hash, err := auth.HashPassword(password)
		if err != nil {
			return err
		}
		uhpp := database.UpdateHashedPasswordParams{
			ID:             inv.ID,
			HashedPassword: sql.NullString{Valid: true, String: hash},
		}
		err = cfg.db.UpdateHashedPassword(context.Background(), uhpp)
		if err != nil {
			return err
		}

		fmt.Println("Password has been updated. Please login again")
		os.Exit(0)
	}

	tries := 0
	for {
		password, err := getSecureInput("Enter password")
		if err != nil {
			return err
		}
		if password == "" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		err = auth.CheckPasswordHash(password, inv.HashedPassword.String)
		if err != nil && strings.Contains(err.Error(), "is not the hash of the given password") {
			// wrong password
			tries++
			if tries == 3 {
				fmt.Println("Too many attempts. Exiting...")
				os.Exit(0)
			} else {
				fmt.Println("Incorrect password. Please try again.")
				continue
			}
		}
		if err != nil {
			// any other error
			return err
		}
		// correct
		break
	}

	position, err := cfg.db.GetUserPosition(context.Background(), inv.Position)
	if err != nil {
		fmt.Println("Error getting logged in user position")
		return err
	}

	cfg.loggedInInvestigator = &inv
	cfg.loggedInPosition = &position

	return nil

}

func (cfg *Config) printLogin() {
	fmt.Printf("Logged in as %s -- %s\n", cfg.loggedInInvestigator.IName, cfg.loggedInPosition.Title)
}

func (cfg *Config) createAdmin() error {
	p, err := cfg.db.CreateAdminPosition(context.Background())
	if err != nil {
		fmt.Println("Error creating admin position")
		return err
	}
	hash, err := auth.HashPassword("admin")
	if err != nil {
		fmt.Println("Error hashing password")
		return err
	}
	caip := database.CreateAdminInvestigatorParams{
		Position:       p.ID,
		HashedPassword: sql.NullString{Valid: true, String: hash},
	}
	i, err := cfg.db.CreateAdminInvestigator(context.Background(), caip)
	if err != nil {
		fmt.Println("Error creating admin staff")
		return err
	}
	if verbose {
		fmt.Println(i)
	}
	return nil
}

// TODO: can we make this show stars instead of the password?
func getNewPassword() (string, error) {
	fmt.Println("No password found")
	for {
		password, err := getSecureInput("Enter password")
		if err != nil {
			return "", err
		}
		if string(password) == "" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		confirm, err := getSecureInput("Confirm password")
		if err != nil {
			return "", err
		}
		if password == "" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		if password != confirm {
			fmt.Println("Passwords did not match. Please try again")
			continue
		}

		return password, nil

	}
}

// Gets input (passwords) without displaying it on the CLI for security
func getSecureInput(prompt string) (string, error) {
	fmt.Println(prompt + " or exit to cancel")
	for {
		fmt.Print("> ")
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		input := strings.TrimSpace(string(bytes))
		if input == "" {
			fmt.Println("No input found. Please try again.")
			continue
		}
		if input == "exit" || input == "cancel" {
			return "", nil
		}
		return input, nil
	}
}

/* removed because it was justed used for logging in, identical to getInv(cfg, name)(inv, err) anyway
func (cfg *Config) getInvestigator(name string) (database.Investigator, error) {
	investigators, err := cfg.db.GetInvestigatorByName(context.Background(), name)
	if err != nil {
		fmt.Println("Error getting investigator")
		return database.Investigator{}, nil
	}
	if len(investigators) > 1 {
		fmt.Println("Error getting investigator")
		return database.Investigator{}, errors.New("vague investigator name")
	}

	return investigators[0], nil
}
*/
