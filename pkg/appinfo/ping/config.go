package ping

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/pflag"
)

var (
	envType      int
	nameFromFile string
	batch        bool
	clusterType  string
	codeList     string
	codeAccept   = mapset.NewSet(200, 401, 403)
)

func AddPingFlag(flags *pflag.FlagSet) {
	flags.IntVar(&envType, "env", 0, "环境ID。")
	flags.StringVar(&nameFromFile, "name-from-file", "apps-list", "从文件中获取APP信息。")
	flags.BoolVar(&batch, "batch", false, "批量拨测所有应用")
	flags.StringVarP(&clusterType, "cluster-type", "c", "k8s", "拨测集群的类型")
	flags.StringVar(&codeList, "code-list", "", "认可的WAF拨测返回值，默认包括：200,401,403")
}
