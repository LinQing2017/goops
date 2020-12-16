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

// 返回当前环境中所有应用的名称
func GetAllClusterByEnv(envId int) []*types.Service {

	var k8sServiceList []*types.Service
	filter := bson.M{"type": 29, "env": envId}
	if err := GetBatch(PortalMongoDB, "service", filter, NdpPortalClient, &k8sServiceList); err != nil {
		logrus.Info("没有查询到K8s集群")
	}
	return k8sServiceList
}
