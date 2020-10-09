package tools

import (
	"fmt"
	"kube-tools/src/config"
	"os"
	"strings"
)

/**
检查磁盘空间使用情况
*/

func Df() {

	pods, _ := GetShellPodList()
	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.Contains(pod.Name, config.ShellPodName) &&
			(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
			fmt.Println("------------------------------> No.", i, " Shell on node: "+pod.Status.HostIP+" <------------------------------")
			catExecOps := ExecOptions{
				Command:       "df -lh ",
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
