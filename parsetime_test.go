package main

import (
	"testing"
	"time"
)

type parsingTest struct {
	input     string
	expected  time.Time
	shouldErr bool
}

func TestParseTime(t *testing.T) {
	tests := []parsingTest{
		{
			input:     "11/18/24",
			expected:  time.Date(2024, 11, 18, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			input:     "1/8/24",
			expected:  time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			input:     "1/8/2024",
			expected:  time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			input:     "01/08/24",
			expected:  time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			shouldErr: false,
		},
		{
			input:     "beans",
			expected:  time.Time{},
			shouldErr: true,
		},
	}

	for i, test := range tests {
		output, err := parseDate(test.input)
		// is throwing an error AND should not be throwing an error
		if err != nil && test.shouldErr != true {
			t.Fatalf("Test %v fail: %s", i+1, err)
		}
		if err == nil && test.shouldErr == true {
			t.Fatalf("Test %v SHOULD have failed but passed", i+1)
		}
		if err == nil && test.shouldErr == false {
			if output != test.expected {
				t.Fatalf("Test %v mismatch. E: %v -- O: %v", i+1, test.expected, output)
			}
		}

		// fmt.Printf("%v: %v -- %v\n", i+1, test.expected, output)

	}

}
