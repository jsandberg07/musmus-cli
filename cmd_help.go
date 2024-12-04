package main

import "fmt"

// for use in a command like card activation or add investigator
// prints available commands or flags
func cmdHelp(scmdMap map[string]Flag) error {
	for _, key := range scmdMap {
		fmt.Printf("* %s\n", key.symbol)
		fmt.Print(key.description)
		if key.takesValue {
			fmt.Print(" Requires value.")
		}
		fmt.Println()
	}
	return nil
}
