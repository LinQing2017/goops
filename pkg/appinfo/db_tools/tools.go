package db_tools

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	mongotools "goops/pkg/util/mongo"
)

var (
	NdpPortalClient *mongo.Client
	K8sDBlClient    *mongo.Client
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

func InitDBClient() {
	portalDBURI := "mongodb://" + PortalMongoUser + ":" + PortalMongoPasswd + "@" + PortalMongoUrl + "/" + PortalMongoDB + "?autoConnectRetry=true"
	NdpPortalClient = mongotools.MongoClient(portalDBURI)

	k8sDBURI := "mongodb://" + K8SPaasMongoUser + ":" + K8SPaasMongoPasswd + "@" + K8SPaasMongoUrl + "/" + K8SPaasMongoDB + "?autoConnectRetry=true"
	K8sDBlClient = mongotools.MongoClient(k8sDBURI)
}

func CloseAllDBClient() {

	mongotools.MongoDisconnect(NdpPortalClient)
	mongotools.MongoDisconnect(K8sDBlClient)
}
