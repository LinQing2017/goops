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
	Name       string
	ROLE       string
	ENV_LABEL  string
	TYPE_LABEL string
	CPU        string
	Memory     string
	POD        string
	SHELL_POD  string
}

func Node() {
	nodes, err := config.KubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeInfoList := make([]NodeInfo, len(nodes.Items))
	shellPods, err := GetShellPodList()
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
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

		shellPod := ""
		if shellPods != nil && err == nil {
			for _, pod := range shellPods.Items {
				if strings.EqualFold(pod.Status.HostIP, node.Name) && pod.Status.Phase == v1.PodRunning {
					shellPod = pod.Name
					break
				}
			}
		}

		pods, _ := config.KubeClientSet.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + node.Name,
		})

		nodeInfo := NodeInfo{
			node.Name,
			role,
			strings.Join(env_lable, ","),
			strings.Join(type_lable, ","),
			node.Status.Capacity.Cpu().String(),
			strconv.FormatFloat(float64(node.Status.Capacity.Memory().Value())/1024/1024/1024, 'f', 2, 64) + " Gi",
			strconv.Itoa(len(pods.Items)) + "/" + node.Status.Capacity.Pods().String(),
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
