package get

import (
	"github.com/spf13/pflag"
)

var (
	envType    int
	outputType string
)

func AddGetFlag(flags *pflag.FlagSet) {
	flags.IntVarP(&envType, "env", "e", 0, "环境ID。")
	flags.StringVarP(&outputType, "output", "o", "brief", "输出格式json|table")
}
