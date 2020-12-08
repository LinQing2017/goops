package db_tools

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOne(mongoDB, collectionName string, filter interface{}, client *mongo.Client, doc interface{}) error {
	// 指定获取要操作的数据集
	collection := client.Database(mongoDB).Collection(collectionName)
	appInfo := collection.FindOne(context.TODO(), filter)
	if err := appInfo.Decode(doc); err != nil {
		return err
	}
	return nil
}

func GetBatch(mongoDB, collectionName string, filter interface{}, client *mongo.Client, docs interface{}) (err error) {
	var cursor *mongo.Cursor
	collection := client.Database(mongoDB).Collection(collectionName)
	if cursor, err = collection.Find(context.TODO(), filter); err != nil {
		logrus.Error("没有查询到Environment信息")
		return err
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), docs)
	return nil
}
