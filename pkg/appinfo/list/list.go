package list

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/client/ews"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/util/metrics"
	systools "goops/pkg/util/sys"
	"goops/pkg/util/table"
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
	ews.AddEWSFlag(cmd.Flags())
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
	case "url":
		printURL(allInformation)
	case "cluster":
		printCluster(allInformation)
	case "app":
		printApp(allInformation)
	default:
		logrus.Error("不支持当前迁移方式")
	}
}

func printApp(allInfo []common.AppInformation) {

	printList := make([]PrintAppInfo, len(allInfo))

	for i, info := range allInfo {
		printList[i] = PrintAppInfo{
			NAME:       info.NAME,
			APPID:      info.APPID,
			Creator:    fmt.Sprintf("%s(%d)", info.PortalInfo.APP.CreatorName, info.PortalInfo.APP.Creator),
			K8SNum:     len(info.PortalInfo.K8SServiceList),
			EWSNum:     len(info.PortalInfo.EWSServiceList),
			MigrateMsg: info.GetMigrateMessage(),
		}
	}
	table.Output(printList)
}

func printCluster(allInfo []common.AppInformation) {
	printList := make([]PrintClusterInfo, 0)

	for _, info := range allInfo {
		var msg string
		if info.PortalInfo.APP.SingleInstance {
			msg = color.HiRedString("单实例")
		}

		// 获取EWS集群信息
		if strings.EqualFold(clusterType, "all") || strings.EqualFold(clusterType, "ews") {
			for _, ewsC := range info.EWSClusterInfo {
				packageUrl, version := ewsC.GetPackageVersion()
				ewsCP := PrintClusterInfo{
					APPNAME:    info.NAME,
					NAME:       ewsC.ShortClusterName(),
					ID:         ewsC.ID,
					Creator:    fmt.Sprintf("%s(%d)", info.PortalInfo.APP.CreatorName, info.PortalInfo.APP.Creator),
					Memory:     metrics.FormatByte(int64(ewsC.Configs.Memory) * 1024 * 1024),
					Instances:  len(ewsC.Instances),
					PackageURL: packageUrl,
					Version:    version,
					MigrateMsg: msg + " " + ewsC.IsPackageNotFound(),
				}
				printList = append(printList, ewsCP)
			}
		}

		// 获取K8S集群信息
		if strings.EqualFold(clusterType, "all") || strings.EqualFold(clusterType, "k8s") {
			for _, k8sC := range info.K8SClusterInfo {
				k8sCP := PrintClusterInfo{
					APPNAME:    info.NAME,
					NAME:       k8sC.ShortClusterName(),
					ID:         k8sC.ID,
					Creator:    fmt.Sprintf("%s(%d)", info.PortalInfo.APP.CreatorName, info.PortalInfo.APP.Creator),
					K8SArea:    k8sC.Area,
					Memory:     k8sC.Config.ContainerTemplate.Resources.Limits.Memory["amount"],
					Instances:  k8sC.Replica,
					NodeType:   k8sC.Config.NodeType,
					NodeName:   k8sC.Config.NodeName,
					MigrateMsg: msg,
				}
				printList = append(printList, k8sCP)
			}
		}
	}

	table.Output(printList)
}

func printURL(allInfo []common.AppInformation) {
	printList := make([]PrintAppURL, len(allInfo))
	for i, info := range allInfo {
		printList[i] = PrintAppURL{
			NAME: info.NAME,
			URL:  info.URL,
		}
	}
	table.Output(printList)
}
