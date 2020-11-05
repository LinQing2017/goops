package check

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/util"
	"strconv"
	"strings"
)

func addDiskFlag(flags *pflag.FlagSet) {

}

func NewCmdDisk() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disk",
		Short:                 "检查Docker容器的磁盘使用情况。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunDisk(cmd, args)
		},
	}
	addDiskFlag(cmd.Flags())
	return cmd
}

func RunDisk(cmd *cobra.Command, args []string) {
	//cli := util.DockerClient("tcp://172.24.135.47:4243")
	////containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	//if err != nil {
	//	panic(err)
	//}

	kubeClientSet, _ := util.KubeClient(cmd)

	allPods, _ := util.GetPodDict(kubeClientSet, "")

	pods := allPods["172.24.135.10"]

	for _, pod := range pods {

		for _, container := range pod.Status.ContainerStatuses {

			containerDataPath := strings.Replace(container.ContainerID, "docker://", "/data/var/lib/docker/containers/", -1)
			size, _ := util.CalculateDirSize(containerDataPath)
			fmt.Printf("%s : %s mb\n", containerDataPath, strconv.Itoa(int(size)))
		}

	}
}
