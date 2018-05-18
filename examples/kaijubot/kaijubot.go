package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	apiKey := ReadApiKey("apikey.txt")
	fmt.print(apiKey)

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	sub := kbc.ListenForNewTextMessages()
	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if err = kbc.SendMessage(msg.Conversation.Id, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}

}

func ReadApiKey(filename) string {
     text, err := ioutil.ReadFile(filename)
     if err != nil {
     fmt.Print(err)
     }
     key := strings.TrimSpace(text)
     return key
}


const twitter = "https://twitter.con/"

func GetTweet(username string, apikey) string {
     return ""
}

