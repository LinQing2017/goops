package util

import (
	"github.com/modood/table"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kube-tools/src/config"
	error2 "kube-tools/src/error"
	"strings"
)

type PodSimpleInfo struct {
	Name   string
	Status string
	Node   string
}

// 获取当前所有 shell pod 的列表
func GetShellPodList(kubeClientSet *kubernetes.Clientset) (*v1.PodList, error) {

	_, err := kubeClientSet.CoreV1().Namespaces().Get(config.ShellNamespace, metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
	})
	if err != nil {
		err = &error2.NodeShellError{500, "Node Shell 工具没有安装或者有异常"}
		return nil, err
	}
	var pods *v1.PodList
	pods, err = kubeClientSet.CoreV1().Pods(config.ShellNamespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: "name=" + config.ShellDaemonset,
	})
	if err != nil || len(pods.Items) == 0 {
		err = &error2.NodeShellError{500, "Node Shell 工具没有安装或者有异常"}
	}
	return pods, err
}

// 根据Node返回当前Pod运行字典
func GetPodDict(kubeClientSet *kubernetes.Clientset, lableSelector string) (podDist map[string][]v1.Pod, err error) {
	listOptions := metav1.ListOptions{}
	if !strings.EqualFold(lableSelector, "") {
		listOptions = metav1.ListOptions{
			TypeMeta:      metav1.TypeMeta{},
			LabelSelector: lableSelector,
		}
	}

	pods, err := kubeClientSet.CoreV1().Pods("").List(listOptions)
	podDist = make(map[string][]v1.Pod)
	for _, pod := range pods.Items {
		key := pod.Spec.NodeName
		podListOnNode := podDist[key]
		if podListOnNode == nil {
			podListOnNode = make([]v1.Pod, 0)
		}
		podListOnNode = append(podListOnNode, pod)
		podDist[key] = podListOnNode
	}
	return
}

// 打印Pod的状态信息
func PrintPodSimpleInfo(kubeClientSet *kubernetes.Clientset, namespace, lableSelector string) {

	pods, err := kubeClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: lableSelector,
	})
	if err != nil {
		panic(err.Error())
	}
	podInfoList := make([]PodSimpleInfo, len(pods.Items))
	for i := 0; i < len(pods.Items); i++ {
		podInfo := PodSimpleInfo{
			Name:   pods.Items[i].Name,
			Status: string(pods.Items[i].Status.Phase),
			Node:   pods.Items[i].Status.HostIP,
		}
		podInfoList[i] = podInfo

	}
	table.Output(podInfoList)
}

// 返回目标节点（Node List）的shell pod列表
func GetShellPodDict(kubeClientSet *kubernetes.Clientset) map[string]*v1.Pod {

	shellPods, _ := GetShellPodList(kubeClientSet)

	// 根据参数判断需要传在那些Node执行Shell
	nodeList := make([]string, 0)
	// 所有Node执行Shell
	if strings.EqualFold(config.ShellNodeName, "") && strings.EqualFold(config.ShellNodeNameFile, "") {
		nodes, _ := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			nodeList = append(nodeList, node.Name)
		}
	}

	// 指定Node执行shell
	if !strings.EqualFold(config.ShellNodeName, "") {
		nodeList = append(nodeList, config.ShellNodeName)
	}

	// nodefile文件指定node执行shell
	if strings.EqualFold(config.ShellNodeName, "") && !strings.EqualFold(config.ShellNodeNameFile, "") {
		nodeList = ReadLine(config.ShellNodeNameFile)
	}
	if len(nodeList) == 0 {
		panic(error2.NodeShellError{500, "选择节点异常"})
	}

	// 获取shellPod
	podTarges := make(map[string]*v1.Pod, 0)
	for _, nodename := range nodeList {
		podTarges[nodename] = nil
	}

	// 填充Shellpod
	for i, shellPod := range shellPods.Items {
		_, isOk := podTarges[shellPod.Spec.NodeName]
		if isOk {
			podTarges[shellPod.Spec.NodeName] = &shellPods.Items[i]
		}
	}
	return podTarges
}
