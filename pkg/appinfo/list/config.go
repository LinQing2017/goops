package list

import (
	"github.com/spf13/pflag"
)

var (
	nameFromFile  string
	envType       int
	outputFormate string
	kubeConfig    string
)

func AddListFlag(flags *pflag.FlagSet) {
	flags.StringVar(&nameFromFile, "name-from-file", "apps-list", "从文件中获取APP信息。")
	flags.IntVar(&envType, "env", 0, "环境ID。")
	flags.StringVarP(&outputFormate, "output-formate", "o", "default", "输出类型。")
	flags.StringVar(&kubeConfig, "kubeconfig", "", "源Kubernete集群的config配置文件。")
}
