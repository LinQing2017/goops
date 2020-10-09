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
	if pods != nil {
		for i := 0; i < len(pods.Items); i++ {
			pod := pods.Items[i]
			tChan := make(chan int)
			if strings.Contains(pod.Name, config.ShellPodName) &&
				(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
				fmt.Println("------------------------------> No.", i, " Shell on node: "+pod.Status.HostIP+" <------------------------------")
				dfExecOps := ExecOptions{
					Command:       "df -lh -x tmpfs -T",
					ContainerName: "",
					In:            nil,
					Out:           os.Stdout,
					Err:           os.Stderr,
					Istty:         false,
				}
				go ExecCmdParallel(&pod, dfExecOps, tChan)
				WaitAllThreadFinish(1, tChan, 10)
			}
		}
	}
}
