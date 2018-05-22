package main

import (
	"testing"
)

func TestIsGreeting(t *testing.T) {
	cases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "Hi",
			Expected: true,
		},
		{
			Input:    "Hi!!!",
			Expected: true,
		},
		{
			Input:    "Durian, king of fruits",
			Expected: false,
		},
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "Hello",
			Expected: true,
		},
		{
			Input:    "Higgldy Piggedy",
			Expected: false,
		},
		{
			Input:    "Konnichiwa",
			Expected: true,
		},
		{
			Input:    "Shiba inu",
			Expected: false,
		},
	}

	for _, c := range cases {
		tIsGreetingHelper(c.Expected, c.Input, t)
	}
}

func tIsGreetingHelper(expected bool, input string, t *testing.T) {
	result := isGreeting(input)
	if expected != result {
		t.Errorf("Result was incorrect on %v, got %v, wanted %v", input, result, expected)
	}
}
