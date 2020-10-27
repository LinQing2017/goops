package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"kube-tools/src/config"
	error2 "kube-tools/src/error"
	"kube-tools/src/util"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	nodename string
	nodefile string

	httpTimeOutInSec int
	currentThreadNum int
)

func RunShell(cmd *cobra.Command, args []string) {

	kubeClientSet, kubeClientConfig := util.KubeClient(cmd)

	cmdStr := strings.Join(args, " ")

	// 返回所有需要运行运行的Node列表
	shellPodTargets := getPodTargets(kubeClientSet)
	i := 0
	threadNum := 0
	total := len(shellPodTargets)
	tChan := make(chan int, len(shellPodTargets))
	outPutBuffers := make([]*bytes.Buffer, len(shellPodTargets))
	defer close(tChan)

	for n, pod := range shellPodTargets {
		outPutBuffer := bytes.NewBufferString("------------------------------> No." + strconv.Itoa(i) + " Shell on node: " + n + " <------------------------------\n")
		outPutBuffers[i] = outPutBuffer
		if pod != nil {
			dfExecOps := ExecOptions{
				Command:       cmdStr,
				ContainerName: "",
				In:            nil,
				Out:           outPutBuffer,
				Err:           os.Stderr,
				Istty:         false,
			}
			go ExecCmdParallel(kubeClientSet, kubeClientConfig, pod, dfExecOps, tChan)
			threadNum += 1
		} else {
			outPutBuffer.WriteString("Can't find shell pod on " + n + "\n")
		}
		i += 1
		if threadNum == currentThreadNum || total == i {
			util.WaitAllThreadFinish(threadNum, tChan, httpTimeOutInSec)
			threadNum = 0
		}
	}

	for _, output := range outPutBuffers {
		fmt.Print(output.String())
	}
}

func getPodTargets(kubeClientSet *kubernetes.Clientset) map[string]*v1.Pod {
	shellPodDict, _ := util.GetAllPodByNodeName(kubeClientSet, "name="+config.ShellDaemonset)
	nodes, _ := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	nodeTargets := make([]string, 0)

	// 所有Node执行Shell
	if strings.EqualFold(nodename, "") && strings.EqualFold(nodefile, "") {
		for _, node := range nodes.Items {
			nodeTargets = append(nodeTargets, node.Name)
		}
	}

	// 指定Node执行shell
	if !strings.EqualFold(nodename, "") {
		nodeTargets = append(nodeTargets, nodename)
	}

	// nodefile文件指定node执行shell
	if strings.EqualFold(nodename, "") && !strings.EqualFold(nodefile, "") {
		nodeTargets = util.ReadLine(nodefile)
	}
	if len(nodeTargets) == 0 {
		panic(error2.NodeShellError{500, "选择节点异常"})
	}

	// 获取shellPod
	podTarges := make(map[string]*v1.Pod, 0)
	for _, n := range nodeTargets {
		if pod, ok := shellPodDict[n]; ok {
			podTarges[n] = &pod[0]
		} else {
			podTarges[n] = nil
		}
	}
	return podTarges
}

type ExecOptions struct {
	Command       string
	ContainerName string
	In            io.Reader
	Out           io.Writer
	Err           io.Writer
	Istty         bool
}

func ExecCmdParallel(kubeClientSet *kubernetes.Clientset, kubeClientConfig *restclient.Config, pod *v1.Pod, execOptions ExecOptions, tChan chan int) {
	err := ExecCmd(kubeClientSet, kubeClientConfig, pod, execOptions)

	if err != nil {
		fmt.Println("请求 API Service 返回异常：", pod.Status.HostIP)
		//panic(err.Error())
	}

	tChan <- 1
}

func ExecCmd(kubeClientSet *kubernetes.Clientset, kubeClientConfig *restclient.Config, pod *v1.Pod, execOptions ExecOptions) error {

	if pod.Status.Phase != v1.PodRunning {
		fmt.Println("Pod 没有就绪：", pod.Name, pod.Status.HostIP)
		err := &error2.NodeShellError{500, "Pod 没有就绪"}
		return err
	}

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
	req := kubeClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Timeout(time.Duration(httpTimeOutInSec))

	req.VersionedParams(&podOptions, scheme.ParameterCodec)

	// 执行命令，并输出到标准输出
	streamOptions := getStreamOptions(&podOptions, execOptions.In, execOptions.Out, execOptions.Err)
	return startStream("POST", req.URL(), kubeClientConfig, streamOptions)
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

func addFlag(flags *pflag.FlagSet) {
	flags.StringVar(&nodename, "node", "", "在指定宿主机节点执行shell命令。")
	flags.StringVar(&nodefile, "nodefile", "", "通过文件指定要运行shell命令的宿主机。")

	flags.IntVar(&httpTimeOutInSec, "timeout", 30, "连接Kubelet超时时间。")
	flags.IntVar(&currentThreadNum, "thread", 1, "执行shell命令的并发数。")
}

func NewCmdSh() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sh [command]",
		Short:                 "在宿主机的客户端中执行Shell命令",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunShell(cmd, args)
		},
	}
	addFlag(cmd.Flags())
	return cmd
}
