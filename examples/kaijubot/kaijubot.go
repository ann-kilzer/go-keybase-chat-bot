package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
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

		link := GetTokugifsLink()
		
		if err = kbc.SendMessage(msg.Conversation.Id, link); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}

	}
}
