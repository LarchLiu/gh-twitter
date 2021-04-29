package model

import (
	"context"
	"fmt"
	"os"
	"twitter-json/src/utils"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

// DbInit init mogodb.
func DbInit() *qmgo.Client {
	uri := os.Getenv("MONGODB_CLOUD_URI")
	client, _ := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: uri})
	return client
}

// DbColl use collection.
func DbColl(client *qmgo.Client, name string) *qmgo.Collection {
	db := client.Database("twitter")
	coll := db.Collection(name)
	return coll
}

// DbClose colse database.
func DbClose(cli *qmgo.Client) {
	if err := cli.Close(context.Background()); err != nil {
		panic(err)
	}
}

// DbCheckImageExist check file exsit.
func DbCheckImageExist(coll *qmgo.Collection, fileName string, user string) (exist bool, out string) {
	one := utils.DbImage{}
	filter := bson.M{"$and": bson.A{bson.M{"filename": fileName}, bson.M{"user": user}}}
	err := coll.Find(context.Background(), filter).One(&one)
	if err != nil {
		fmt.Println(err)
	}
	if one.Status == "success" {
		return true, one.URL
	}
	return false, ""
}

// DbInsertImage insert image info to db.
func DbInsertImage(coll *qmgo.Collection, img utils.DbImage) {
	coll.InsertOne(context.Background(), img)
}

// DbUpdateImage update image info to db
func DbUpdateImage(coll *qmgo.Collection, img utils.DbImage) {
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "key", Value: img.Key},
			{Key: "url", Value: img.URL},
			{Key: "status", Value: img.Status},
		}},
	}
	coll.UpdateOne(context.Background(), bson.M{"$and": bson.A{bson.M{"filename": img.FileName}, bson.M{"user": img.User}}}, update)
}
