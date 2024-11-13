package main

type Flag struct {
	symbol      string
	description string
	takesValue  bool
}

type Argument struct {
	flag  string
	value string
}

type Command struct {
	name        string
	description string
	function    func(args []Argument) error
	flags       map[string]Flag
}
