package get

import (
	"github.com/spf13/pflag"
)

var (
	envType int
)

func AddGetFlag(flags *pflag.FlagSet) {
	flags.IntVar(&envType, "env", 0, "环境ID。")
}
