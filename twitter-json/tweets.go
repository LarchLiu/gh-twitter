package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type Twitter struct {
	Profile twitterscraper.Profile
	Tweets  []twitterscraper.Tweet
}

func main() {
	var twitter Twitter
	var tweets []twitterscraper.Tweet
	maxTweetsNbr := 50
	profile, err := twitterscraper.GetProfile("FCBarcelona")
	if err != nil {
		fmt.Println(err)
	}
	twitter.Profile = profile
	count := 0
	for tweet := range twitterscraper.GetTweets(context.Background(), "FCBarcelona", maxTweetsNbr) {
		tweets = append(tweets, tweet.Tweet)
		count++
	}
	twitter.Tweets = tweets
	// fmt.Println(twitter)

	jsonBytes, err := json.Marshal(twitter)
	if err != nil {
		fmt.Println(err)
	}
	file, er := os.Open("./raw/twitter/FCBarcelona.json")
	defer func() { file.Close() }()
	if er != nil && os.IsNotExist(er) {
		file, err = os.Create("./raw/twitter/FCBarcelona.json")
	}
	_, err = file.Write(jsonBytes)
	if err != nil {
		fmt.Println(err)
	}
}
