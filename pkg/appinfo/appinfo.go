package appinfo

import (
	"context"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	mapset "github.com/deckarep/golang-set"
	"github.com/modood/table"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"kube-tools/pkg/appinfo/entity"
	k8stools "kube-tools/pkg/util/kubernetes"
	"kube-tools/pkg/util/mongo"
	systools "kube-tools/pkg/util/sys"
	"sort"
	"strconv"
	"strings"
)

func Main(cmd *cobra.Command, args []string) {
	clusterInfoDict := getClusterInfo()
	kubeClient, _ := k8stools.KubeClientAndConfig(kubeConfigStr)
	podNSDict := getPodNSDict(kubeClient)

	appinfoList := make([]entity.AppInfo, 0)
	for namespace, podsOnNS := range podNSDict {

		clist := clusterInfoDict[namespace]
		if clist == nil {
			continue
		}
		appinfo := entity.AppInfo{
			Name:    clist[0].AppName,
			APPID:   namespace,
			PodNode: make([]string, 0),
			Type:    make([]string, 0),
		}
		for _, pod := range podsOnNS {
			if pod.Status.Phase != v1.PodRunning {
				continue
			}
			appinfo.PodNode = append(appinfo.PodNode, pod.Status.HostIP)
			for k, v := range pod.Spec.NodeSelector {
				if strings.EqualFold(v, "type") {
					appinfo.Type = append(appinfo.Type, k)
				}
			}
		}
		appinfo.PodNode = systools.RemoveRepeatedElement(appinfo.PodNode)
		sort.Strings(appinfo.PodNode)
		appinfo.Type = systools.RemoveRepeatedElement(appinfo.Type)
		appinfoList = append(appinfoList, appinfo)
	}

	appinfoFilter := filter(appinfoList)
	excelAppInfo(appinfoFilter)
	if isPrint {
		table.Output(appinfoFilter)
	}
}

func getClusterInfo() map[string][]entity.ClusterInfo {

	// 设置客户端连接配置
	mongoURI := "mongodb://" + mongoUser + ":" + mongoPasswd + "@" + mongoUrl + "/" + mongoDB + "?autoConnectRetry=true"
	client := mongo.MongoClient(mongoURI)

	// 指定获取要操作的数据集
	collection := client.Database(mongoDB).Collection("cluster")
	records, _ := collection.Find(context.TODO(), bson.M{})
	k8sClusters := make([]entity.ClusterInfo, 0)
	records.All(context.TODO(), &k8sClusters)
	dict := map[string][]entity.ClusterInfo{}
	for _, c := range k8sClusters {
		clusters := dict[c.AppId]
		if clusters == nil {
			clusters = make([]entity.ClusterInfo, 0)
		}
		dict[c.AppId] = append(clusters, c)
	}
	// 断开连接
	mongo.MongoDisconnect(client)

	return dict
}

func getPodNSDict(kubeClient *kubernetes.Clientset) map[string][]v1.Pod {
	pods, _ := k8stools.GetPodList(kubeClient, "", "")
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
	return podNSDict
}

func filter(appInfoList []entity.AppInfo) []entity.AppInfo {
	if strings.EqualFold(nodeFilter, "") {
		return appInfoList
	}
	nodeFilterSet := mapset.NewSet()
	for _, v := range strings.Split(nodeFilter, ",") {
		nodeFilterSet.Add(v)
	}

	appInfoListFilter := make([]entity.AppInfo, 0)
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

func excelAppInfo(appInfoList []entity.AppInfo) {

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
		logrus.Errorf(err.Error())
	}
}
