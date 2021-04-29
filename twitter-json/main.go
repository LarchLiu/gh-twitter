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

func checkImageExist(picBed string, coll *qmgo.Collection, user string, path string, fileName string) (exist bool, url string, key string, ghPath string) {
	pathReg := regexp.MustCompile(`(?s)twitter-json/raw(.*?)$`)
	ghPath = "." + pathReg.FindStringSubmatch(path)[1]
	switch picBed {
	case "qiniu":
		e, u := model.DbCheckImageExist(coll, fileName, user)
		exist = e
		if e {
			url = u
		} else {
			prefix := model.QiniuGetResourcePrefix()
			key = fmt.Sprintf("%s/%s/%s/%s", prefix, "images", user, fileName)
			// 如果 prefix 为空需删除开头的 /
			key = strings.TrimPrefix(key, "/")
		}
	case "smms":
		e, u := model.DbCheckImageExist(coll, fileName, user)
		exist = e
		if e {
			url = u
		} else {
			key = fmt.Sprintf("%s_%s", user, fileName)
		}
	default:
		_, err := os.Stat(path)
		if err == nil {
			exist = true
		} else {
			exist = false
		}
		url = ghPath
	}

	return exist, url, key, ghPath
}

func uploadImage(picBed string, coll *qmgo.Collection, path string, key string, dbImage utils.DbImage, reUpload bool) (string, error) {
	var url string = ""
	var e error = nil
	dbImage.PicBed = picBed
	switch picBed {
	case "qiniu":
		retQiniu, err := model.QiniuUpload(path, key)
		if err != nil {
			dbImage.Status = "fail"
			e = err
		} else {
			url = retQiniu.URL
			dbImage.Key = key
			dbImage.URL = url
			dbImage.Status = "success"
			os.Remove(path)
		}
	case "smms":
		retSmms, err := model.SmmsUpload(path, key)
		if err != nil {
			dbImage.Status = "fail"
			e = err
		} else {
			url = retSmms.URL
			dbImage.Key = retSmms.Hash
			dbImage.URL = url
			dbImage.Status = "success"
			os.Remove(path)
		}
	}
	fmt.Println(dbImage)
	if reUpload {
		model.DbUpdateImage(coll, dbImage)
	} else {
		model.DbInsertImage(coll, dbImage)
	}
	return url, e
}

func uploadJSON(picBed string, path string, key string) (url string, err error) {
	if picBed == "qiniu" {
		prefix := model.QiniuGetResourcePrefix()
		joinKey := fmt.Sprintf("%s/%s/%s", prefix, "json", key)
		// 如果 prefix 为空需删除开头的 /
		joinKey = strings.TrimPrefix(joinKey, "/")
		retQiniu, err := model.QiniuUpload(path, joinKey)
		if err != nil {
			return "", err
		}
		url = retQiniu.URL
		os.Remove(path)
	}
	return url, nil
}

func imageProcess(coll *qmgo.Collection, picBed string, picDir string, user string, src string, dbImage utils.DbImage, limiter *utils.RateLimiter) (htmlSrc string) {
	if picBed == "cfworkers" {
		domain := os.Getenv("CF_WORKERS_IMG_DOMAIN")
		return utils.ReplaceDomain(src, domain)
	}
	localPath, fileName, err := utils.GetImageInfo(picDir+user, src)
	if err != nil {
		fmt.Println(err)
	}
	dbImage.FileName = fileName
	exit, url, key, ghPath := checkImageExist(picBed, coll, user, localPath, fileName)
	if exit {
		htmlSrc = url
	} else {
		err := utils.WriteImage(localPath, src)
		if err != nil {
			fmt.Println(err)
		} else {
			if picBed == "github" {
				htmlSrc = ghPath
				fmt.Printf("filePath %s\n", ghPath)
			} else {
				if limiter != nil {
					limiter.Limit()
				}
				url, err := uploadImage(picBed, coll, localPath, key, dbImage, false)
				if err != nil {
					// 上传失败用 github 资源
					htmlSrc = ghPath
				} else {
					htmlSrc = url
				}
			}
		}
	}
	return htmlSrc
}

