package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
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

func randFileName(dirPath string, suffix string) (fileName string) {
	var randChar string = "abcdefghimnjkqzxyt0123456789ioedkaldncalew0129387iue"
	var name string
	size := len(randChar)
	for i := 0; i < 12; i++ {
		name += string(randChar[rand.Intn(size)])
	}
	fileName = path.Join(dirPath, name)
	fileName += suffix
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fileName
		} else {
			return randFileName(dirPath, suffix)
		}
	}
	return randFileName(dirPath, suffix)
}

func writeImage(dirPath string, href string) string {
	i := strings.LastIndex(href, ".")
	suffix := href[i:]
	fileName := randFileName(dirPath, suffix)
	imageFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[writeImage create file]: fileName: %s\n href: %s\nerror: %s\n", fileName, href, err.Error())
		return ""
	}
	resp, err := http.Get(href)
	if err != nil {
		fmt.Printf("[writeImage http.Get: fileName]: %s\nhref: %s\n error: %s\n", fileName, href, err.Error())
		os.Remove(fileName)
		return ""
	}
	size, err := io.Copy(imageFile, resp.Body)
	if err != nil {
		fmt.Printf("io.Copy:error: %s  href: %s\n", err.Error(), href)
		os.Remove(fileName)
	}
	fmt.Printf("Get From %s: %d bytes\n", href, size)
	return fileName
}

func getCorrectHref(href string, url string) (correctHref string, exist bool) {
	completeValidator := regexp.MustCompile(`\A((http)|(https))(.+)((.jpg)|(.png)|(.jpeg))\z`)
	imcompleteValidator := regexp.MustCompile(`\A(||)(.+)((.jpg)|(.jpeg)|(.png))\z`)
	if completeValidator.MatchString(href) {
		//fmt.Println(href)
		return href, true
	} else if imcompleteValidator.MatchString(href) {
		ret := href[1:]
		ret = url + ret
		//fmt.Println(ret)
		return string(ret), true
	}
	return "", false
}

func main() {
	var users []string
	twitterPath := "./raw/twitter/"
	picPath := twitterPath + "images"
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
			if tweet.Tweet.Photos != nil {
				var photos []string
				for _, src := range tweet.Tweet.Photos {
					correctHref, exist := getCorrectHref(src, "")
					if exist {
						photos = append(photos, writeImage(picPath, correctHref))
					}
				}
				tweet.Tweet.Photos = photos
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
