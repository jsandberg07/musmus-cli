package main

import "errors"

// we settin state
// B^)

func getSetStateCmd() Command {

	setStateFlags := make(map[string]Flag)

	aFlag := Flag{
		symbol:      "a",
		description: "Sets state to CC processing",
		takesValue:  false,
	}
	setStateFlags["-"+aFlag.symbol] = aFlag

	mFlag := Flag{
		symbol:      "m",
		description: "Sets state to main",
		takesValue:  false,
	}
	setStateFlags["-"+mFlag.symbol] = mFlag

	setStateCmd := Command{
		name:        "setstate",
		description: "Sets the state to a different one",
		function:    setStateCommand,
		flags:       setStateFlags,
	}

	return setStateCmd
}

func setStateCommand(cfg *Config, args []Argument) error {
	if len(args) == 0 {
		return errors.New("Set State requires a state flag")
	}
	if len(args) != 1 {
		return errors.New("Set state only takes 1 flag")
	}

	switch args[0].flag {
	case "-a":
		cfg.nextState = getProcessingState()
	case "-m":
		cfg.nextState = getMainState()
	default:
		return errors.New("Whoops a fake flag slipped into setStateCommand")
	}

	return nil
}
