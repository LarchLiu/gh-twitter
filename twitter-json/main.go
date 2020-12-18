package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"twitter-json/src/utils"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
	var users []string
	dir, _ := os.Getwd()
	twitterPath := path.Join(dir, "./raw/twitter/") + "/"
	picPath := twitterPath + "images"
	imgReg := regexp.MustCompile(`(?s)(\<br\>\<a href=(.*)\>\<img(.*)\<\/a\>)(?s)`)
	pathReg := regexp.MustCompile(`(?s)twitter-json/raw(.*?)$`)
	f, err := utils.GetFileContent(twitterPath + "userList.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		f = strings.Replace(f, " ", "", -1) // 删除空格
		users = strings.Split(f, ",")
	}

	for _, user := range users {
		var twitter utils.Twitter
		var tweets []twitterscraper.Tweet
		maxTweetsNbr := 50
		profile, err := twitterscraper.GetProfile(user)
		if err != nil {
			fmt.Println(err)
		}

		if profile.Avatar != "" {
			filePath := utils.WriteImage(picPath, profile.Avatar)
			filePath = "." + pathReg.FindStringSubmatch(filePath)[1]
			fmt.Print(filePath)
			profile.Avatar = filePath
		}
		twitter.Profile = profile
		count := 0
		for tweet := range twitterscraper.GetTweets(context.Background(), user, maxTweetsNbr) {
			if tweet.Tweet.Photos != nil && len(tweet.Tweet.Photos) > 0 {
				var photos []string
				for _, src := range tweet.Tweet.Photos {
					filePath := utils.WriteImage(picPath, src)
					filePath = "." + pathReg.FindStringSubmatch(filePath)[1]
					fmt.Print(filePath)
					photos = append(photos, filePath)
				}
				tweet.Tweet.Photos = photos
			}
			if tweet.Tweet.HTML != "" {
				tweet.Tweet.HTML = imgReg.ReplaceAllString(tweet.Tweet.HTML, "")
			}
			tweets = append(tweets, tweet.Tweet)
			count++
		}
		twitter.Tweets = tweets
		// fmt.Println(twitter)

		jsonBytes, err := json.Marshal(twitter)
		if err != nil {
			fmt.Println(err)
		}
		file, er := os.OpenFile(twitterPath+user+".json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		defer func() { file.Close() }()
		if er != nil && os.IsNotExist(er) {
			file, err = os.Create(twitterPath + user + ".json")
		}
		_, err = file.Write(jsonBytes)
		if err != nil {
			fmt.Println(err)
		}
	}
}