func failImageProcess(collImage *qmgo.Collection, collTweet *qmgo.Collection, collProfile *qmgo.Collection, imgDir string, user string, picBed string, limiter *utils.RateLimiter) bool {
	var imageChange bool = false
	dbImages := []utils.DbImage{}
	filter := bson.M{"$and": bson.A{bson.M{"status": "fail"}, bson.M{"user": user}}}

	collImage.Find(context.Background(), filter).All(&dbImages)
	if len(dbImages) > 0 {
		for _, img := range dbImages {
			filePath := path.Join(imgDir+user, img.FileName)
			exist := false
			_, err := os.Stat(filePath)
			if err == nil {
				exist = true
			} else {
				exist = false
			}
			if exist {
				key := ""
				switch picBed {
				case "qiniu":
					key = img.Key
				case "smms":
					key = fmt.Sprintf("%s_%s", user, img.FileName)
				}
				if limiter != nil {
					limiter.Limit()
				}
				url, err := uploadImage(picBed, collImage, filePath, key, img, true)
				if err != nil {
					// 上传失败
					imageChange = imageChange || false
				} else {
					if img.Type == "photo" {
						filter := bson.M{"$and": bson.A{bson.M{"id": img.TweetID}, bson.M{"username": user}}}
						one := twitterscraper.Tweet{}
						err := collTweet.Find(context.Background(), filter).One(&one)
						if err != nil {
							fmt.Println(err)
							imageChange = imageChange || false
						} else {
							one.Photos[img.Idx] = url
							update := bson.D{
								{Key: "$set", Value: bson.D{
									{Key: "photos", Value: one.Photos},
								}},
							}
							err := collTweet.UpdateOne(context.Background(), filter, update)
							if err == nil {
								imageChange = imageChange || true
							} else {
								imageChange = imageChange || false
							}
						}
					} else if img.Type == "video" {
						filter := bson.M{"$and": bson.A{bson.M{"id": img.TweetID}, bson.M{"username": user}}}
						one := twitterscraper.Tweet{}
						err := collTweet.Find(context.Background(), filter).One(&one)
						if err != nil {
							fmt.Println(err)
							imageChange = imageChange || false
						} else {
							one.Videos[img.Idx].Preview = url
							update := bson.D{
								{Key: "$set", Value: bson.D{
									{Key: "Videos", Value: one.Videos},
								}},
							}
							err := collTweet.UpdateOne(context.Background(), filter, update)
							if err == nil {
								imageChange = imageChange || true
							} else {
								imageChange = imageChange || false
							}
						}
					} else {
						filter := bson.M{"userinfo.username": user}
						one := utils.DbProfile{}
						err := collProfile.Find(context.Background(), filter).One(&one)
						if err != nil {
							fmt.Println(err)
							imageChange = imageChange || false
						} else {
							if img.Idx == 0 {
								one.Avatar = url
								one.Profile.Avatar = url
								update := bson.D{
									{Key: "$set", Value: bson.D{
										{Key: "userinfo", Value: one.UserInfo},
										{Key: "profile", Value: one.Profile},
									}},
								}
								err := collProfile.UpdateOne(context.Background(), filter, update)
								if err == nil {
									imageChange = imageChange || true
								} else {
									imageChange = imageChange || false
								}
							}
						}
					}
				}
			} else {
				imageChange = imageChange || false
			}
		}

	}
	return imageChange
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
			filePath := fmt.Sprintf("%s%s%d%s", jsonUserDir, "/", i+1, ".json")
			file, er := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer func() { file.Close() }()
			if er != nil && os.IsNotExist(er) {
				file, err = os.Create(filePath)
			}
			_, err = file.Write(jsonBytes)
			if err != nil {
				fmt.Println(err)
			}
			key := fmt.Sprintf("%s/%d%s", selUser, i+1, ".json")
			uploadJSON(picBed, filePath, key)
		}
	}
}

func main() {
	videoDomain := os.Getenv("CF_WORKERS_VIDEO_DOMAIN")
	info := utils.UserChangeInfo{}
	if err := env.Parse(&info); err != nil {
		fmt.Printf("%+v\n", err)
	}

	cfg := utils.Settings{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)

	var lim *utils.RateLimiter = nil
	if cfg.PicBed == "smms" {
		lim = utils.NewRateLimiter(time.Minute, 20)
	}
	ctx := context.Background()
	cli := model.DbInit()
	collProfile := model.DbColl(cli, "profile")
	collTweet := model.DbColl(cli, "tweets")
	collImage := model.DbColl(cli, "images")
	defer func() {
		model.DbClose(cli)
	}()

	if info.Type == "addusers" || info.Type == "delusers" {
		users := strings.Split(info.Users, ",")
		for _, user := range users {
			one := utils.DbProfile{}
			count, err := collProfile.Find(ctx, bson.M{"userinfo.username": user}).Count()
			if err != nil {
				fmt.Println(err)
			}
			if count > 0 {
				if info.Type == "delusers" {
					collProfile.RemoveAll(ctx, bson.M{"userinfo.username": user})
					collTweet.RemoveAll(ctx, bson.M{"username": user})
					collImage.RemoveAll(ctx, bson.M{"user": user})
				}
				continue
			}

			if info.Type == "addusers" {
				one.Username = user
				collProfile.InsertOne(ctx, utils.DbProfile{UserInfo: one.UserInfo})
			}
		}
	}

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
		imgChange := failImageProcess(collImage, collTweet, collProfile, picDir, user, cfg.PicBed, lim)
		var twitter utils.Twitter
		var tweets []twitterscraper.Tweet
		maxTweetsNbr := 200
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
						for i, src := range tweet.Tweet.Photos {
							dbImage := utils.DbImage{Type: "photo", TweetID: tweet.Tweet.ID, Idx: i, User: user}
							htmlSrc := imageProcess(collImage, cfg.PicBed, picDir, user, src, dbImage, lim)
							photos = append(photos, htmlSrc)
						}
						tweet.Tweet.Photos = photos
					}
					if tweet.Tweet.Videos != nil && len(tweet.Tweet.Videos) > 0 {
						var videos []twitterscraper.Video
						for i, video := range tweet.Tweet.Videos {
							dbImage := utils.DbImage{Type: "video", TweetID: tweet.Tweet.ID, Idx: i, User: user}
							htmlSrc := imageProcess(collImage, cfg.PicBed, picDir, user, video.Preview, dbImage, lim)
							video.Preview = htmlSrc
							video.URL = utils.ReplaceDomain(video.URL, videoDomain)
							videos = append(videos, video)
						}
						tweet.Tweet.Videos = videos
					}
					if tweet.Tweet.HTML != "" {
						tweet.Tweet.HTML = imgReg.ReplaceAllString(tweet.Tweet.HTML, "")
					}
					tweets = append(tweets, tweet.Tweet)
					count++
				}
			}

			if imgChange || count > 0 {
				isUpdate = true
				updateUsers = append(updateUsers, user)
				lastTweetTime = curLastTweetTime

				if count > 0 {
					collTweet.InsertMany(ctx, tweets)
				}

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
						dbImage := utils.DbImage{Type: "profile", Idx: 0, User: user}
						profile.Avatar = imageProcess(collImage, cfg.PicBed, picDir, user, profile.Avatar, dbImage, lim)
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

	if isUpdate || info.Type == "delusers" {
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
