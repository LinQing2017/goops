package list

import (
	"encoding/json"
	"github.com/modood/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	"goops/pkg/appinfo/ews_client"
	mongotools "goops/pkg/util/mongo"
	systools "goops/pkg/util/sys"
	"strings"
)

var (
	ndpPortalClient *mongo.Client
	k8sDBlClient    *mongo.Client
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "批量获取应用信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	ews_client.AddEWSFlag(cmd.Flags())
	AddListFlag(cmd.Flags())
	return cmd
}

func Main(cmd *cobra.Command, args []string) {

	if strings.EqualFold(nameFromFile, "") || envType == 0 {
		logrus.Error("请指定环境变量和输入文件")
		return
	}

	portalDBURI := "mongodb://" + db_tools.PortalMongoUser + ":" + db_tools.PortalMongoPasswd + "@" + db_tools.PortalMongoUrl + "/" + db_tools.PortalMongoDB + "?autoConnectRetry=true"
	ndpPortalClient = mongotools.MongoClient(portalDBURI)

	k8sDBURI := "mongodb://" + db_tools.K8SPaasMongoUser + ":" + db_tools.K8SPaasMongoPasswd + "@" + db_tools.K8SPaasMongoUrl + "/" + db_tools.K8SPaasMongoDB + "?autoConnectRetry=true"
	k8sDBlClient = mongotools.MongoClient(k8sDBURI)

	infoList := make([]ListAppInfo, 0)
	for _, appname := range systools.ReadLine(nameFromFile) {
		appname = strings.TrimSpace(appname)
		portalInfo := db_tools.GetPortalInfo(appname, envType, ndpPortalClient)
		infoPrint := getListAppInfo(portalInfo)
		infoList = append(infoList, infoPrint)
	}
	table.Output(infoList)
	mongotools.MongoDisconnect(ndpPortalClient)
	mongotools.MongoDisconnect(k8sDBlClient)
}

func getListAppInfo(portalInfo types.AppPortalInfo) ListAppInfo {

	listAppInfo := ListAppInfo{
		NAME:       portalInfo.APP.Name,
		APPID:      portalInfo.APP.ID.Hex(),
		Single:     portalInfo.APP.SingleInstance,
		EWSPackage: false,
		K8SNum:     len(portalInfo.K8SServiceList),
		EWSNum:     len(portalInfo.EWSServiceList),
		URL:        "",
	}

	if withUrl {
		listAppInfo.URL = "https://da.sdp.101.com/#/ndpfront/applicationManagement/applicationList/serviceInformation/" + portalInfo.APP.ID.Hex() + "/" + portalInfo.APP.Name
	}
	packageNum := 0 // 能够找到包路径的弹性web集群数目
	for _, ewsServer := range portalInfo.EWSServiceList {
		if ewsCluster, err := ews_client.GetCluster(portalInfo.APP.Name, ewsServer.ClusterId); err == nil {
			if len(ewsCluster.Instances) > 0 && !strings.EqualFold(ewsCluster.Instances[0].PackageUrl, "") {
				packageNum++
			}
		} else {
			logrus.Error("获取弹性Web集群信息失败：", portalInfo.APP.Name, ewsServer.ClusterId)
			continue
		}
	}

	// 如果有一个集群找不到包地址，那么判断该应用的包地址为空
	listAppInfo.EWSPackage = packageNum == listAppInfo.EWSNum

	// 打印K8S信息
	GetK8SClusterInfo(&listAppInfo, portalInfo.K8SServiceList)

	return listAppInfo
}

func GetK8SClusterInfo(listAppInfo *ListAppInfo, k8sClusterList []*types.Service) {
	// 获取K8S集群的Node-Select信息
	listAppInfo.NodeSelector = make([]string, 0)
	listAppInfo.NodeType = make([]string, 0)
	listAppInfo.NodeName = make([]string, 0)
	listAppInfo.K8SAREA = make([]string, 0)
	for _, k8sCluster := range k8sClusterList {
		k8sClusterInfo := db_tools.GetK8SClusterInfo(k8sCluster.ClusterId, k8sDBlClient)
		var nodeSelectorStr []byte
		if k8sClusterInfo.Config.NodeSelector != nil {
			if err := json.Unmarshal(nodeSelectorStr, &k8sClusterInfo.Config.NodeSelector); err != nil {
				logrus.Error("解析nodeSelector错误")
			}
		}

		if !strings.EqualFold(string(nodeSelectorStr), "") {
			listAppInfo.NodeSelector = append(listAppInfo.NodeSelector, string(nodeSelectorStr))
		}
		if !strings.EqualFold(k8sClusterInfo.Config.NodeType, "") {
			listAppInfo.NodeType = append(listAppInfo.NodeType, k8sClusterInfo.Config.NodeType)
		}
		if !strings.EqualFold(k8sClusterInfo.Config.NodeName, "") {
			listAppInfo.NodeName = append(listAppInfo.NodeName, k8sClusterInfo.Config.NodeName)
		}
		if !strings.EqualFold(k8sClusterInfo.Area, "") {
			listAppInfo.K8SAREA = append(listAppInfo.K8SAREA, k8sClusterInfo.Area)
		}
	}

}
