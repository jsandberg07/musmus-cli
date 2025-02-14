package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

// prints available command names or flags.
// sorted now! neat!
func cmdHelp(input map[string]Flag) {
	flags := slices.Collect(maps.Values(input))
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].printOrder < flags[j].printOrder
	})
	for _, flag := range flags {
		fmt.Printf("* %s\n", flag.symbol)
		if flag.description != "" {
			fmt.Print(flag.description)
			if flag.takesValue {
				fmt.Print(". Requires value")
			}
			fmt.Println()
		}
	}
}
