package main

import (
	//"golang.org/x/net/context"
	"golang.org/x/oauth2"
	//	"golang.org/x/oauth2/google"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	accessToken := ReadSecret("config/accessToken.txt")
	MakeToken(accessToken)

}

func ReadSecret(filename string) string {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	key := strings.TrimSpace(string(text))
	return key
}

const twitter = "https://twitter.com/"

func GetTweet(username string) string {
	//html := ReadWebpage(twitter + username)
	json := ReadRecentTweets(username)

	re := regexp.MustCompile(`https://video.twimg.com/.*\.mp4`)
	//re := regexp.MustCompile(`https://pbs.twimg.com/.*\.jpg`)
	pics := re.FindAllString(json, -1)

	numPics := len(pics)
	if numPics == 0 {
		return json
	}
	magicNumber := rand.Intn(numPics)
	fmt.Println(len(pics))
	return pics[magicNumber]
}

func MakeToken(accessToken string) *oauth2.Token {
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: accessToken}
	// OAuth2 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext, token)

	username := "tokugifs"
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=%v", username)

	resp, err := httpClient.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(resp.Body)

	return token

}

func handler(w http.ResponseWriter, r *http.Request, apiKey string) {
	//appendtoken := MakeToken("todo")
	//var tokenSrc oauth2.TokenSource
	// okay to get if nil
	//tokenSrc := oauth2.ReuseTokenSource(token, tokenSrc)

	//client := oauth2.NewClient()
	//client.Get("...")
}

func ReadRecentTweets(username string) string {

	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=%v", username)
	resp, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func ReadWebpage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	html := string(bytes)
	fmt.Println("HTML:\n\n", html)

	resp.Body.Close()

	return html
}
