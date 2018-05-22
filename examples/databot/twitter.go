package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func GetTokugifsLink(client *twitter.Client) string {
	tweet := GetRandomTweet(client, "tokugifs")
	if tweet == nil {
		return ""
	}
	return ExtractVideo(tweet)

}

func GetCatsuLink(client *twitter.Client) string {
	tweet := GetRandomTweet(client, "catsu")
	if tweet == nil {
		return ""
	}
	return ExtractPhoto(tweet)
}

func BuildClient(ta *TwitterAuth) *twitter.Client {
	// gotta make sure we get random cats
	rand.Seed(time.Now().Unix())

	config := oauth1.NewConfig(ta.ConsumerKey, ta.ConsumerSecret)
	token := oauth1.NewToken(ta.AccessToken, ta.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	return twitter.NewClient(httpClient)
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
	count := len(tweets)
	if count == 0 {
		return nil
	}
	lucky := rand.Intn(count)
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
