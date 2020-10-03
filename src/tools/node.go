package tools

import (
	"fmt"
	"github.com/modood/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"strings"
)

type NodeInfo struct {
	Name       string
	ENV_LABEL  string
	TYPE_LABEL string
}

func Node() {
	nodes, err := config.KubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeInfoList := make([]NodeInfo, len(nodes.Items))
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
		env_lable := make([]string, 0)
		type_lable := make([]string, 0)
		for k, v := range node.Labels {
			switch v {
			case "env":
				env_lable = append(env_lable, k)
			case "type":
				type_lable = append(type_lable, k)
			default:
				continue
			}
		}
		nodeInfo := NodeInfo{
			node.Name,
			strings.Join(env_lable, ","),
			strings.Join(type_lable, ","),
		}
		nodeInfoList[i] = nodeInfo
	}
	nodeInfoStr := table.Table(nodeInfoList)
	fmt.Println(nodeInfoStr)
}
