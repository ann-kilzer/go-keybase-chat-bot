package tweets

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ann-kilzer/go-keybase-chat-bot/examples/databot/toml"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TweetResponder struct {
	Client   *twitter.Client
	Accounts []WatchedAccount
}

type WatchedAccount struct {
	Keyword string
	Account string
	Type    string // is there enum?
}

func (t *TweetResponder) GetVideoLink(username string) string {
	tweet := t.GetRandomTweet(username)
	if tweet == nil {
		return ""
	}
	return ExtractVideo(tweet)

}

func (t *TweetResponder) GetPictureLink(username string) string {
	tweet := t.GetRandomTweet(username)
	if tweet == nil {
		return ""
	}
	return ExtractPhoto(tweet)
}

func (t *TweetResponder) GetText(username string) string {
	tweet := t.GetRandomTweet(username)
	if tweet == nil {
		return ""
	}
	return tweet.Text
}

func NewTweetResponder(ta *toml.TwitterAuth, configFile string) *TweetResponder {
	// gotta make sure we get random cats
	rand.Seed(time.Now().Unix())

	config := oauth1.NewConfig(ta.ConsumerKey, ta.ConsumerSecret)
	token := oauth1.NewToken(ta.AccessToken, ta.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	return &TweetResponder{
		Client:   twitter.NewClient(httpClient),
		Accounts: LoadConfig(configFile),
	}
}

// We expect a three column csv:
// keyword,username,type (one of "text","pic", or "video")
func LoadConfig(filename string) []WatchedAccount {
	a := []WatchedAccount{}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return a
	}

	csv := csv.NewReader(file)
	for row, err := csv.Read(); row != nil; row, err = csv.Read() {
		if err != nil {
			fmt.Println(err.Error())
			return a
		}
		a = append(a, WatchedAccount{
			Keyword: strings.TrimSpace(row[0]),
			Account: strings.TrimSpace(row[1]),
			Type:    strings.TrimSpace(row[2]),
		})
	}
	return a
}

func (t *TweetResponder) ReadRecentTweets(username string) []twitter.Tweet {
	excludeReplies := true
	params := twitter.UserTimelineParams{
		ScreenName:     username,
		ExcludeReplies: &excludeReplies,
		Count:          100,
	}
	tweets, _, err := t.Client.Timelines.UserTimeline(&params)
	if err != nil {
		fmt.Errorf(err.Error())
		return []twitter.Tweet{}
	}

	return tweets
}

func (t *TweetResponder) GetRandomTweet(username string) *twitter.Tweet {
	tweets := t.ReadRecentTweets(username)
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
