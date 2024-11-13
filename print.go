package main

func getPrintCmd() Command {

	printFlags := make(map[string]Flag)

	cFlag := Flag{
		symbol:      "c",
		description: "Allows custom text.",
		takesValue:  true,
	}
	printFlags["-"+cFlag.symbol] = cFlag

	bFlag := Flag{
		symbol:      "b",
		description: "Makes text uppercase.",
		takesValue:  false,
	}
	printFlags["-"+bFlag.symbol] = bFlag

	printCmd := Command{
		name:        "print",
		description: "prints wow or sometimes something other than wow.",
		function:    printCommand,
		flags:       printFlags,
	}

	return printCmd
}
