package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"twitter-json/src/model"
	"twitter-json/src/utils"

	"github.com/caarlos0/env/v6"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

func checkImageExist(picBed string, coll *qmgo.Collection, user string, path string, name string) (exist bool, ret utils.DbImage, ghPath string) {
	pathReg := regexp.MustCompile(`(?s)twitter-json/raw(.*?)$`)
	ghPath = "." + pathReg.FindStringSubmatch(path)[1]
	ret = utils.DbImage{User: user}
	if picBed == "github" {
		_, err := os.Stat(path)
		if err == nil {
			exist = true
		} else {
			exist = false
		}
		ret.Key = ghPath
	} else {
		e, key := model.DbCheckImageExist(coll, fmt.Sprintf("%s/%s/%s", "resource/twitter", user, name))
		exist = e
		ret.Key = key
	}

	return exist, ret, ghPath
}

func uploadImage(picBed string, coll *qmgo.Collection, user string, path string, key string) (ret utils.DbImage, err error) {
	ret = utils.DbImage{Key: key, User: user}
	if picBed == "qiniu" {
		retQiniu, err := model.QiniuUpload(path, key)
		if err != nil {
			return ret, err
		}
		ret.Key = retQiniu.URL
		os.Remove(path)
	}
	dbImg := utils.DbImage{Key: key, User: user}
	model.DbInsertImage(coll, dbImg)
	return ret, nil
}

func imageProcess(coll *qmgo.Collection, picBed string, picDir string, user string, src string) (htmlSrc string) {
	localPath, fileName, err := utils.GetImageInfo(picDir+user, src)
	if err != nil {
		fmt.Println(err)
	}
	exit, ret, ghPath := checkImageExist(picBed, coll, user, localPath, fileName)
	if exit {
		if picBed == "github" {
			htmlSrc = ret.Key
		} else {
			if picBed == "qiniu" {
				htmlSrc = model.QiniuMakeURL(ret.Key)
			} else {
				htmlSrc = ret.Key
			}
		}
	} else {
		err := utils.WriteImage(localPath, src)
		if err != nil {
			fmt.Println(err)
		} else {
			if picBed == "github" {
				htmlSrc = ret.Key
				fmt.Printf("filePath %s\n", ret.Key)
			} else {
				retUpload, err := uploadImage(picBed, coll, user, localPath, ret.Key)
				if err != nil {
					// 上传失败用 github 资源
					htmlSrc = ghPath
				} else {
					htmlSrc = retUpload.Key
				}
			}
		}
	}
	return htmlSrc
}

