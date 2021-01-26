package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"goops/pkg/appinfo/client/ews"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	"strings"
)

// 从数据库获取所有信息
func GetAppInfo(appname string, envType int) (AppInformation, error) {

	appname = strings.TrimSpace(appname)
	appinformation := AppInformation{
		NAME:               appname,
		ClusterBindDomains: make(map[string][]*types.RMDomains),
	}
	// 从Portal数据库获取相应信息
	portalInfo, err := db_tools.GetPortalInfo(appname, envType, db_tools.NdpPortalClient)
	if err != nil {
		return appinformation, err
	}
	appinformation.APPID = portalInfo.APP.ID.Hex()
	appinformation.URL = GetAppUrl(appinformation.APPID, appname)
	appinformation.PortalInfo = portalInfo
	// 通过弹性Web环境获取弹性web集群信息
	appinformation.EWSClusterInfo = make([]ews.EWSCluster, 0)
	for _, ewsServer := range portalInfo.EWSServiceList {
		if ewsCluster, err := ews.GetCluster(portalInfo.APP.Name, ewsServer.ClusterId); err == nil {
			appinformation.EWSClusterInfo = append(appinformation.EWSClusterInfo, ewsCluster)
		} else {
			logrus.Error("获取弹性Web集群信息失败：", portalInfo.APP.Name, ewsServer.ClusterId, err.Error())
			continue
		}
		// 获取绑定的域名信息
		if domains, err := GetClusterBindDomains(ewsServer.ClusterId); err == nil && domains != nil && len(domains) > 0 {
			appinformation.ClusterBindDomains[ewsServer.ClusterId] = domains
		}
	}

	// 从K8S数据库获取相关信息
	appinformation.K8SClusterInfo = make([]types.K8SCluster, 0)
	for _, k8sCluster := range portalInfo.K8SServiceList {
		k8sClusterInfo := db_tools.GetK8SClusterInfo(k8sCluster.ClusterId, db_tools.K8sDBlClient)
		domains, err := GetClusterBindDomains(k8sCluster.ClusterId)
		if err == nil && domains != nil && len(domains) > 0 {
			appinformation.ClusterBindDomains[k8sCluster.ClusterId] = domains
		}
		appinformation.K8SClusterInfo = append(appinformation.K8SClusterInfo, k8sClusterInfo)
	}

	return appinformation, nil
}

func GetClusterBindDomains(clusterId string) ([]*types.RMDomains, error) {
	var domainBind []*types.RMDomains
	var cmptDomainBind []*types.RMDomains
	if !strings.EqualFold(clusterId, "") {
		// 从domains表中查询域名信息
		if err := db_tools.GetBatch(db_tools.RMMongoDB, "domains",
			bson.M{"cluster_id": clusterId}, db_tools.RMClient, &domainBind); err != nil {
			logrus.Error("域名信息查询失败", err.Error())
			return nil, err
		}
		// 从component_domains表中查询域名信息
		if err := db_tools.GetBatch(db_tools.RMMongoDB, "component_domains",
			bson.M{"cluster_id": clusterId}, db_tools.RMClient, &cmptDomainBind); err != nil {
			logrus.Error("域名信息查询失败", err.Error())
			return nil, err
		}
	}
	allDomain := append(domainBind, cmptDomainBind...)

	return allDomain, nil
}

func GetAppUrl(appId, appName string) string {
	return fmt.Sprintf("https://da.sdp.101.com/#/ndpfront/applicationManagement/applicationList/serviceInformation/%s/%s", appId, appName)
}
