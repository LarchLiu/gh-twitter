package model

import (
	"context"
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
func DbCheckImageExist(coll *qmgo.Collection, key string) (exist bool, out string) {
	cnt, _ := coll.Find(context.Background(), bson.M{"key": key}).Count()
	if cnt > 0 {
		return true, key
	}
	return false, key
}

// DbInsertImage insert image info to db.
func DbInsertImage(coll *qmgo.Collection, img utils.DbImage) {
	coll.InsertOne(context.Background(), img)
}
