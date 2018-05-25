package memes

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

// Memes is a lookup table of keywords to meme hyperlinks
// It's populated by a two-column csv with keyword,link
// If you want multiple memes to correspond to a keyword,
// simply add multiple rows
type Memes struct {
	Links map[string][]string
	Files map[string][]string
}

func LoadMemes(filename string) Memes {
	m := Memes{
		Links: make(map[string][]string),
		Files: make(map[string][]string),
	}
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
	m.DownloadAll()

	return m
}

func (m *Memes) DownloadAll() error {
	for label, links := range m.Links {
		for count, url := range links {
			// make a unique name
			name, err := BuildFileName(url, label, count)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			// download it
			err = DownloadURL(name, url)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			// Now keep track of where you put it
			files, found := m.Files[label]
			if found {
				m.Files[label] = append(files, name)
			} else {
				m.Files[label] = []string{name}
			}
		}
	}
	return nil
}

const prefix = "downloads/"

func BuildFileName(url, name string, count int) (string, error) {
	suffix := ParseFileSuffix(url)
	if suffix == "" {
		return "", fmt.Errorf("Bad suffix for %v", url)
	}
	clean := strings.Replace(name, " ", "", -1)
	return fmt.Sprintf("%v%v%d%v", prefix, clean, count, suffix), nil
}

// todo: don't download files we already have
func DownloadURL(name, url string) error {
	fmt.Printf("Downloading %v\n", name)
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ParseFileSuffix(url string) string {
	parts := strings.Split(url, ".")
	last := len(parts) - 1
	return "." + parts[last]
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
