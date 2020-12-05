package db_tools

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/db_tools/types"
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

func GetPortalInfo(appname string, envType int, client *mongo.Client) types.AppPortalInfo {
	appPortalInfo := types.AppPortalInfo{}

	// 获取APP表数据
	var app types.App
	if err := GetOne(PortalMongoDB, "app", bson.M{"name": appname}, client, &app); err != nil {
		logrus.Error("没有查询到App信息")
	}
	appPortalInfo.APP = app

	//获取Environment表记录
	var environments []*types.Environment
	var filter interface{}
	if envType == 0 {
		filter = bson.M{"appId": app.ID.Hex()}
	} else {
		filter = bson.M{"appId": app.ID.Hex(), "envType": envType}
	}
	if err := GetBatch(PortalMongoDB, "environment", filter, client, &environments); err != nil {
		logrus.Error("没有查询到environment信息")
	}
	appPortalInfo.Environments = environments

	// 获取Service信息
	var ewsServiceList []*types.Service
	if envType == 0 {
		filter = bson.M{"name": app.Name, "type": 1}
	} else {
		filter = bson.M{"name": app.Name, "type": 1, "env": envType}
	}
	if err := GetBatch(PortalMongoDB, "service", filter, client, &ewsServiceList); err != nil {
		logrus.Info("没有查询到弹性web集群")
	}
	appPortalInfo.EWSServiceList = ewsServiceList

	var k8sServiceList []*types.Service
	if envType == 0 {
		filter = bson.M{"name": app.Name, "type": 29}
	} else {
		filter = bson.M{"name": app.Name, "type": 29, "env": envType}
	}
	if err := GetBatch(PortalMongoDB, "service", filter, client, &k8sServiceList); err != nil {
		logrus.Info("没有查询到K8s集群")
	}
	appPortalInfo.K8SServiceList = k8sServiceList
	return appPortalInfo
}
