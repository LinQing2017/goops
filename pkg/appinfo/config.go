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
	metricFile    string
)

func AddAppInfoFlag(flags *pflag.FlagSet) {
	flags.StringVar(&mongoUser, "mongo-user", "", "Mongodb 用户名。")
	flags.StringVar(&mongoPasswd, "mongo-password", "", "Mongodb 密码。")
	flags.StringVar(&mongoDB, "mongo-database", "ndp_portal", "Mongodb 数据库名称。")
	//flags.StringVar(&mongoUrl, "mongo-url", "m3.ndp.prod.mongod3.sdp:34002,m2.ndp.prod.mongod3.sdp:34002,m1.ndp.prod.mongod3.sdp:34002", "Mongodb 地址。")
	flags.StringVar(&mongoUrl, "mongo-url", "", "Mongodb 地址。")
	flags.StringVar(&metricFile, "metrics-file", "", "流量统计execel文件。")
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Kubernete集群的config配置文件。")
}
