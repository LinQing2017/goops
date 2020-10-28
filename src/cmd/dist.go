package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"kube-tools/src/util"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	distNodename string
	distNodefile string
	timeOut      int
)

func RunDist(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		panic("请填写源路径和目标路径")
	}
	srcPath := path.Clean(args[0])
	destPath := path.Clean(args[1])
	if !path.IsAbs(destPath) {
		panic("目标文件请指定绝对路径")
	}
	destDir := path.Dir(destPath)
	cmdArr := []string{"tar", "-xmf", "-", "-C", destDir}

	kubeClientSet, kubeClientConfig := util.KubeClient(cmd)
	// 返回所有需要运行运行的Node列表
	shellPodTargets := util.GetShellPodDict(kubeClientSet)
	i := 0
	for _, v := range shellPodTargets {
		fmt.Printf("------------------------------> No." + strconv.Itoa(i) + " Dist on node: " + v.Spec.NodeName + " <------------------------------\n")

		reader, writer := io.Pipe()
		tarExecOps := util.ExecOptions{
			Command:       strings.Join(cmdArr, " "),
			ContainerName: "",
			In:            reader,
			Out:           os.Stdout,
			Err:           os.Stderr,
			Istty:         false,
			TimeOut:       timeOut,
		}

		go func() {
			defer writer.Close()
			util.MakeTar(srcPath, destDir, writer)
		}()

		util.ExecCmd(kubeClientSet, kubeClientConfig, v, tarExecOps)
		i += 1
	}
}

func addDistFlag(flags *pflag.FlagSet) {
	flags.IntVar(&timeOut, "timeout", 15, "单个文件传输超时时间。")
}

func NewCmdDist() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "dist [local file path] [host file path]",
		Short:                 "分发文件到宿主机目录",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunDist(cmd, args)
		},
	}
	addDistFlag(cmd.Flags())
	return cmd
}
