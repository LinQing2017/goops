package list

import (
	"github.com/spf13/pflag"
)

var (
	nameFromFile string
	envType      int
	envTypeList  map[int]string
)

func AddListFlag(flags *pflag.FlagSet) {
	flags.StringVar(&nameFromFile, "name-from-file", "apps-list", "从文件中获取APP信息。")
	flags.IntVar(&envType, "env", 0, "环境ID。")
}
