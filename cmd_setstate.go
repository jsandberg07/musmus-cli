package main

import "errors"

// we settin state
// B^)

func getSetStateCmd() Command {

	setStateFlags := make(map[string]Flag)

	aFlag := Flag{
		symbol:      "a",
		description: "Sets state to CC processing.",
		takesValue:  false,
	}
	setStateFlags["-"+aFlag.symbol] = aFlag

	iFlag := Flag{
		symbol:      "i",
		description: "Sets state to investigator.",
		takesValue:  false,
	}
	setStateFlags["-"+iFlag.symbol] = iFlag

	setStateCmd := Command{
		name:        "setstate",
		description: "Sets the state to a different one.",
		function:    setStateCommand,
		flags:       setStateFlags,
	}

	return setStateCmd
}

func setStateCommand(cfg *Config, args []Argument) error {
	if len(args) == 0 {
		return errors.New("setstate requires a state flag")
	}
	if len(args) != 1 {
		return errors.New("setstate only takes 1 flag")
	}

	switch args[0].flag {
	case "-a":
		cfg.nextState = getProcessingState()
	case "-m":
		cfg.nextState = getMainState()
	case "-i":
		cfg.nextState = getInvesitatorsState()
	default:
		return errors.New("whoops a fake flag slipped into setStateCommand")
	}

	return nil
}
