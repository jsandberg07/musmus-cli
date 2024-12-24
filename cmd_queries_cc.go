package main

import (
	"bufio"
	"fmt"
	"os"
)

func getCCQueriesCmd() Command {
	CCQueriesFlags := make(map[string]Flag)
	CCQueriesCmd := Command{
		name:        "cc",
		description: "Run queries on cage cards",
		function:    CCQueriesFunction,
		flags:       CCQueriesFlags,
	}

	return CCQueriesCmd
}

// [i]nvestigator, [p]rotocol, [d]ate, active by default, [de]activated, [all],
// dag how do i wanna do this
// thinking is king innit?
// maybe prompts
// like quick ones >all active cards
// or limit it like just by PI or ivnestigator. its a csv
// easier to short than write and parse a bunch of unique queries
// problem is there are so many not nulls that need to have a match
// it would be easy if protocol could be null
// or investigator could be null, then it could be optional that way
// param structs
// where its like investigator, protocol, activated on, deactivated on
// then have a bunch of cases for which query to run
// no no
// investigator, protocol, PI are optional. you set and case those.
// check if that protocol is under that PI and eliminate one query
// dont do that for investigator, in case they are removed from that protocol
// then you can do activated on, deactivated as optional like null or not null
// depending on if you want active, all, deactivated
// how to do for particular dates? fuck that, that's all narrow it down later on your own time
// fine this is ok

func getCCQueriesFlags() map[string]Flag {
	XXXFlags := make(map[string]Flag)
	XFlag := Flag{
		symbol:      "X",
		description: "Sets X",
		takesValue:  false,
	}
	XXXFlags["-"+XFlag.symbol] = XFlag

	// ect as needed or remove the "-"+ for longer ones

	fmt.Println("If you see this, you accidentally left the template flag function in")
	return XXXFlags

}

// look into removing the args thing, might have to stay
func CCQueriesFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getXXXFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "-X":
				exit = true
			default:
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}
