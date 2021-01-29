package ops

import "github.com/spf13/pflag"

var (
	envType int
	action  string
)

func AddOpsFlag(flags *pflag.FlagSet) {
	flags.IntVarP(&envType, "env", "e", 0, "环境ID。")
	flags.StringVarP(&action, "action", "a", "", "操作类型，包括：sdomain（域名切换）")
}
