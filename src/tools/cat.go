package tools

import (
	"bytes"
	"flag"
	"fmt"
	"kube-tools/src/config"
	"os"
	"strings"
)

func Cat() {
	filePath := flag.Arg(1)
	if strings.EqualFold(filePath, "") {
		return
	}
	pods, _ := GetShellPodList()
	outPutBuffers := make([]*bytes.Buffer, len(pods.Items))

	//tChan := make(chan int, len(pods.Items))
	//threadNum := 0

	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.Contains(pod.Name, config.ShellPodName) &&
			(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
			outPutBuffer := bytes.NewBufferString("------------------------------> Shell on node: " + pod.Status.HostIP + " <------------------------------\n")
			outPutBuffers[i] = outPutBuffer
			catExecOps := ExecOptions{
				Command:       "cat " + filePath,
				ContainerName: "",
				In:            nil,
				Out:           outPutBuffer,
				Err:           os.Stderr,
				Istty:         false,
			}
			ExecCmd(&pod, catExecOps)
			//threadNum += 1
			//go ExecCmdParallel(&pod, catExecOps, tChan)
			//if threadNum%config.ConcurrentThreadNum == 0 {
			//	// 5个线程进行并发
			//	WaitAllThreadFinish(config.ConcurrentThreadNum, tChan, 60)
			//} else if threadNum == len(pods.Items) {
			//	// 等待剩余线程完成
			//	WaitAllThreadFinish(threadNum%config.ConcurrentThreadNum, tChan, 60)
			//}
		}

	}

	//close(tChan)
	for i := 0; i < len(outPutBuffers); i++ {
		outPutBuffer := outPutBuffers[i]
		if outPutBuffer != nil {
			fmt.Println(outPutBuffer.String())
		}
	}
}
