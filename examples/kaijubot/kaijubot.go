package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type Chatbot struct {
     Mux sync.Mutex
     Location string
     Kbc *kbchat.API
     Client *twitter.Client
}

// make him a real boy
func InitChatbot() *Chatbot {
	rand.Seed(time.Now().Unix())

	var err error
	var kbc *kbchat.API
	var kbLoc string
	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	return &Chatbot{
	       Location: kbLoc,
	       Kbc: kbc,
	       Client: BuildClient(),
	}
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	bot := InitChatbot()

	sub := bot.Kbc.ListenForNewTextMessages()
	for {
		bot.Mux.Lock()
		msg, err := sub.Read()
		bot.Mux.Unlock()

		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		bot.Respond(msg)
	}
}

func (bot *Chatbot) Respond(msg kbchat.SubscriptionMessage) {
        bot.Mux.Lock()
	response := ProcessMessage(bot.Client, msg)
	if response != "" {	   	
		fmt.Printf("Sending response %v", response)
		if err := bot.Kbc.SendMessage(msg.Conversation.Id, response); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
	bot.Mux.Unlock()
}

// Read the message and decide what to do with it.
// Todo: Add user whitelist
// Handle channels
func ProcessMessage(client *twitter.Client, msg kbchat.SubscriptionMessage) string {
	text := msg.Message.Content.Text.Body
	fmt.Printf("Handling %v...\n", text)

	if strings.HasPrefix(text, "help") {
		return "You can ask me things like 'kaiju' or 'cat'"
	}
	if strings.HasPrefix(text, "kaiju") {
		return GetTokugifsLink(client)
	}
	if strings.HasPrefix(text, "cat") {
		return GetCatsuLink(client)
	}
	return ""
}
