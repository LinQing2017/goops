package check

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
	v1 "k8s.io/api/core/v1"
	"kube-tools/src/cmd/check/entity"
	"kube-tools/src/util"
	"strconv"
	"strings"
)

type AppInfo struct {
	Name    string
	APPID   string
	PodNode []string
	Type    []string
}

var (
	mongoUser   string
	mongoPasswd string
	mongoUrl    string
	mongoDB     string
	nodeFilter  string
)

func addAppInfoFlag(flags *pflag.FlagSet) {
	flags.StringVar(&mongoUser, "mongo-user", "qa_mdb_k8s_paas_api", "SDP-K8S Mongodb 用户名。")
	flags.StringVar(&mongoPasswd, "mongo-password", "lcsX5zq9QEUW", "SDP-K8S Mongodb 密码。")
	flags.StringVar(&mongoDB, "mongo-database", "qa_mdb_k8s_paas_api", "SDP-K8S Mongodb 数据库名称。")
	flags.StringVar(&mongoUrl, "mongo-url", "m3.all.debug2.mongod3.sdp:34003,m1.all.debug2.mongod3.sdp:34003,m2.all.debug2.mongod3.sdp:34003", "SDP-K8S Mongodb 地址。")
	flags.StringVar(&nodeFilter, "node-filter", "172.24.135.12", "根据node名称过滤信息，使用逗号分隔")
}

func NewCmdAppInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "appinfo",
		Short:                 "获取集群应用信息。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunAppInfo(cmd, args)
		},
	}
	addAppInfoFlag(cmd.Flags())
	return cmd
}

func RunAppInfo(cmd *cobra.Command, args []string) {
	clusterInfoDict := getClusterInfo()
	kubeClient, _ := util.KubeClient(cmd)
	appinfoList := make([]AppInfo, 0)

	pods, _ := util.GetPodList(kubeClient, "", "")
	podNSDict := make(map[string][]v1.Pod)
	for _, pod := range pods.Items {
		key := pod.Namespace
		podListOnNS := podNSDict[key]
		if podListOnNS == nil {
			podListOnNS = make([]v1.Pod, 0)
		}
		podListOnNS = append(podListOnNS, pod)
		podNSDict[key] = podListOnNS
	}

	for namespace, podsOnNS := range podNSDict {

		clist := clusterInfoDict[namespace]
		if clist == nil {
			continue
		}
		appinfo := AppInfo{
			Name:    clist[0].AppName,
			APPID:   namespace,
			PodNode: make([]string, 0),
		}

		for _, pod := range podsOnNS {
			if pod.Status.Phase != v1.PodRunning {
				continue
			}
			appinfo.Type = labelParse(pod.Spec.NodeSelector)
			appinfo.PodNode = append(appinfo.PodNode, pod.Status.HostIP)
		}
		appinfo.PodNode = util.RemoveRepeatedElement(appinfo.PodNode)
		appinfoList = append(appinfoList, appinfo)

	}

	appinfoFilter := filter(appinfoList)
	excelOutput(appinfoFilter)
}

func getClusterInfo() map[string][]entity.ClusterInfo {

	// 设置客户端连接配置
	mongoURI := "mongodb://" + mongoUser + ":" + mongoPasswd + "@" + mongoUrl + "/" + mongoDB + "?autoConnectRetry=true"
	client := util.MongoClient(mongoURI)

	// 指定获取要操作的数据集
	collection := client.Database(mongoDB).Collection("cluster")
	//filter := bson.M{"env_type": "wx"}
	filter := bson.M{}
	records, _ := collection.Find(context.TODO(), filter)

	dict := map[string][]entity.ClusterInfo{}

	for records.Next(context.TODO()) {
		d := entity.ParseClusterInfo(&records.Current)
		clusters := dict[d.AppId]
		if clusters == nil {
			clusters = make([]entity.ClusterInfo, 0)
		}
		dict[d.AppId] = append(clusters, d)
	}
	// 断开连接
	util.MongoDisconnect(client)

	return dict
}

func labelParse(labels map[string]string) []string {
	typeLabels := make([]string, 0)
	for k, v := range labels {
		if strings.EqualFold(v, "type") {
			typeLabels = append(typeLabels, k)
		}
	}
	return typeLabels
}

func filter(appInfoList []AppInfo) []AppInfo {
	if strings.EqualFold(nodeFilter, "") {
		return appInfoList
	}
	nodeFilterSet := mapset.NewSet()
	for _, v := range strings.Split(nodeFilter, ",") {
		nodeFilterSet.Add(v)
	}

	appInfoListFilter := make([]AppInfo, 0)
	for _, appInfo := range appInfoList {
		for _, node := range appInfo.PodNode {
			if nodeFilterSet.Contains(node) {
				appInfoListFilter = append(appInfoListFilter, appInfo)
				break
			}
		}
	}
	return appInfoListFilter
}

func excelOutput(appInfoList []AppInfo) {

	title := map[string]string{"A1": "序号", "B1": "应用名称", "C1": "命名空间", "D1": "Pod运行节点", "E1": "应用类型"}
	f := excelize.NewFile()
	for k, v := range title {
		f.SetCellValue("Sheet1", k, v)
	}
	id := 1
	for _, appinfo := range appInfoList {

		rowNum := strconv.Itoa(id + 1)
		row := map[string]string{
			"A" + rowNum: strconv.Itoa(id),
			"B" + rowNum: appinfo.Name,
			"C" + rowNum: appinfo.APPID,
			"D" + rowNum: strings.Join(appinfo.PodNode, ","),
			"E" + rowNum: strings.Join(appinfo.Type, ","),
		}
		for k, v := range row {
			f.SetCellValue("Sheet1", k, v)
		}
		id += 1
	}
	if err := f.SaveAs("output" + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
