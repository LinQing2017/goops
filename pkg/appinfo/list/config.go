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
	flags.StringVarP(&nameFromFile, "name-from-file", "f", "apps-list", "从文件中获取APP信息。")
	flags.IntVarP(&envType, "env", "e", 0, "环境ID。")
	flags.StringVarP(&outputFormate, "output-formate", "o", "app", "输出类型: app|url|cluster。")
}
