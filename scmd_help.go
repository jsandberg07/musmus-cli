package main

import "fmt"

func scmdHelp(scmdMap map[string]Flag) error {
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
