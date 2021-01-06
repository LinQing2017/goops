package export

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	"strings"
)

func ExportK8sCluster() []*ClusterInfo {

	services := db_tools.GetAllClusterByEnv(envType)
	cinfos := make([]*ClusterInfo, 0)
	for i, s := range services {
		clusterInfo := db_tools.GetK8SClusterInfo(s.ClusterId, db_tools.K8sDBlClient)
		if strings.EqualFold(clusterInfo.ID, "") {
			continue
		}
		var app types.App
		if err := db_tools.GetOne(db_tools.PortalMongoDB, "app", bson.M{"name": clusterInfo.AppName}, db_tools.NdpPortalClient, &app); err != nil {
			logrus.Error("没有查询到App信息")
			continue
		}

		cinfo := ClusterInfo{
			AppName:     clusterInfo.AppName,
			ClusterName: fmt.Sprintf("%s-%s", "k8s", clusterInfo.ID[len(clusterInfo.ID)-6:]),
			Replica:     clusterInfo.Replica,
			Creator:     fmt.Sprintf("%s(%d)", app.CreatorName, app.Creator),
			MaxReplica:  clusterInfo.MaxReplica,
			MemRequest:  clusterInfo.Config.ContainerTemplate.Resources.Requests.Memory["amount"],
			MemLimit:    clusterInfo.Config.ContainerTemplate.Resources.Limits.Memory["amount"],
		}

		if withUrl {
			cinfos[i].AppUrl = common.GetAppUrl(s.ClusterId, clusterInfo.AppName)
		}
		cinfos = append(cinfos, &cinfo)
	}
	return cinfos
}
