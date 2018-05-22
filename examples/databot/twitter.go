package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Be sure to call rand.Seed beforehand
func GetTokugifsLink(client *twitter.Client) string {
	tweet := GetRandomTweet(client, "tokugifs")
	return ExtractVideo(tweet)

}

func GetCatsuLink(client *twitter.Client) string {
	tweet := GetRandomTweet(client, "catsu")
	return ExtractPhoto(tweet)
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
		Count:          100,
	}
	tweets, _, err := client.Timelines.UserTimeline(&params)
	if err != nil {
		fmt.Errorf(err.Error())
		return []twitter.Tweet{}
	}

	return tweets
}

func GetRandomTweet(client *twitter.Client, username string) *twitter.Tweet {
	tweets := ReadRecentTweets(client, username)
	lucky := rand.Intn(len(tweets))
	return &tweets[lucky]
}

func ExtractVideo(tweet *twitter.Tweet) string {
	for _, ent := range tweet.ExtendedEntities.Media {
		for _, v := range ent.VideoInfo.Variants {
			// just return the first one for now
			return v.URL
		}
	}
	return ""
}

func ExtractPhoto(tweet *twitter.Tweet) string {
	for _, ent := range tweet.ExtendedEntities.Media {
		return ent.MediaURL
	}
	return ""
}
