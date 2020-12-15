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
	timeOut      int
)

func AddPingFlag(flags *pflag.FlagSet) {
	flags.IntVarP(&envType, "env", "e", 0, "环境ID。")
	flags.StringVarP(&nameFromFile, "name-from-file", "f", "apps-list", "从文件中获取APP信息。")
	flags.BoolVar(&batch, "batch", true, "批量拨测所有应用")
	flags.StringVarP(&clusterType, "cluster-type", "c", "all", "拨测集群的类型：all|k8s|ews")
	flags.StringVar(&codeList, "code-list", "", "认可的WAF拨测返回值，默认包括：200,401,403")
	flags.IntVar(&timeOut, "ping-timeout", 3, "拨测超时时间")
}
