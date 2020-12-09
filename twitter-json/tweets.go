package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type Twitter struct {
	Profile twitterscraper.Profile
	Tweets  []twitterscraper.Tweet
}

func GetFileContent(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	//读取文件全部内容
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	var users []string
	twitterPath := "./raw/twitter/"
	f, err := GetFileContent(twitterPath + "userList.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		f = strings.Replace(f, " ", "", -1) // 删除空格
		users = strings.Split(f, ",")
	}

	for _, user := range users {
		var twitter Twitter
		var tweets []twitterscraper.Tweet
		maxTweetsNbr := 50
		profile, err := twitterscraper.GetProfile(user)
		if err != nil {
			fmt.Println(err)
		}
		twitter.Profile = profile
		count := 0
		for tweet := range twitterscraper.GetTweets(context.Background(), user, maxTweetsNbr) {
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
