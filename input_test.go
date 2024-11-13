package main

import (
	"testing"
)

func TestInput(t *testing.T) {
	input := "print c meowmeow\n"
	cmdName, arguments, _ := getInput(input)
	eCmdName := "print"
	eArguments := []Argument{{flag: "c", value: "meowmeow"}}

	if cmdName != eCmdName {
		t.Fatalf("Argument Name: %s -- %s", cmdName, eCmdName)
	}

	if arguments[0].flag != eArguments[0].flag {
		t.Fatalf("Argument Flag: %s -- %s", arguments[0].flag, eArguments[0].flag)
	}

	if arguments[0].value != eArguments[0].value {
		t.Fatalf("Argument Flag: %s -- %s", arguments[0].value, eArguments[0].value)
	}
}
