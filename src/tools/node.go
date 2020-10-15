package tools

import (
	"fmt"
	"github.com/modood/table"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"sort"
	"strconv"
	"strings"
)

type NodeInfo struct {
	Name          string
	Role          string
	Unschedulable string
	Env           string
	Type          string
	CPU           string
	Memory        string
	Pod           string
	Shell         string
}

func Node() {
	nodes, err := config.KubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeInfoList := make([]NodeInfo, len(nodes.Items))
	shellPods, err := GetShellPodList()
	allPodDist, _ := getAllPodByNodeName()
	for i, node := range nodes.Items {
		// 获取Role以及Label信息
		role := ""
		env_lable := make([]string, 0)
		type_lable := make([]string, 0)
		for k, v := range node.Labels {
			switch v {
			case "env":
				env_lable = append(env_lable, k)
			case "type":
				type_lable = append(type_lable, k)
			}
			if strings.HasPrefix(k, "node-role.kubernetes.io") {
				role = strings.Split(k, "/")[1]
			}
		}

		// 获取ShellPod名称
		shellPod := ""
		if shellPods != nil && err == nil {
			for _, pod := range shellPods.Items {
				if strings.EqualFold(pod.Status.HostIP, node.Name) && pod.Status.Phase == v1.PodRunning {
					shellPod = pod.Name
					break
				}
			}
		}

		// 列出获取该节点的所有Pod
		podListOnNode := allPodDist[node.Name]

		// 获取节点的状态

		nodeInfo := NodeInfo{
			node.Name,
			role,
			strconv.FormatBool(node.Spec.Unschedulable),
			strings.Join(env_lable, ","),
			strings.Join(type_lable, ","),
			node.Status.Capacity.Cpu().String(),
			strconv.FormatFloat(float64(node.Status.Capacity.Memory().Value())/1024/1024/1024, 'f', 2, 64) + " Gi",
			strconv.Itoa(len(podListOnNode)) + "/" + node.Status.Capacity.Pods().String(),
			shellPod,
		}
		nodeInfoList[i] = nodeInfo
	}
	sort.Slice(nodeInfoList, func(i, j int) bool {
		return nodeInfoList[i].ROLE < nodeInfoList[j].ROLE
	})
	nodeInfoStr := table.Table(nodeInfoList)
	fmt.Println(nodeInfoStr)
}

func getAllPodByNodeName() (podDist map[string][]v1.Pod, err error) {
	pods, err := config.KubeClientSet.CoreV1().Pods("").List(metav1.ListOptions{})
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
