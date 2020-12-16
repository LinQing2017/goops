package export

import (
	"github.com/spf13/pflag"
)

var (
	envType      int
	outputFormat string
	withUrl      bool
)

func AddExportFlag(flags *pflag.FlagSet) {

	flags.IntVarP(&envType, "env", "e", 0, "环境ID。")
	flags.StringVarP(&outputFormat, "output-formate", "o", "k8s-cluster", "输出类型。")
	flags.BoolVar(&withUrl, "url", false, "打印连接地址。")
}
