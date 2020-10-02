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
	"os"
	"strings"
)

func ExecCmd(pod *v1.Pod, command string, containerName string) {

	// 获取pod中的目标Container
	container, _ := containerToExec(containerName, pod)

	// 创建运行表达式
	execOptions := v1.PodExecOptions{
		Command:   strings.Fields(command),
		Container: container.Name,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}

	// 创建客户端请求
	req := config.KubeClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&execOptions, scheme.ParameterCodec)

	// 执行命令，并输出到标准输出
	streamOptions := getStreamOptions(&execOptions, nil, os.Stdout, os.Stderr)
	startStream("POST", req.URL(), config.KubeConfig, streamOptions)
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
