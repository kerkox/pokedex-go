package main

import (
	"testing"
)

var cases = []struct {
	input    string
	expected []string
}{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	// add more cases here
}

func TestCleanInput(t *testing.T) {

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			if len(actual) != len(c.expected) {
				t.Errorf("For input '%s', expected length %d but got %d", c.input, len(c.expected), len(actual))
			}
			if word != expectedWord {
				t.Errorf("For input '%s', expected word '%s' but got '%s'", c.input, expectedWord, word)
			}
		}
	}    
}