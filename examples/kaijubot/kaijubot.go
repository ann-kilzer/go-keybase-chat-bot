package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

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

	rand.Seed(time.Now().Unix())
	client := BuildClient()

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

		response := ProcessMessage(client, msg)

		if err = kbc.SendMessage(msg.Conversation.Id, response); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}

	}
}

// Read the message and decide what to do with it.
// Todo: Add user whitelist
// Handle channels
func ProcessMessage(client *twitter.Client, msg kbc.SubscriptionMessage) string {
	text := msg.Message.Content.Text
	if strings.HasPrefix(text, "kaiju") {
		return GetTokugifsLink(client)
	}
	if strings.HasPrefix(text, "cat") {
		return GetCatsuLink(client)
	}
	return "I do not understand your command."
}
