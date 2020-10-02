package tools

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"strings"
)

func Node() {
	nodes, err := config.KubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < len(nodes.Items); i++ {
		node := nodes.Items[i]
		fmt.Println("=======================", nodes.Items[i].Name, "=======================")
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
		fmt.Println("ENV:", strings.Join(env_lable, ","))
		fmt.Println("TYPE:", strings.Join(type_lable, ","))
	}

}
