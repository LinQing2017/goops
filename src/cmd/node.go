package cmd

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
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
	Label  string
	CPU    string
	Memory string
	Pod    string
	Shell  string
}

var labelFilter = mapset.NewSet(
	"beta.kubernetes.io/arch",
	"beta.kubernetes.io/os",
	"kubernetes.io/arch",
	"kubernetes.io/hostname",
	"kubernetes.io/os",
)

func RunNode(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := util.KubeClient(cmd)

	nodes, err := kubeClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeInfoList := make([]NodeInfo, len(nodes.Items))
	shellPods, err := util.GetShellPodList(kubeClientSet)
	allPodDist, _ := util.GetPodDict(kubeClientSet, "")
	for i, node := range nodes.Items {
		// 获取Role以及Label信息
		role := ""
		envLabel := make([]string, 0)
		typeLabel := make([]string, 0)
		commonLabel := make([]string, 0)
		for k, v := range node.Labels {
			if strings.EqualFold(v, "env") {
				envLabel = append(envLabel, k)
				continue
			}
			if strings.EqualFold(v, "type") {
				typeLabel = append(typeLabel, k)
				continue
			}
			if strings.HasPrefix(k, "node-role.kubernetes.io") {
				role = strings.Split(k, "/")[1]
				continue
			}
			if !labelFilter.Contains(k) {
				commonLabel = append(commonLabel, k+"="+v)
				continue
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
			strings.Join(envLabel, ","),
			strings.Join(typeLabel, ","),
			strings.Join(commonLabel, ","),
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