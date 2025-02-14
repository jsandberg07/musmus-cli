package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	// fmt.Println("Parse test //")
	flags := getActivationFlags()
	input := "-d 12/24/24 help exit"
	inputs, err := readSubcommandInput(input)
	if err != nil {
		t.Fatalf("Error readSubcommandInput: %s\n", err)
	}
	// fmt.Println(inputs)

	_, err = parseArguments(flags, inputs)
	if err != nil {
		t.Fatalf("Error parseSubcommand: %s", err)
	}
	// fmt.Println(args)
}

func TestParseProtocol(t *testing.T) {
	// fmt.Println("Parsing with - //")
	flags := getAddInvestToProtFlags()
	input := "-p 12-24-32"
	inputs, err := readSubcommandInput(input)
	if err != nil {
		t.Fatalf("Error readSubcommandInput: %s\n", err)
	}
	// fmt.Println(inputs)

	_, err = parseArguments(flags, inputs)
	if err != nil {
		t.Fatalf("Error parseSubcommand: %s", err)
	}
	// fmt.Println(args)
}
