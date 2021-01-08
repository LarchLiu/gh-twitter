package model

import (
	"context"
	"fmt"
	"twitter-json/src/utils"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

// DbInit init mogodb.
func DbInit() *qmgo.Client {
	client, _ := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: "mongodb://localhost:27017"})
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
