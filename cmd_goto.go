package main

import "errors"

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

	gotoCmd := Command{
		name:        "goto",
		description: "Goes to another menu.",
		function:    gotoCommand,
		flags:       gotoFlags,
	}

	return gotoCmd
}

func gotoCommand(cfg *Config, args []Argument) error {
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
	default:
		return errors.New("whoops a fake flag slipped into gotoCommand")
	}

	return nil
}
