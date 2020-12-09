package list

import (
	"github.com/modood/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/ews_client"
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

	db_tools.InitDBClient()
	allInformation := make([]common.AppInformation, 0)
	for _, appname := range systools.ReadLine(nameFromFile) {
		appInformation, _ := common.GetAppInfo(appname, envType)
		allInformation = append(allInformation, appInformation)
	}
	db_tools.CloseAllDBClient()

	switch outputFormate {
	case "k8s":
		printK8s(allInformation)
	case "ews":
		printEws(allInformation)
	default:
		printDefault(allInformation)
	}
}

func printDefault(allInfo []common.AppInformation) {

	printList := make([]PrintAppInfo, len(allInfo))

	for i, info := range allInfo {
		printList[i] = PrintAppInfo{
			NAME:   info.NAME,
			APPID:  info.APPID,
			Single: "",
			K8SNum: len(info.PortalInfo.K8SServiceList),
			EWSNum: len(info.PortalInfo.EWSServiceList),
			URL:    info.URL,
		}
		if info.PortalInfo.APP.SingleInstance {
			printList[i].Single = "Y"
		}
	}
	table.Output(printList)
}

func printEws(allInfo []common.AppInformation) {

	printList := make([]PrintEWSInfo, len(allInfo))

	for i, info := range allInfo {
		printList[i] = PrintEWSInfo{
			NAME: info.NAME,
			Num:  len(info.PortalInfo.EWSServiceList),
		}
		packageNum := 0 // 能够找到包路径的弹性web集群数目
		for _, ewsCluster := range info.EWSClusterInfo {
			if len(ewsCluster.Instances) > 0 && !strings.EqualFold(ewsCluster.Instances[0].PackageUrl, "") {
				packageNum++
			}
		}
		// 所有集群都能找到一个对应实例，有包地址。否则认为应用缺少war包
		if packageNum < len(info.PortalInfo.EWSServiceList) {
			printList[i].PackageURL = "Not Found"
		}
	}
	table.Output(printList)
}

func printK8s(allInfo []common.AppInformation) {

	printList := make([]PrintK8SInfo, len(allInfo))
	for i, info := range allInfo {
		printList[i] = PrintK8SInfo{
			NAME:     info.NAME,
			Num:      len(info.PortalInfo.K8SServiceList),
			NodeType: make([]string, 0),
			K8SAREA:  make([]string, 0),
		}
		for _, k8sCluster := range info.K8SClusterInfo {
			if !strings.EqualFold(k8sCluster.Config.NodeType, "") {
				printList[i].NodeType = append(printList[i].NodeType, k8sCluster.Config.NodeType)
			}
			if !strings.EqualFold(k8sCluster.Area, "") {
				printList[i].K8SAREA = append(printList[i].K8SAREA, k8sCluster.Area)
			}

		}
	}
	table.Output(printList)
}
