package dresscode

import (
	"testing"
)

const otenba = "Otenba"
const ametora = "Ametora"
const GLaS = "Gothic Lolita at Sea"

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
