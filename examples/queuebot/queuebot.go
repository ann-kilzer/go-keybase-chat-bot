// This good good boy reads from a message queue and delivers messages to a preconfigured destination
// useful for when keybase mysteriously fails to send via your shell script ¯\_(ツ)_//
// idk
//
// * inspired by smtp mailers :) *
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func main() {
	var configfile string
	if len(os.Args) >= 2 {
		configfile = os.Args[1]
	} else {
		configfile = "queuebot.toml"
	}

	qb := InitQueueBot(configfile)
	for {
		HandleNewMessges(qb)
		time.Sleep(time.Minute)
	}
}

type QueueBot struct {
	Kbc    *kbchat.API
	Config *TomlConfig
}

type TomlConfig struct {
	TeamName string `toml:"team_name"`
	Channel  string `toml:"channel"`
	QueueDir string `toml:"queue_directory"`
}

func InitQueueBot(filename string) *QueueBot {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	return &QueueBot{
		Kbc:    kbc,
		Config: ReadConfig(filename),
	}
}

func (qb *QueueBot) Send(contents string) error {
	if err := qb.Kbc.SendMessageByTeamName(qb.Config.TeamName, contents, &qb.Config.Channel); err != nil {
		return err
	}
	return nil
}

func ReadConfig(filename string) *TomlConfig {
	var conf TomlConfig
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		fail("Unable to read config %v, RIP", filename)
		return nil
	}
	return &conf
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

// hail mary, try to alert the channel
func failWithAlert(qb *QueueBot, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	qb.Send("Bot down! " + msg)
	os.Exit(3)
}

func HandleNewMessges(qb *QueueBot) {
	dir := qb.Config.QueueDir
	msgFiles, err := ListDir(dir)
	if err != nil {
		failWithAlert(qb, "unable to read directory "+dir, err.Error())
	}
	for _, mf := range msgFiles {
		fmt.Printf("Handling message %v\n", mf)
		fullpath := filepath.Join(dir, mf)
		SendMessage(fullpath, qb)
	}

}

func ListDir(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return []string{}, err
	}
	return names, nil
}

// SendMessage reads the contents at filepath and sends it to the
// Queuebot's keybase channel
func SendMessage(filename string, qb *QueueBot) {
	contents, err := ReadMessage(filename)
	if err != nil {
		failWithAlert(qb, "unable to read message", err.Error())
	}

	if contents == "" {
		fmt.Fprintf(os.Stderr, "Contents empty, nothing to do\n")
		return
	}

	// test logging
	fmt.Printf("sending %v\n", contents)
	err = qb.Send(contents)
	if err != nil {
		failWithAlert(qb, "error sending message: %s", err.Error())
	}

	// clean up your mess
	err = os.Remove(filename)
	if err != nil {
		failWithAlert(qb, "Unable to clean up file %v", err.Error())
	}
}

func ReadMessage(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	return string(bytes), err
}
