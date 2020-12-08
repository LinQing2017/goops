package db_tools

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/db_tools/types"
)

func GetK8SClusterInfo(clusterid string, client *mongo.Client) types.K8SCluster {

	var k8sCluster types.K8SCluster
	if err := GetOne(K8SPaasMongoDB, "cluster", bson.M{"_id": clusterid}, client, &k8sCluster); err != nil {
		logrus.Error("没有查询到集群信息", err.Error())
	}
	return k8sCluster
}
