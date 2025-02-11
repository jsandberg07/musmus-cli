package main

import (
	"bufio"
	"fmt"
	"os"
)

/* removed to make goto more consistent with other commands, and remove unused parameters from literally every command
// separate out flags for consistency maybe
func getSetStateCmd() Command {

	gotoFlags := make(map[string]Flag)

	ccFlag := Flag{
		symbol:      "cc",
		description: "Goes to CC processing menu.",
		takesValue:  false,
	}
	gotoFlags["-"+ccFlag.symbol] = ccFlag

	psFlag := Flag{
		symbol:      "ps",
		description: "Goes to positions menu.",
		takesValue:  false,
	}
	gotoFlags["-"+psFlag.symbol] = psFlag

	iFlag := Flag{
		symbol:      "in",
		description: "Goes to investigator menu.",
		takesValue:  false,
	}
	gotoFlags["-"+iFlag.symbol] = iFlag

	prFlag := Flag{
		symbol:      "pr",
		description: "Goes to protocol menu.",
		takesValue:  false,
	}
	gotoFlags["-"+prFlag.symbol] = prFlag

	seFlag := Flag{
		symbol:      "se",
		description: "Goes to settings menu.",
		takesValue:  false,
	}
	gotoFlags["-"+seFlag.symbol] = seFlag

	stFlag := Flag{
		symbol:      "st",
		description: "Goes to the strains menu.",
		takesValue:  false,
	}
	gotoFlags["-"+stFlag.symbol] = stFlag

	quFlag := Flag{
		symbol:      "qu",
		description: "Goes to the queries menu.",
		takesValue:  false,
	}
	gotoFlags["-"+quFlag.symbol] = quFlag

	orFlag := Flag{
		symbol:      "or",
		description: "Goes to the orders menu.",
		takesValue:  false,
	}
	gotoFlags["-"+orFlag.symbol] = orFlag

	rmFlag := Flag{
		symbol:      "rm",
		description: "Goes to the reminders menu.",
		takesValue:  false,
	}
	gotoFlags["-"+rmFlag.symbol] = rmFlag

	gotoCmd := Command{
		name:        "goto",
		description: "Goes to another menu.",
		function:    gotoCommand,
		flags:       gotoFlags,
	}

	return gotoCmd
}


func gotoCommand(cfg *Config) error {
	if len(args) == 0 {
		return errors.New("goto requires a state flag")
	}
	if len(args) != 1 {
		return errors.New("goto only takes 1 flag")
	}

	switch args[0].flag {
	case "-cc":
		cfg.nextState = getProcessingState()
	case "-ps":
		cfg.nextState = getPositionState()
	case "-in":
		cfg.nextState = getInvesitatorsState()
	case "-pr":
		cfg.nextState = getProtocolState()
	case "-se":
		cfg.nextState = getSettingsState()
	case "-st":
		cfg.nextState = getStrainsState()
	case "-qu":
		cfg.nextState = getQueriesState()
	case "-or":
		cfg.nextState = getOrdersState()
	case "-rm":
		cfg.nextState = getRemindersState()
	default:
		return errors.New("whoops a fake flag slipped into gotoCommand")
	}

	return nil
}
*/

func getGotoCmd() Command {
	gotoFlags := make(map[string]Flag)
	gotoCmd := Command{
		name:        "goto",
		description: "Used for changing menus",
		function:    gotoFunction,
		flags:       gotoFlags,
	}

	return gotoCmd
}

func getGotoFlags() map[string]Flag {
	gotoFlags := make(map[string]Flag)

	ccFlag := Flag{
		symbol:      "cc",
		description: "Goes to CC processing menu.",
		takesValue:  false,
	}
	gotoFlags["-"+ccFlag.symbol] = ccFlag

	psFlag := Flag{
		symbol:      "ps",
		description: "Goes to positions menu.",
		takesValue:  false,
	}
	gotoFlags["-"+psFlag.symbol] = psFlag

	iFlag := Flag{
		symbol:      "in",
		description: "Goes to investigator menu.",
		takesValue:  false,
	}
	gotoFlags["-"+iFlag.symbol] = iFlag

	prFlag := Flag{
		symbol:      "pr",
		description: "Goes to protocol menu.",
		takesValue:  false,
	}
	gotoFlags["-"+prFlag.symbol] = prFlag

	seFlag := Flag{
		symbol:      "se",
		description: "Goes to settings menu.",
		takesValue:  false,
	}
	gotoFlags["-"+seFlag.symbol] = seFlag

	stFlag := Flag{
		symbol:      "st",
		description: "Goes to the strains menu.",
		takesValue:  false,
	}
	gotoFlags["-"+stFlag.symbol] = stFlag

	quFlag := Flag{
		symbol:      "qu",
		description: "Goes to the queries menu.",
		takesValue:  false,
	}
	gotoFlags["-"+quFlag.symbol] = quFlag

	orFlag := Flag{
		symbol:      "or",
		description: "Goes to the orders menu.",
		takesValue:  false,
	}
	gotoFlags["-"+orFlag.symbol] = orFlag

	rmFlag := Flag{
		symbol:      "rm",
		description: "Goes to the reminders menu.",
		takesValue:  false,
	}
	gotoFlags["-"+rmFlag.symbol] = rmFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available flags",
		takesValue:  false,
	}
	gotoFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without changing menu",
		takesValue:  false,
	}
	gotoFlags[exitFlag.symbol] = exitFlag

	return gotoFlags

}

// extremely simple function turned into some weird loop, all to remove a parameter from all the other functions
func gotoFunction(cfg *Config) error {
	// get flags
	flags := getGotoFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the flag for the menu you'd like to go to. Enter 'help' to see list of available flags")
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

		if len(args) > 1 {
			fmt.Println("Too many flags entered. Only really need one.")
		}

		for _, arg := range args {
			switch arg.flag {
			case "-cc":
				cfg.nextState = getProcessingState()
				exit = true
			case "-ps":
				cfg.nextState = getPositionState()
				exit = true
			case "-in":
				cfg.nextState = getInvesitatorsState()
				exit = true
			case "-pr":
				cfg.nextState = getProtocolState()
				exit = true
			case "-se":
				cfg.nextState = getSettingsState()
				exit = true
			case "-st":
				cfg.nextState = getStrainsState()
				exit = true
			case "-qu":
				cfg.nextState = getQueriesState()
				exit = true
			case "-or":
				cfg.nextState = getOrdersState()
				exit = true
			case "-rm":
				cfg.nextState = getRemindersState()
				exit = true
			case "exit":
				exit = true
			case "help":
				cmdHelp(flags)
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
