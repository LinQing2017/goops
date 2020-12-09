package common

import (
	"github.com/sirupsen/logrus"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	"goops/pkg/appinfo/ews_client"
	"strings"
)

// 从数据库获取所有信息
func GetAppInfo(appname string, envType int) (AppInformation, error) {

	appname = strings.TrimSpace(appname)
	appinformation := AppInformation{
		NAME: appname,
	}
	// 从Portal数据库获取相应信息
	portalInfo, err := db_tools.GetPortalInfo(appname, envType, db_tools.NdpPortalClient)
	if err != nil {
		return appinformation, err
	}
	appinformation.APPID = portalInfo.APP.ID.Hex()
	appinformation.URL = "https://da.sdp.101.com/#/ndpfront/applicationManagement/applicationList/serviceInformation/" + portalInfo.APP.ID.Hex() + "/" + portalInfo.APP.Name
	appinformation.PortalInfo = portalInfo
	// 通过弹性Web环境获取弹性web集群信息
	appinformation.EWSClusterInfo = make([]ews_client.EWSCluster, 0)
	for _, ewsServer := range portalInfo.EWSServiceList {
		if ewsCluster, err := ews_client.GetCluster(portalInfo.APP.Name, ewsServer.ClusterId); err == nil {
			appinformation.EWSClusterInfo = append(appinformation.EWSClusterInfo, ewsCluster)
		} else {
			logrus.Error("获取弹性Web集群信息失败：", portalInfo.APP.Name, ewsServer.ClusterId)
			continue
		}
	}

	// 从K8S数据库获取相关信息
	appinformation.K8SClusterInfo = make([]types.K8SCluster, 0)
	for _, k8sCluster := range portalInfo.K8SServiceList {
		k8sClusterInfo := db_tools.GetK8SClusterInfo(k8sCluster.ClusterId, db_tools.K8sDBlClient)
		appinformation.K8SClusterInfo = append(appinformation.K8SClusterInfo, k8sClusterInfo)
	}

	return appinformation, nil
}
