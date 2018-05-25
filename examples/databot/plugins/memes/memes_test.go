package memes

import (
	"testing"
)

func TestParseFileSuffix(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "https://media.giphy.com/media/3bvmkU1dGHLTa/giphy.gif",
			Expected: ".gif",
		},
		{
			Input:    "http://www.developermemes.com/wp-content/uploads/2013/12/Worked-Fine-In-Dev-Ops-Problem-Now.jpg",
			Expected: ".jpg",
		},
	}

	for _, c := range cases {
		tParseFilenameFromURLHelper(c.Expected, c.Input, t)
	}
}

func tParseFilenameFromURLHelper(expected, input string, t *testing.T) {
	actual := ParseFileSuffix(input)
	if actual != expected {
		t.Errorf("Wanted %v, Got %v", expected, actual)
	}
}
