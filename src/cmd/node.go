package cmd

import (
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/util"
	"sort"
	"strconv"
	"strings"
)

type NodeInfo struct {
	Name   string
	Role   string
	UnSche string
	Env    string
	Type   string
	CPU    string
	Memory string
	Pod    string
	Shell  string
}

func RunNode(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := util.KubeClient(cmd)

	nodes, err := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeInfoList := make([]NodeInfo, len(nodes.Items))
	shellPods, err := util.GetShellPodList(kubeClientSet)
	allPodDist, _ := util.GetAllPodByNodeName(kubeClientSet, "")
	for i, node := range nodes.Items {
		// 获取Role以及Label信息
		role := ""
		env_lable := make([]string, 0)
		type_lable := make([]string, 0)
		for k, v := range node.Labels {
			switch v {
			case "env":
				env_lable = append(env_lable, k)
			case "type":
				type_lable = append(type_lable, k)
			}
			if strings.HasPrefix(k, "node-role.kubernetes.io") {
				role = strings.Split(k, "/")[1]
			}
		}

		// 获取ShellPod名称
		shellPod := ""
		if shellPods != nil && err == nil {
			for _, pod := range shellPods.Items {
				if strings.EqualFold(pod.Status.HostIP, node.Name) && pod.Status.Phase == v1.PodRunning {
					shellPod = pod.Name
					break
				}
			}
		}

		// 列出获取该节点的所有Pod
		podListOnNode := allPodDist[node.Name]

		// 获取节点的状态
		var unschedulable = ""
		if node.Spec.Unschedulable {
			unschedulable = "Y"
		}

		nodeInfo := NodeInfo{
			node.Name,
			role,
			unschedulable,
			strings.Join(env_lable, ","),
			strings.Join(type_lable, ","),
			node.Status.Capacity.Cpu().String(),
			strconv.FormatFloat(float64(node.Status.Capacity.Memory().Value())/1024/1024/1024, 'f', 2, 64) + " Gi",
			strconv.Itoa(len(podListOnNode)) + "/" + node.Status.Capacity.Pods().String(),
			shellPod,
		}
		nodeInfoList[i] = nodeInfo
	}
	sort.Slice(nodeInfoList, func(i, j int) bool {
		return nodeInfoList[i].Role < nodeInfoList[j].Role
	})
	nodeInfoStr := table.Table(nodeInfoList)
	fmt.Println(nodeInfoStr)
}

func NewCmdNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "node",
		Short:                 "打印更加细致的集群信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunNode(cmd, args)
		},
	}
	return cmd
}
