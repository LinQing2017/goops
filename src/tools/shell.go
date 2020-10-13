package tools

import (
	"flag"
	"fmt"
	"kube-tools/src/config"
	"os"
	"strings"
)

func Shell() {
	pods, _ := GetShellPodList()
	cmdStr := strings.Join(flag.Args()[1:], " ")
	if pods != nil {
		for i := 0; i < len(pods.Items); i++ {
			pod := pods.Items[i]
			tChan := make(chan int)
			if strings.Contains(pod.Name, config.ShellPodName) &&
				(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
				fmt.Println("------------------------------> No.", i, " Shell on node: "+pod.Status.HostIP+" <------------------------------")
				dfExecOps := ExecOptions{
					Command:       cmdStr,
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
