package tools

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"strings"
)

func Cat() {
	filePath := flag.Arg(1)
	if strings.EqualFold(filePath, "") {
		return
	}
	pods, err := config.KubeClientSet.CoreV1().Pods(*config.ShellNamespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if pods.Size() == 0 {
		fmt.Println("Node Shell 工具没有安装")
		return
	}
	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.Contains(pod.Name, "node-shell-tool") {
			fmt.Println("------------------------------>", pod.Status.HostIP, "<------------------------------")
			ExecCmd(&pod, "cat "+filePath, "")
		}

	}
}
