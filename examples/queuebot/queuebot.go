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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type QueueBot struct {
	Kbc    *kbchat.API
	Config *TomlConfig
}

type TomlConfig struct {
	ConvoID  string `toml:"conversation_id"`
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

func ReadConfig(filename string) *TomlConfig {
	var conf TomlConfig
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		fail("Unable to read config, RIP")
		return nil
	}
	return &conf
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	qb := InitQueueBot("queuebot.toml")

	for {
		HandleNewMessges(qb)
		time.Sleep(time.Minute)
	}
}

func HandleNewMessges(qb *QueueBot) {
	// list directory contents
	messages := MockListDir(qb.Config.QueueDir)
	for _, msg := range messages {
		fmt.Printf("Handling message %v\n", msg)
		SendMessage(msg, qb)
	}

}

// todo: replace this with something that does ls when I get wifi and can read godocs
func MockListDir(dir string) []string {
	return []string{"a", "b", "c"}
}

func SendMessage(filename string, qb *QueueBot) {
	contents, err := ReadMessage(filename)
	if err != nil {
		fail("unable to read message", err.Error())
	}

	if contents == "" {
		fmt.Fprintf(os.Stderr, "Contents empty, nothing to do\n")
		return
	}

	// test logging
	fmt.Printf("sending %v\n", contents)
	if err = qb.Kbc.SendMessage(qb.Config.ConvoID, contents); err != nil {
		fail("error echo'ing message: %s", err.Error())
	}

	// clean up your mess
	err = os.Remove(filename)
	if err != nil {
		fail("Unable to clean up file %v", err.Error())
	}
}

// todo: review this when you can internet
func ReadMessage(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	return string(bytes), err
}
