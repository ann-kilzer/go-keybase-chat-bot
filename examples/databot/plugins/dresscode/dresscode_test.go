package dresscode

import (
	"testing"
)

const otenba = "Otenba"
const ametora = "Ametora"
const GLaS = "Gothic Lolita at Sea"

func TestNewHoliday(t *testing.T) {
	cases := []struct{
		month int
		day int
		style string
		expected Holiday
	} {
		{
			month: 1,
			day: 24,
			style: "Shibukaji",
			expected: Holiday{
				DateStr: "Jan 24",
				Style: "Shibukaji",
			},
		},
	}

	for c := range cases {
		h := NewHoliday(c.month, c.day, c.style)
		if h.DateStr != c.expected.DateStr {
			t.Fail()
		}
		if h.Style != c.expected.Style {
			t.Fail()
		}
	}
}

func TestDresscodes(t *testing.T) {
	d := Dresscodes{
		Styles : []string{
			otenba,
			ametora,
			GLaS,
		},
	}


	resp1 := d.RespondToDresscode("anything")
	resp2 := d.RespondToDresscode("anything")
	if resp1 != resp2 {
		t.Fail()
	}
}
