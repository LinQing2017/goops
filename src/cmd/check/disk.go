package check

import (
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kube-tools/src/util"
	"os"
	"path"
	"strings"
)

var (
	nodename       string
	dockerRootPath string
	logRootPath    string
	cleanLog       bool
)

func addDiskFlag(flags *pflag.FlagSet) {
	flags.StringVar(&nodename, "nodename", "", "当前节点Node名称。")
	flags.StringVar(&dockerRootPath, "docker-root-path", "/data/var/lib/docker", "docker服务数据目录。")
	flags.StringVar(&logRootPath, "log-root-path", "/data/container_logs", "自定义容器日志根目录。")
	flags.BoolVar(&cleanLog, "cleanLog", false, "是否清理容器日志。")
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

type DiskUsageInfo struct {
	PodName       string
	ContainerSize string
	DockerLog     string
	UserLog       string
}

func RunDisk(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := util.KubeClient(cmd)
	allPods, _ := util.GetPodDict(kubeClientSet, "")

	if strings.EqualFold(nodename, "") {
		nodename, _ = os.Hostname()
	}
	pods := allPods[nodename]
	if pods == nil || len(pods) == 0 {
		fmt.Printf("%s 不是Kubernetes集群的节点或者该节点没有Pod运行")
		return
	}

	diskUsageInfos := make([]DiskUsageInfo, 0)
	// 统计运行Pod磁盘使用情况
	dockerContainerBasePath := path.Join(dockerRootPath, "containers/")
	dockercli := util.DockerClient("")
	for _, pod := range pods {

		var containerSize int64
		var dockerLogSize int64

		for _, container := range pod.Status.ContainerStatuses {
			containerDataPath := path.Join(dockerContainerBasePath, strings.Replace(container.ContainerID, "docker://", "", -1))
			logSize, _ := util.CalculateDirSize(containerDataPath)
			dockerLogSize += logSize
			cSize := util.ContainerSize(container.ContainerID, dockercli)
			containerSize += cSize
		}
		logDirPath := path.Join(logRootPath, pod.Name)

		_, isExist := os.Stat(logDirPath)
		var logDirSize int64
		if isExist == nil {
			logDirSize, _ = util.CalculateDirSize(logDirPath)
		}
		info := DiskUsageInfo{
			PodName:       pod.Name,
			ContainerSize: util.FormatByte(containerSize),
			DockerLog:     util.FormatByte(logDirSize),
			UserLog:       util.FormatByte(logDirSize),
		}
		diskUsageInfos = append(diskUsageInfos, info)
	}
	nodeInfoStr := table.Table(diskUsageInfos)
	fmt.Println(nodeInfoStr)
}
