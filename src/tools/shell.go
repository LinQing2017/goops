package tools

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"os"
	"strings"
)

func Shell() {
	pods, err := config.KubeClientSet.CoreV1().Pods(config.ShellNamespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: "name=" + config.ShellPodName,
	})

	if err != nil {
		fmt.Println("请先安装shell工具")
		panic(err.Error())
	}

	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.EqualFold(pod.Status.HostIP, *config.NodeIP) {
			ExecCmd(&pod, "bash", "", os.Stdin, true)
			return
		}
	}
}
