package model

import (
	"context"
	"testing"
	"twitter-json/src/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func TestDB(t *testing.T) {
	cli := DbInit()
	// collProfile := DbColl(cli, "profile")
	// collTweet := DbColl(cli, "tweets")
	collImage := DbColl(cli, "images")
	dbImage := utils.DbImage{Type: "tweet", TweetID: "123", Idx: 0, User: "FCBarcelona", Status: "success",
		FileName: "Ep_NssfXYAA6KwJ.jpg", Key: "resource/twitter/FCBarcelona/Ep_NssfXYAA6KwJ.jpg"}
	collImage.InsertOne(context.Background(), dbImage)
	b, key := DbCheckImageExist(collImage, "Ep_NssfXYAA6KwJ.jpg", "FCBarcelona")
	collImage.RemoveAll(context.Background(), bson.M{"tweetid": "123"})
	if !b || key != "resource/twitter/FCBarcelona/Ep_NssfXYAA6KwJ.jpg" {
		t.Error("res ", b, key)
	}
}
