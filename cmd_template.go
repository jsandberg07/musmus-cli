package main

// template for creating new commands
/* commented out because staticcheck hates it. copy and paste a skeleton for new commands

func getXXXCmd() Command {
	XXXFlags := make(map[string]Flag)
	XXXCmd := Command{
		name:        "XXX",
		description: "Used for XXX",
		function:    XXXFunction,
		flags:       XXXFlags,
		printOrder:  1,
	}

	return XXXCmd
}

func getXXXFlags() map[string]Flag {
	XXXFlags := make(map[string]Flag)
	XFlag := Flag{
		symbol:      "-X",
		description: "Sets X",
		takesValue:  false,
		printOrder:  1,
	}
	XXXFlags[XFlag.symbol] = XFlag

	fmt.Println("If you see this, you accidentally left the template flag function in")
	return XXXFlags

}

func XXXFunction(cfg *Config) error {
	// get flags
	flags := getXXXFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// loop
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

		// do between loop behavior here

		// regular parsing
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
				fmt.Printf("%s%s\n", DefaultFlagMsg, arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

*/
