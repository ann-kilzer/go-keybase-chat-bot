package dresscode

import (
	"encoding/csv"
	"os"
	"fmt"
	"time"
	"strings"
	"hash/fnv"
)


type Dresscodes struct {
	Styles []string
}

func LoadDresscodes(filename string) Dresscodes {
	d := Dresscodes{Styles:[]string{}}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return d
	}
	csv := csv.NewReader(file)
	for row, err := csv.Read(); row != nil; row, err = csv.Read() {
		if err != nil {
			fmt.Println(err.Error())
			return d
		}
		style := strings.TrimSpace(row[0])
		d.Styles = append(d.Styles, style)
	}
	return d
}

func (d *Dresscodes) RespondToDresscode(msg string) string {
	// what day is it?
	datestr := time.Now().Format("Mon Jan 2 2006")
	fmt.Println(datestr)
	// compute the hash
	h := hash(datestr)
	// modulo t
	idx := h % len(d.Styles)
	return fmt.Sprintf("Today's dress code is %v.", d.Styles[idx])
}

func hash(s string) int {
	h := fnv.New32a()
        h.Write([]byte(s))
        return int(h.Sum32())
}
