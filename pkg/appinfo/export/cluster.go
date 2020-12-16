package export

import (
	"fmt"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
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

		cinfo := ClusterInfo{
			AppName:     clusterInfo.AppName,
			MemRequest:  clusterInfo.Config.ContainerTemplate.Resources.Requests.Memory["amount"],
			MemLimit:    clusterInfo.Config.ContainerTemplate.Resources.Limits.Memory["amount"],
			ClusterName: fmt.Sprintf("%s-%s", "k8s", clusterInfo.ID[len(clusterInfo.ID)-6:]),
		}

		if withUrl {
			cinfos[i].AppUrl = common.GetAppUrl(s.ClusterId, clusterInfo.AppName)
		}
		cinfos = append(cinfos, &cinfo)
	}
	return cinfos
}
