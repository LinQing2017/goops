package tools

import (
	"bytes"
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
	tChan := make(chan int, len(pods.Items))
	outPutBuffers := make([]*bytes.Buffer, len(pods.Items))
	threadNum := 0

	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		if strings.Contains(pod.Name, config.ShellPodName) &&
			(strings.EqualFold(*config.NodeIP, "") || strings.EqualFold(*config.NodeIP, pod.Status.HostIP)) {
			outPutBuffer := bytes.NewBufferString("------------------------------> Shell on node: " + pod.Status.HostIP + "<------------------------------\n")
			outPutBuffers[i] = outPutBuffer
			catExecOps := ExecOptions{
				Command:       "cat " + filePath,
				ContainerName: "",
				In:            nil,
				Out:           outPutBuffer,
				Err:           os.Stderr,
				Istty:         false,
			}
			go ExexCmdParallel(&pod, catExecOps, tChan)
			threadNum += 1
			if threadNum%10 == 0 {
				// 5个线程进行并发
				WaitAllThreadFinish(10, tChan, 30)
			} else if threadNum == len(pods.Items) {
				// 等待剩余线程完成
				WaitAllThreadFinish(threadNum%10, tChan, 30)
			}
		}
	}
	//WaitAllThreadFinish(threadNum, tChan, 30)
	close(tChan)

	for i := 0; i < threadNum; i++ {
		outPutBuffer := outPutBuffers[i]
		fmt.Println(outPutBuffer.String())
	}
}
