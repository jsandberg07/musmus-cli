package main

import "fmt"

// for use in a command like card activation or add investigator
// prints available commands or flags
// TODO: unordered. add a numeric or enum prio to Flag struct.
// Sort the keys, then iterate through that
// does sort.Ints() work on enums?
func cmdHelp(flags map[string]Flag) error {
	for _, flag := range flags {
		fmt.Printf("* %s\n", flag.symbol)
		fmt.Print(flag.description)
		if flag.takesValue {
			fmt.Print(". Requires value")
		}
		fmt.Println()
	}
	return nil
}