func main() {
	cfg := utils.Settings{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)

	ctx := context.Background()
	cli := model.DbInit()
	collProfile := model.DbColl(cli, "profile")
	collTweet := model.DbColl(cli, "tweets")
	collImage := model.DbColl(cli, "images")
	defer func() {
		model.DbClose(cli)
	}()

	var users []string
	dir, _ := os.Getwd()
	twitterDir := path.Join(dir, "./raw/twitter/") + "/"
	picDir := twitterDir + "images/"
	// pathReg := regexp.MustCompile(`(?s)twitter-json/raw(.*?)$`)
	imgReg := regexp.MustCompile(`(?s)(\<br\>\<a href=(.*)\>\<img(.*)\<\/a\>)(?s)`)
	f, err := utils.GetFileContent(twitterDir + "userList.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		f = strings.Replace(f, " ", "", -1) // 删除空格
		users = strings.Split(f, ",")
	}

	for _, user := range users {
		scraper := twitterscraper.New()
		var twitter utils.Twitter
		var tweets []twitterscraper.Tweet
		maxTweetsNbr := 50
		count := 0
		one := utils.DbProfile{}
		err = collProfile.Find(ctx, bson.M{"username": user}).One(&one)
		lastTweetTime := one.LastTweetTime
		curLastTweetTime := lastTweetTime
		for tweet := range scraper.GetTweets(context.Background(), user, maxTweetsNbr) {
			if tweet.Error != nil {
				fmt.Println(tweet.Error)
			} else {
				if lastTweetTime >= tweet.Tweet.Timestamp {
					break
				}
				if count == 0 {
					curLastTweetTime = tweet.Tweet.Timestamp
				}
				if tweet.Tweet.Photos != nil && len(tweet.Tweet.Photos) > 0 {
					var photos []string
					for _, src := range tweet.Tweet.Photos {
						htmlSrc := imageProcess(collImage, cfg.PicBed, picDir, user, src)
						photos = append(photos, htmlSrc)
					}
					tweet.Tweet.Photos = photos
				}
				if tweet.Tweet.HTML != "" {
					tweet.Tweet.HTML = imgReg.ReplaceAllString(tweet.Tweet.HTML, "")
				}
				tweets = append(tweets, tweet.Tweet)
				count++
			}
		}
		lastTweetTime = curLastTweetTime
		// twitter.Tweets = tweets
		if count > 0 {
			collTweet.InsertMany(ctx, tweets)
		}
		// get user tweets count
		tweetsCnt, err := collTweet.Find(ctx, bson.M{"username": user}).Count()

		profile, err := scraper.GetProfile(user)
		if err != nil {
			fmt.Println(err)
			one.TweetsCount = tweetsCnt
			one.LastTweetTime = lastTweetTime
			one.LastUpdateTime = time.Now().Unix()
			twitter.DbProfile = one
		} else {
			if profile.Avatar != "" {
				profile.Avatar = imageProcess(collImage, cfg.PicBed, picDir, user, profile.Avatar)
			}
			dbProfile := utils.DbProfile{Username: user, Profile: profile, LastTweetTime: lastTweetTime, LastUpdateTime: time.Now().Unix()}
			dbProfile.TweetsCount = tweetsCnt
			twitter.DbProfile = dbProfile

			cnt, err := collProfile.Find(ctx, bson.M{"username": profile.Username}).Count()
			if err != nil {
				fmt.Println(err)
			}
			if cnt > 0 {
				update := bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "profile", Value: dbProfile.Profile},
						{Key: "tweetscount", Value: dbProfile.TweetsCount},
						{Key: "lasttweettime", Value: dbProfile.LastTweetTime},
						{Key: "lastupdatetime", Value: dbProfile.LastUpdateTime},
					}},
				}
				e := collProfile.UpdateOne(ctx, bson.M{"username": profile.Username}, update)
				if e != nil {
					fmt.Println(e)
				}
			} else {
				collProfile.InsertOne(ctx, dbProfile)
			}
		}

		pages := int64(math.Ceil(float64(tweetsCnt) / float64(cfg.PageSize)))
		for i := int64(0); i < pages; i++ {
			skipCnt := int64(cfg.PageSize * i)
			dbTweets := []twitterscraper.Tweet{}
			collTweet.Find(ctx, bson.M{"username": user}).Skip(skipCnt).Sort("-timestamp").Limit(cfg.PageSize).All(&dbTweets)
			twitter.Tweets = dbTweets

			jsonBytes, err := json.Marshal(twitter)
			if err != nil {
				fmt.Println(err)
			}
			fileName := fmt.Sprintf("%s%s%s%d%s", twitterDir, user, "-", i+1, ".json")
			file, er := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer func() { file.Close() }()
			if er != nil && os.IsNotExist(er) {
				file, err = os.Create(fileName)
			}
			_, err = file.Write(jsonBytes)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if !cfg.DbInit {
		s, err := utils.GetFileContent(cfg.SettingsPath)
		settings := utils.Settings{}
		err = json.Unmarshal([]byte(s), &settings)
		if err != nil {
			fmt.Println(err)
		} else {
			settings.DbInit = true
			jsonBytes, err := json.Marshal(settings)
			if err != nil {
				fmt.Println(err)
			}
			file, er := os.OpenFile(cfg.SettingsPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer func() { file.Close() }()
			if er != nil && os.IsNotExist(er) {
				file, err = os.Create(cfg.SettingsPath)
			}
			_, err = file.Write(jsonBytes)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
