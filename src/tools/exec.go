package tools

import (
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"kube-tools/src/config"
	"net/url"
	"strings"
)

type ExecOptions struct {
	Command       string
	ContainerName string
	In            io.Reader
	Out           io.Writer
	Err           io.Writer
	Istty         bool
}

func ExexCmdParallel(pod *v1.Pod, execOptions ExecOptions, tChan chan int) {
	err := ExecCmd(pod, execOptions)

	if err != nil {
		fmt.Println("请求 API Service 返回异常：", pod.Status.HostIP)
		panic(err.Error())
	}

	tChan <- 1
}

func ExecCmd(pod *v1.Pod, execOptions ExecOptions) error {

	// 获取pod中的目标Container
	container, _ := containerToExec(execOptions.ContainerName, pod)
	// 创建运行表达式
	podOptions := v1.PodExecOptions{
		Command:   strings.Fields(execOptions.Command),
		Container: container.Name,
		Stdin:     execOptions.In != nil,
		Stdout:    execOptions.Out != nil,
		Stderr:    execOptions.Err != nil,
		TTY:       execOptions.Istty,
	}

	// 创建客户端请求
	req := config.KubeClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&podOptions, scheme.ParameterCodec)

	// 执行命令，并输出到标准输出
	streamOptions := getStreamOptions(&podOptions, execOptions.In, execOptions.Out, execOptions.Err)
	return startStream("POST", req.URL(), config.KubeConfig, streamOptions)
}

func containerToExec(container string, pod *v1.Pod) (*v1.Container, error) {
	if len(container) > 0 {
		for i := range pod.Spec.Containers {
			if pod.Spec.Containers[i].Name == container {
				return &pod.Spec.Containers[i], nil
			}
		}
		for i := range pod.Spec.InitContainers {
			if pod.Spec.InitContainers[i].Name == container {
				return &pod.Spec.InitContainers[i], nil
			}
		}
		return nil, fmt.Errorf("container not found (%s)", container)
	}
	return &pod.Spec.Containers[0], nil
}

func getStreamOptions(execOptions *v1.PodExecOptions, stdin io.Reader, stdout, stderr io.Writer) remotecommand.StreamOptions {

	var streamOptions remotecommand.StreamOptions

	if execOptions.Stdin {
		streamOptions.Stdin = stdin
	}

	if execOptions.Stdout {
		streamOptions.Stdout = stdout
	}

	if execOptions.Stderr {
		streamOptions.Stderr = stderr
	}

	return streamOptions
}

func startStream(method string, url *url.URL, config *restclient.Config, streamOptions remotecommand.StreamOptions) error {

	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}

	return exec.Stream(streamOptions)
}
