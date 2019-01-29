package dresscode

import (
	"encoding/csv"
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	"hash/fnv"
)

const timezone = "Asia/Tokyo"

var location *time.Location

type Holiday struct {
	DateStr string
	Style   string
}

func NewHoliday(month int, day int, style string) Holiday {
	var dateStr string
	switch month {
	case 1:
		dateStr = "Jan " + strconv.Itoa(day)
	case 2:
		dateStr = "Feb " + strconv.Itoa(day)
	case 3:
		dateStr = "Mar " + strconv.Itoa(day)
	case 4:
		dateStr = "Apr " + strconv.Itoa(day)
	case 5:
		dateStr = "May " + strconv.Itoa(day)
	case 6:
		dateStr = "Jun " + strconv.Itoa(day)
	case 7:
		dateStr = "Jul " + strconv.Itoa(day)
	case 8:
		dateStr = "Aug " + strconv.Itoa(day)
	case 9:
		dateStr = "Sep " + strconv.Itoa(day)
	case 10:
		dateStr = "Oct " + strconv.Itoa(day)
	case 11:
		dateStr = "Nov " + strconv.Itoa(day)
	case 12:
		dateStr = "Dec " + strconv.Itoa(day)

	}
	
	return Holiday{DateStr:dateStr, Style:style}
}

type Dresscodes struct {
	Styles []string
	Holidays []Holiday
}

func LoadDresscodes(dresscodeFilename string) Dresscodes {
	location, _ = time.LoadLocation(timezone)
	d := Dresscodes{Styles:[]string{}, Holidays:[]Holiday{}}
	file, err := os.Open(dresscodeFilename)
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
	fmt.Printf("Loaded %d dresscodes\n", len(d.Styles)) 
	return d
}

func (d *Dresscodes) RespondToDresscode(msg string) string {
	// what day is it?
	dateStr := time.Now().In(location).Format("Mon Jan 2 2006")
	if strings.Contains(dateStr, "Oct 31") {
		return "Today's dress code is HALLOWEEN!!!!"
	}
	// compute the hash
	h := hash(dateStr)
	// modulo t
	idx := h % len(d.Styles)
	return fmt.Sprintf("Today's dress code is %v.", d.Styles[idx])
}

func hash(s string) int {
	h := fnv.New32a()
        h.Write([]byte(s))
        return int(h.Sum32())
}
