package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/ann-kilzer/go-keybase-chat-bot/examples/databot/plugins/memes"
	"github.com/ann-kilzer/go-keybase-chat-bot/examples/databot/plugins/tweets"
	"github.com/ann-kilzer/go-keybase-chat-bot/examples/databot/plugins/dresscode"
	"github.com/ann-kilzer/go-keybase-chat-bot/examples/databot/toml"
	"github.com/ann-kilzer/go-keybase-chat-bot/kbchat"
)

type Chatbot struct {
	Mux      sync.Mutex
	Location string
	Kbc      *kbchat.API
	Friends  map[string]bool // friends we accept messages from
	Memes    memes.Memes
	Tweets   *tweets.TweetResponder
	Dresscode dresscode.Dresscodes
}

// make data a real boy
func InitChatbot() *Chatbot {
	config := toml.ReadConfig("config/config.toml")

	var err error
	var kbc *kbchat.API
	var kbLoc string
	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	// build whitelist set
	friends := make(map[string]bool)
	for _, user := range config.Whitelist.AllowedUsers {
		friends[user] = true
	}

	return &Chatbot{
		Location: kbLoc,
		Kbc:      kbc,
		Friends:  friends,
		Memes:    memes.LoadMemes("config/memes.csv"),
		Tweets:   tweets.NewTweetResponder(&config.Twitter, "config/tweets.csv"),
		Dresscode: dresscode.LoadDresscodes("config/dresscodes.csv"),
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
		} else { 
			bot.Respond(msg)
		}
	}
}

// the locking here is really overkill now that this is single threaded.
// Still, something is up with the channels, because occasionally things deadlock
// augh this is so wrong it hurts
func (bot *Chatbot) Respond(msg kbchat.SubscriptionMessage) {
	bot.Mux.Lock()
	response := ProcessMessage(bot, msg)
	if response != "" {
		fmt.Printf("Sending response %v\n", response)
		if err := bot.Kbc.SendMessage(msg.Conversation.Id, response); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
	bot.Mux.Unlock()
}

// Read the message and decide what to do with it.
func ProcessMessage(bot *Chatbot, msg kbchat.SubscriptionMessage) string {
	text := strings.ToLower(msg.Message.Content.Text.Body)
	username := msg.Message.Sender.Username
	// check if the user is a friend
	_, ok := bot.Friends[username]
	if !ok {
		fmt.Printf("Ignoring message %v from stranger %v", text, username)
		return ""
	}

	fmt.Printf("Handling %v from %v\n", text, username)

	if isGreeting(text) {
		return fmt.Sprintf("Hello %v!", username)
	}
	if strings.HasPrefix(text, "who are you") {
		return GetRandomResponse(
			"I am an android.",
			"My name is Lt. Commander Data.")
	}
	if strings.HasPrefix(text, "help") {
		return "You can ask me things like 'kaiju' or 'cat'"
	}
	if strings.Contains(text, "dresscode") {
		return bot.Dresscode.RespondToDresscode(text)
	}
	if resp := bot.Tweets.RespondToKeywords(text); resp != "" {
		return resp
	}

	return bot.Memes.RespondToMemes(text)
}

func isGreeting(text string) bool {
	lower := strings.ToLower(text)
	re, _ := regexp.Compile("(hi[^a-z]+.*)|(^hi$)|(hello)|(konnichiwa)")
	return re.MatchString(lower)
}

func GetRandomResponse(responses ...string) string {
	count := len(responses)
	if count == 0 {
		return ""
	}
	return responses[rand.Intn(count)]
}
