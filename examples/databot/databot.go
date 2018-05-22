package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type Chatbot struct {
	Mux          sync.Mutex
	Location     string
	Kbc          *kbchat.API
	Client       *twitter.Client
	AllowedUsers []string // friends we accept messages from
}

// make data a real boy
func InitChatbot() *Chatbot {
	config := ReadConfig("config/config.toml")

	var err error
	var kbc *kbchat.API
	var kbLoc string
	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	return &Chatbot{
		Location:     kbLoc,
		Kbc:          kbc,
		Client:       BuildClient(&config.twitter),
		AllowedUsers: config.whitelist.AllowedUsers,
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
	text := strings.ToLower(msg.Message.Content.Text.Body)
	username := msg.Message.Sender.Username
	fmt.Printf("Handling %v...\n", text)

	if isGreeting(text) {
		return fmt.Sprintf("Hello %v!", username)
	}
	if strings.HasPrefix(text, "who are you") {
		return "I am an android."
	}
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

func isGreeting(text string) bool {
	lower := strings.ToLower(text)
	re, _ := regexp.Compile("(hi[^a-z]+.*)|(^hi$)|(hello)|(konnichiwa)")
	return re.MatchString(lower)
}
