package list

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/ews_client"
	k8s_tools "goops/pkg/util/kubernetes"
	systools "goops/pkg/util/sys"
	"goops/pkg/util/table"
	"strconv"
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
			NAME:    info.NAME,
			APPID:   info.APPID,
			Creator: fmt.Sprintf("%s(%d)", info.PortalInfo.APP.CreatorName, info.PortalInfo.APP.Creator),
			Single:  "",
			K8SNum:  len(info.PortalInfo.K8SServiceList),
			EWSNum:  len(info.PortalInfo.EWSServiceList),
			URL:     "",
		}
		if withUrl {
			printList[i].URL = info.URL
		}
		if info.PortalInfo.APP.SingleInstance {
			printList[i].Single = "Y"
		}

		cmptDomain := 0
		domainNum := 0
		for _, domains := range info.ClusterBindDomains {

			for _, domain := range domains {
				if domain.IsCmptDomain() {
					cmptDomain++
					continue
				}
				if !domain.IsTestDomain() {
					domainNum++
				}

			}
		}
		printList[i].DomainNum = strconv.Itoa(domainNum)
		if cmptDomain > 0 {
			printList[i].DomainNum = fmt.Sprintf("%d(cmpt:%d)", domainNum, cmptDomain)
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
		// 所有集群都能找到一个对应实例，有包地址。否则认为应用缺少war包
		if info.IsPackageNotFound() {
			printList[i].PackageURL = color.HiRedString("Not Found")
		}
	}
	table.Output(printList)
}

func printK8s(allInfo []common.AppInformation) {

	printList := make([]PrintK8SInfo, len(allInfo))
	kubeClient, _ := k8s_tools.KubeClientAndConfig(kubeConfig)
	podDict, _ := k8s_tools.GetPodDictByNamespace(kubeClient, "")

	for i, info := range allInfo {

		pods := podDict[info.APPID]
		running := 0
		for _, pod := range pods {
			if pod.Status.Phase == "Running" {
				running++
			}
		}
		var podStatus string
		if running == len(pods) {
			podStatus = color.GreenString(fmt.Sprintf("%d/%d", running, len(pods)))
		} else {
			podStatus = color.HiRedString(fmt.Sprintf("%d/%d", running, len(pods)))
		}
		printList[i] = PrintK8SInfo{
			NAME:      info.NAME,
			Num:       len(info.PortalInfo.K8SServiceList),
			PodStatus: podStatus,
			NodeType:  make([]string, 0),
			K8SAREA:   make([]string, 0),
		}
		for _, k8sCluster := range info.K8SClusterInfo {
			if !strings.EqualFold(k8sCluster.Config.NodeType, "") {
				printList[i].NodeType = append(printList[i].NodeType, k8sCluster.Config.NodeType)
			}
			if !strings.EqualFold(k8sCluster.Area, "") {
				printList[i].K8SAREA = append(printList[i].K8SAREA, k8sCluster.Area)
			}
			if !strings.EqualFold(k8sCluster.Config.NodeName, "") {
				printList[i].NodeName = append(printList[i].NodeName, k8sCluster.Config.NodeName)
			}
			if k8sCluster.Config.Hpa.Enable {
				printList[i].HPA += fmt.Sprintf("k8s-%s", k8sCluster.ID[len(k8sCluster.ID)-6:])
			}
		}
	}

	table.Output(printList)
}
