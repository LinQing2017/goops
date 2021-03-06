package migrate_k8s

import (
	"context"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	k8stools "goops/pkg/util/kubernetes"
	mongotools "goops/pkg/util/mongo"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func Main(cmd *cobra.Command, args []string) {
	// 设置客户端连接配置
	mongoURI := "mongodb://" + mongoUser + ":" + mongoPasswd + "@" + mongoUrl + "/" + mongoDB + "?autoConnectRetry=true"
	ndpPortalClient := mongotools.MongoClient(mongoURI)

	k8sClient, _ := k8stools.KubeClientAndConfig(kubeConfigStr)
	podDict, _ := k8stools.GetPodDict(k8sClient, "")
	podDictByNamespace, _ := k8stools.GetPodDictByNamespace(k8sClient, "")

	// 读取流量信息
	appMetrics := getAppMetrics()

	appInfoDict := make(map[string][]AppInformathion, 0)
	for nodename, pods := range podDict {

		appInfoList := make([]AppInformathion, 0)
		for _, pod := range pods {
			namespace := pod.Namespace
			// 过滤掉转码任务
			nodeSelector := getNodeSelectors(pod)
			if strings.EqualFold(nodeSelector, "cpu") {
				continue
			}
			// 检查其他Pod运行节点
			otherPods := podDictByNamespace[namespace]
			otherHostIP := make([]string, 0)
			for _, otherPod := range otherPods {
				if !strings.EqualFold(otherPod.Name, pod.Name) {
					otherHostIP = append(otherHostIP, otherPod.Status.HostIP)
				}
			}
			// 获取APP信息
			var app types.App
			objectId, _ := primitive.ObjectIDFromHex(namespace)
			if err := db_tools.GetOne(mongoDB, "app", bson.M{"_id": objectId}, ndpPortalClient, &app); err == nil {
				appInfo := AppInformathion{
					AppId:             namespace,
					Name:              app.Name,
					HostIP:            nodename,
					OtherIP:           strings.Join(otherHostIP, "\n"),
					NodeSelectorLabel: nodeSelector,
					Metric:            appMetrics[app.Name],
					Creator:           app.CreatorName + "(" + strconv.Itoa(app.Creator) + ")",
					URL:               "https://da.sdp.101.com/#/ndpfront/applicationManagement/applicationList/serviceInformation/" + namespace + "/" + app.Name,
					Single:            app.SingleInstance,
				}
				appInfoList = append(appInfoList, appInfo)
			}
		}
		if len(appInfoList) > 0 {
			appInfoDict[nodename] = appInfoList
		}
		//table.Output(appInfoList)
	}
	excelAppInfo(appInfoDict)
	ndpPortalClient.Disconnect(context.TODO())
	getAppMetrics()
}

func getAppMetrics() map[string]string {

	appMetrics := make(map[string]string)
	if f, err := excelize.OpenFile(metricFile, excelize.Options{}); err != nil {
		logrus.Error(err.Error())
	} else {
		if rows, err := f.Rows("nginx"); err == nil {
			rows.Columns()
			for rows.Next() {
				columns, _ := rows.Columns()
				appMetrics[columns[0]] = columns[2]
			}

		}
	}
	return appMetrics
}

func getNodeSelectors(pod v1.Pod) string {
	for k, v := range pod.Spec.NodeSelector {
		if strings.EqualFold(v, "type") {
			return k
		}
	}
	return ""
}

func excelAppInfo(appInfoDict map[string][]AppInformathion) {

	sheetTitle := map[string]string{"A1": "应用名称", "B1": "运行节点", "C1": "其他实例运行节点", "D1": "标签", "E1": "创建人", "F1": "单实例", "G1": "访问量", "H1": "链接"}
	f := excelize.NewFile()
	for nodename, appInfoList := range appInfoDict {
		f.NewSheet(nodename)
		for k, v := range sheetTitle {
			f.SetCellValue(nodename, k, v)
		}
		id := 1
		for _, appinfo := range appInfoList {
			rowNum := strconv.Itoa(id + 1)
			row := map[string]string{
				"A" + rowNum: appinfo.Name,
				"B" + rowNum: appinfo.HostIP,
				"C" + rowNum: appinfo.OtherIP,
				"D" + rowNum: appinfo.NodeSelectorLabel,
				"E" + rowNum: appinfo.Creator,
				"F" + rowNum: strconv.FormatBool(appinfo.Single),
				"G" + rowNum: appinfo.Metric,
				"H" + rowNum: "链接",
			}
			for k, v := range row {
				f.SetCellValue(nodename, k, v)
			}
			f.SetCellHyperLink(nodename, "H"+rowNum, appinfo.URL, "External")
			id += 1
		}
	}
	if err := f.SaveAs("output" + ".xlsx"); err != nil {
		logrus.Errorf(err.Error())
	}
}

func NewCmdMigrateK8s() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "migrate-k8s",
		Short:                 "Kubernetes集群迁移Kubernetes集群时相关信息导出",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	AddMigrateFlag(cmd.Flags())
	return cmd
}
