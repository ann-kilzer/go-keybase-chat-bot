package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	client := BuildClient()
	tweet := GetRandomTweet(client, "tokugifs")
	fmt.Println(tweet.ID)

}

func BuildClient() *twitter.Client {
	accessToken := ReadSecret("config/accessToken.txt")
	accessSecret := ReadSecret("config/accessSecret.txt")
	consumerKey := ReadSecret("config/consumerKey.txt")
	consumerSecret := ReadSecret("config/consumerSecret.txt")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	return twitter.NewClient(httpClient)
}

func ReadSecret(filename string) string {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	key := strings.TrimSpace(string(text))
	return key
}

func ReadRecentTweets(client *twitter.Client, username string) []twitter.Tweet {
	excludeReplies := true
	params := twitter.UserTimelineParams{
		ScreenName:     username,
		ExcludeReplies: &excludeReplies,
		Count:          30,
	}
	tweets, _, err := client.Timelines.UserTimeline(&params)
	if err != nil {
		fmt.Errorf(err.Error())
		return []twitter.Tweet{}
	}

	return tweets
}

func GetRandomTweet(client *twitter.Client, username string) twitter.Tweet {
	tweets := ReadRecentTweets(client, username)
	lucky := rand.Intn(len(tweets))
	return tweets[lucky]
}
