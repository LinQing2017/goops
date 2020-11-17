package appinfo

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	mongoUser     string
	mongoPasswd   string
	mongoUrl      string
	mongoDB       string
	nodeFilter    string
	isPrint       bool
	kubeConfigStr string
)

func AddAppInfoFlag(flags *pflag.FlagSet) {
	flags.StringVar(&mongoUser, "mongo-user", "qa_mdb_k8s_paas_api", "SDP-K8S Mongodb 用户名。")
	flags.StringVar(&mongoPasswd, "mongo-password", "lcsX5zq9QEUW", "SDP-K8S Mongodb 密码。")
	flags.StringVar(&mongoDB, "mongo-database", "qa_mdb_k8s_paas_api", "SDP-K8S Mongodb 数据库名称。")
	flags.StringVar(&mongoUrl, "mongo-url", "m3.all.debug2.mongod3.sdp:34003,m1.all.debug2.mongod3.sdp:34003,m2.all.debug2.mongod3.sdp:34003", "SDP-K8S Mongodb 地址。")
	flags.StringVar(&nodeFilter, "node-filter", "", "根据node名称过滤信息，使用逗号分隔")
	flags.BoolVar(&isPrint, "print", false, "输出到屏幕。")
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")

}
