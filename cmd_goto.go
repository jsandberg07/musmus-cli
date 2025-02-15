package main

import (
	"bufio"
	"fmt"
	"os"
)

func getGotoCmd() Command {
	gotoFlags := make(map[string]Flag)
	gotoCmd := Command{
		name:        "goto",
		description: "Used for changing menus",
		function:    gotoFunction,
		flags:       gotoFlags,
		printOrder:  1,
	}
	return gotoCmd
}

func getGotoFlags() map[string]Flag {
	gotoFlags := make(map[string]Flag)

	cagecardFlag := Flag{
		symbol:      "cagecard",
		description: "",
		takesValue:  false,
		printOrder:  1,
	}
	gotoFlags[cagecardFlag.symbol] = cagecardFlag

	ccFlag := Flag{
		symbol:      "-cc",
		description: "Goes to CC processing menu.",
		takesValue:  false,
		printOrder:  2,
	}
	gotoFlags[ccFlag.symbol] = ccFlag

	positionFlag := Flag{
		symbol:      "position",
		description: "",
		takesValue:  false,
		printOrder:  3,
	}
	gotoFlags[positionFlag.symbol] = positionFlag
	psFlag := Flag{
		symbol:      "-ps",
		description: "Goes to positions menu.",
		takesValue:  false,
		printOrder:  4,
	}
	gotoFlags[psFlag.symbol] = psFlag

	investigatorFlag := Flag{
		symbol:      "investigator",
		description: "",
		takesValue:  false,
		printOrder:  5,
	}
	gotoFlags[investigatorFlag.symbol] = investigatorFlag
	iFlag := Flag{
		symbol:      "-in",
		description: "Goes to investigator menu.",
		takesValue:  false,
		printOrder:  6,
	}
	gotoFlags[iFlag.symbol] = iFlag

	protocolFlag := Flag{
		symbol:      "protocol",
		description: "",
		takesValue:  false,
		printOrder:  7,
	}
	gotoFlags[protocolFlag.symbol] = protocolFlag
	prFlag := Flag{
		symbol:      "-pr",
		description: "Goes to protocol menu.",
		takesValue:  false,
		printOrder:  8,
	}
	gotoFlags[prFlag.symbol] = prFlag

	settingFlag := Flag{
		symbol:      "setting",
		description: "",
		takesValue:  false,
		printOrder:  9,
	}
	gotoFlags[settingFlag.symbol] = settingFlag
	seFlag := Flag{
		symbol:      "-se",
		description: "Goes to settings menu.",
		takesValue:  false,
		printOrder:  10,
	}
	gotoFlags[seFlag.symbol] = seFlag

	strainFlag := Flag{
		symbol:      "strain",
		description: "",
		takesValue:  false,
		printOrder:  11,
	}
	gotoFlags[strainFlag.symbol] = strainFlag
	stFlag := Flag{
		symbol:      "-st",
		description: "Goes to the strains menu.",
		takesValue:  false,
		printOrder:  12,
	}
	gotoFlags[stFlag.symbol] = stFlag

	queryFlag := Flag{
		symbol:      "query",
		description: "",
		takesValue:  false,
		printOrder:  13,
	}
	gotoFlags[queryFlag.symbol] = queryFlag
	quFlag := Flag{
		symbol:      "-qu",
		description: "Goes to the queries menu.",
		takesValue:  false,
		printOrder:  14,
	}
	gotoFlags[quFlag.symbol] = quFlag

	orderFlag := Flag{
		symbol:      "orders",
		description: "",
		takesValue:  false,
		printOrder:  15,
	}
	gotoFlags[orderFlag.symbol] = orderFlag
	orFlag := Flag{
		symbol:      "-or",
		description: "Goes to the orders menu.",
		takesValue:  false,
		printOrder:  16,
	}
	gotoFlags[orFlag.symbol] = orFlag

	reminderFlag := Flag{
		symbol:      "reminder",
		description: "",
		takesValue:  false,
		printOrder:  17,
	}
	gotoFlags[reminderFlag.symbol] = reminderFlag
	rmFlag := Flag{
		symbol:      "-rm",
		description: "Goes to the reminders menu.",
		takesValue:  false,
		printOrder:  18,
	}
	gotoFlags[rmFlag.symbol] = rmFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints list of available flags",
		takesValue:  false,
		printOrder:  100,
	}
	gotoFlags[helpFlag.symbol] = helpFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without changing menu",
		takesValue:  false,
		printOrder:  100,
	}
	gotoFlags[exitFlag.symbol] = exitFlag

	return gotoFlags

}

func gotoFunction(cfg *Config) error {
	flags := getGotoFlags()

	exit := false

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the flag for the menu you'd like to go to. Enter 'help' to see list of available flags")
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
			case "cagecard":
				fallthrough
			case "-cc":
				cfg.nextState = getProcessingState()
				exit = true
			case "position":
				fallthrough
			case "-ps":
				cfg.nextState = getPositionState()
				exit = true
			case "investigator":
				fallthrough
			case "-in":
				cfg.nextState = getInvesitatorsState()
				exit = true
			case "protocol":
				fallthrough
			case "-pr":
				cfg.nextState = getProtocolState()
				exit = true
			case "setting":
				fallthrough
			case "-se":
				cfg.nextState = getSettingsState()
				exit = true
			case "strain":
				fallthrough
			case "-st":
				cfg.nextState = getStrainsState()
				exit = true
			case "query":
				fallthrough
			case "-qu":
				cfg.nextState = getQueriesState()
				exit = true
			case "order":
				fallthrough
			case "-or":
				cfg.nextState = getOrdersState()
				exit = true
			case "reminder":
				fallthrough
			case "-rm":
				cfg.nextState = getRemindersState()
				exit = true
			case "exit":
				exit = true
			case "help":
				cmdHelp(flags)
			default:
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}
		if exit {
			break
		}
	}
	return nil
}
