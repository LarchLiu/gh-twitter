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
	if picBed == "qiniu" {
		e, key := model.DbCheckImageExist(coll, fmt.Sprintf("%s/%s/%s/%s", os.Getenv("QINIU_RESOURCE_PREFIX"), "images", user, name))
		exist = e
		ret.Key = key
	} else {
		_, err := os.Stat(path)
		if err == nil {
			exist = true
		} else {
			exist = false
		}
		ret.Key = ghPath
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

func uploadJSON(picBed string, path string, key string) (url string, err error) {
	if picBed == "qiniu" {
		key := fmt.Sprintf("%s/%s/%s", os.Getenv("QINIU_RESOURCE_PREFIX"), "json", key)
		retQiniu, err := model.QiniuUpload(path, key)
		if err != nil {
			return "", err
		}
		url = retQiniu.URL
		os.Remove(path)
	}
	return url, nil
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

func jsonTwitterFromDB(coll *qmgo.Collection, picBed string, dir string, selUser string, count int64, pageSize int64, twitter utils.Twitter) {
	jsonUserDir := dir + selUser
	exist := utils.CheckDirExist(jsonUserDir)
	if exist {
		pages := int64(math.Ceil(float64(count) / float64(pageSize)))
		for i := int64(0); i < pages; i++ {
			skipCnt := int64(pageSize * i)
			dbTweets := []twitterscraper.Tweet{}
			sel := bson.M{}
			if selUser != "@all" {
				sel = bson.M{"username": selUser}
			}
			coll.Find(context.Background(), sel).Skip(skipCnt).Sort("-timestamp").Limit(pageSize).All(&dbTweets)

			twitter.Tweets = dbTweets

			jsonBytes, err := json.Marshal(twitter)
			if err != nil {
				fmt.Println(err)
			}
			fileName := fmt.Sprintf("%s%s%d%s", jsonUserDir, "/", i+1, ".json")
			file, er := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer func() { file.Close() }()
			if er != nil && os.IsNotExist(er) {
				file, err = os.Create(fileName)
			}
			_, err = file.Write(jsonBytes)
			if err != nil {
				fmt.Println(err)
			}
			key := fmt.Sprintf("%s/%d%s", selUser, i+1, ".json")
			uploadJSON(picBed, fileName, key)
		}
	}
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
	var userInfoList []utils.UserInfo
	dir, _ := os.Getwd()
	twitterDir := path.Join(dir, "./raw/twitter") + "/"
	picDir := twitterDir + "images/"
	jsonDir := twitterDir + "json/"
	// pathReg := regexp.MustCompile(`(?s)twitter-json/raw(.*?)$`)
	imgReg := regexp.MustCompile(`(?s)(\<br\>\<a href=(.*)\>\<img(.*)\<\/a\>)(?s)`)
	if cfg.DbInit {
		all := []utils.DbProfile{}
		err := collProfile.Find(ctx, bson.D{{}}).Sort("_id").All(&all)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, user := range all {
				users = append(users, user.Username)
			}
		}
	} else {
		f, err := utils.GetFileContent(twitterDir + "userList.txt")
		if err != nil {
			fmt.Println(err)
		} else {
			f = strings.Replace(f, " ", "", -1) // 删除空格
			users = strings.Split(f, ",")
		}
	}

	scraper := twitterscraper.New()
	isUpdate := false
	updateUsers := []string{}
	for _, user := range users {
		var twitter utils.Twitter
		var tweets []twitterscraper.Tweet
		maxTweetsNbr := 50
		count := 0
		one := utils.DbProfile{}
		err := collProfile.Find(ctx, bson.M{"userinfo.username": user}).One(&one)
		if err != nil {
			fmt.Println(err)
		}
		if err == nil || !cfg.DbInit {
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

			if count > 0 {
				isUpdate = true
				updateUsers = append(updateUsers, user)
				lastTweetTime = curLastTweetTime

				collTweet.InsertMany(ctx, tweets)

				// get user tweets count
				tweetsCnt, err := collTweet.Find(ctx, bson.M{"username": user}).Count()
				pages := int(math.Ceil(float64(tweetsCnt) / float64(cfg.PageSize)))

				profile, err := scraper.GetProfile(user)
				if err != nil {
					fmt.Println(err)
					one.Username = user
					one.TweetsCount = tweetsCnt
					one.LastTweetTime = lastTweetTime
					one.LastUpdateTime = time.Now().Unix()
					one.Pages = pages
					twitter.DbProfile = one
					userInfoList = append(userInfoList, one.UserInfo)
					if !cfg.DbInit {
						collProfile.InsertOne(ctx, utils.DbProfile{UserInfo: one.UserInfo})
					}
				} else {
					if profile.Avatar != "" {
						profile.Avatar = imageProcess(collImage, cfg.PicBed, picDir, user, profile.Avatar)
					}
					userInfo := utils.UserInfo{Avatar: profile.Avatar, Username: user, Name: profile.Name,
						LastTweetTime: lastTweetTime, LastUpdateTime: time.Now().Unix(), TweetsCount: tweetsCnt, Pages: pages}
					dbProfile := utils.DbProfile{UserInfo: userInfo, Profile: profile}
					twitter.DbProfile = dbProfile
					userInfoList = append(userInfoList, userInfo)

					cnt, err := collProfile.Find(ctx, bson.M{"userinfo.username": profile.Username}).Count()
					if err != nil {
						fmt.Println(err)
					}
					if cnt > 0 {
						update := bson.D{
							{Key: "$set", Value: bson.D{
								{Key: "profile", Value: dbProfile.Profile},
								{Key: "userinfo", Value: dbProfile.UserInfo},
							}},
						}
						e := collProfile.UpdateOne(ctx, bson.M{"userinfo.username": profile.Username}, update)
						if e != nil {
							fmt.Println(e)
						}
					} else {
						collProfile.InsertOne(ctx, dbProfile)
					}
				}

				jsonTwitterFromDB(collTweet, cfg.PicBed, jsonDir, user, tweetsCnt, cfg.PageSize, twitter)
			} else {
				err := collProfile.Find(ctx, bson.M{"userinfo.username": user}).One(&one)
				if err != nil {
					fmt.Println(err)
				} else {
					tweetsCnt, err := collTweet.Find(ctx, bson.M{"username": user}).Count()
					if err != nil {
						tweetsCnt = 0
					}
					pages := int(math.Ceil(float64(tweetsCnt) / float64(cfg.PageSize)))

					one.Username = user
					one.TweetsCount = tweetsCnt
					one.LastTweetTime = lastTweetTime
					one.LastUpdateTime = time.Now().Unix()
					one.Pages = pages
					userInfoList = append(userInfoList, one.UserInfo)
					if !cfg.DbInit {
						collProfile.InsertOne(ctx, utils.DbProfile{UserInfo: one.UserInfo})
					}
				}
			}
		}
	}

	if isUpdate {
		name := "@all"
		twitter := utils.Twitter{}
		count, err := collTweet.Find(ctx, bson.M{}).Count()
		if err != nil {
			fmt.Println(err)
		} else {
			pages := int(math.Ceil(float64(count) / float64(cfg.PageSize)))
			one := twitterscraper.Tweet{}
			err = collTweet.Find(ctx, bson.M{}).Sort("-timestamp").Limit(1).One(&one)
			userInfo := utils.UserInfo{Username: name, Name: "All", LastTweetTime: one.Timestamp, LastUpdateTime: time.Now().Unix(), TweetsCount: count, Pages: pages}
			twitter.DbProfile = utils.DbProfile{UserInfo: userInfo}
			userInfoList = append([]utils.UserInfo{userInfo}, userInfoList...)

			jsonTwitterFromDB(collTweet, cfg.PicBed, jsonDir, name, count, cfg.PageSize, twitter)
			updateUsers = append([]string{name}, updateUsers...)
		}

		fileName := "userList.json"
		jsonBytes, err := json.Marshal(userInfoList)
		if err != nil {
			fmt.Println(err)
		}
		utils.WriteJSONFile(jsonDir, fileName, jsonBytes)
		uploadJSON(cfg.PicBed, jsonDir+fileName, fileName)
	}

	// write and upload json file of updataInfo
	updateInfo := utils.UpdateInfo{UpdateTime: time.Now().Unix(), IsUpdate: isUpdate, Users: updateUsers}
	fileName := "updateInfo.json"
	jsonBytes, err := json.Marshal(updateInfo)
	if err != nil {
		fmt.Println(err)
	}
	utils.WriteJSONFile(jsonDir, fileName, jsonBytes)
	uploadJSON(cfg.PicBed, jsonDir+fileName, fileName)

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
			utils.WriteJSONFileByPath(cfg.SettingsPath, jsonBytes)
		}
	}
}
