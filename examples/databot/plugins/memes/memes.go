package memes

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// Memes is a lookup table of keywords to meme hyperlinks
// It's populated by a two-column csv with keyword,link
// If you want multiple memes to correspond to a keyword,
// simply add multiple rows
type Memes struct {
	Links map[string][]string
}

func LoadMemes(filename string) Memes {
	m := Memes{Links: make(map[string][]string)}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return m
	}
	csv := csv.NewReader(file)
	for row, err := csv.Read(); row != nil; row, err = csv.Read() {
		if err != nil {
			fmt.Println(err.Error())
			return m
		}
		label := strings.TrimSpace(row[0])
		url := strings.TrimSpace(row[1])
		links, found := m.Links[label]
		if found {
			m.Links[label] = append(links, url)
		} else {
			m.Links[label] = []string{url}
		}
	}
	return m
}

func (m *Memes) RespondToMemes(msg string) string {
	for keyword, links := range m.Links {
		if strings.Contains(msg, keyword) {
			count := len(links)
			if count == 1 {
				return links[0]
			}
			lucky := rand.Intn(count)
			return links[lucky]
		}
	}
	return ""
}
