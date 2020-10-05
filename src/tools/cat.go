package tools

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"os"
	"strings"
)

func Cat() {
	filePath := flag.Arg(1)
	if strings.EqualFold(filePath, "") {
		return
	}
	pods, err := config.KubeClientSet.CoreV1().Pods(config.ShellNamespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: "name=" + config.ShellPodName,
	})
	if err != nil {
		panic(err.Error())
	}
	if pods.Size() == 0 {
		fmt.Println("Node Shell 工具没有安装")
		return
	}
	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.Contains(pod.Name, config.ShellPodName) &&
			(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
			fmt.Println("------------------------------>", pod.Status.HostIP, "<------------------------------")
			catExecOps := ExecOptions{
				Command:       "cat " + filePath,
				ContainerName: "",
				In:            nil,
				Out:           os.Stdout,
				Err:           os.Stderr,
				Istty:         false,
			}
			ExecCmd(&pod, catExecOps)
		}
	}
}
