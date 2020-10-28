package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	kubeConfigStr     string
	ShellNamespace    string
	ShellDaemonset    string
	ShellNodeName     string
	ShellNodeNameFile string
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")

	flags.StringVar(&ShellNamespace, "shellnamespace", "node-shell", "安装Shell客户端的命名空间名。")
	flags.StringVar(&ShellDaemonset, "shelldeploy", "node-shell-tool", "安装Shell客户端的控制器名称。")

	flags.StringVar(&ShellNodeName, "node", "", "在指定宿主机节点执行操作。")
	flags.StringVar(&ShellNodeNameFile, "nodefile", "", "通过文件指定要运行命令的宿主机。")
}
