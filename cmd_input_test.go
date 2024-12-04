package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	fmt.Println("Parse test //")
	flags := getActivationFlags()
	input := "-d 12/24/24 help exit"
	inputs, err := readSubcommandInput(input)
	if err != nil {
		t.Fatalf("Error readSubcommandInput: %s\n", err)
	}
	fmt.Println(inputs)

	args, err := parseSubcommand(flags, inputs)
	if err != nil {
		t.Fatalf("Error parseSubcommand: %s", err)
	}
	fmt.Println(args)
}
