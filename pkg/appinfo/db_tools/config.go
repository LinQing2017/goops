package db_tools

import (
	"github.com/spf13/pflag"
	"os"
)

var (
	PortalMongoUser    string
	PortalMongoPasswd  string
	PortalMongoUrl     string
	PortalMongoDB      string
	K8SPaasMongoUser   string
	K8SPaasMongoPasswd string
	K8SPaasMongoUrl    string
	K8SPaasMongoDB     string
)

func AddDBFlag(flags *pflag.FlagSet) {
	flags.StringVar(&PortalMongoUser, "portal-db-user", os.Getenv("PORTAL_MONGODB_USER"), "Mongodb 用户名。")
	flags.StringVar(&PortalMongoPasswd, "portal-db-password", os.Getenv("PORTAL_MONGODB_PASSWORD"), "Mongodb 密码。")
	flags.StringVar(&PortalMongoUrl, "portal-db-url", os.Getenv("PORTAL_MONGODB_URL"), "Mongodb 地址。")
	flags.StringVar(&PortalMongoDB, "portal-db-database", os.Getenv("PORTAL_MONGODB_DATABASE"), "Mongodb 数据库名称。")

	flags.StringVar(&K8SPaasMongoUser, "k8spaas-db-user", os.Getenv("K8S_MONGODB_USER"), "Mongodb 用户名。")
	flags.StringVar(&K8SPaasMongoPasswd, "k8spaas-db-password", os.Getenv("K8S_MONGODB_PASSWORD"), "Mongodb 密码。")
	flags.StringVar(&K8SPaasMongoUrl, "k8spaas-db-url", os.Getenv("K8S_MONGODB_URL"), "Mongodb 地址。")
	flags.StringVar(&K8SPaasMongoDB, "k8spaas-db-database", os.Getenv("K8S_MONGODB_DATABASE"), "Mongodb 数据库名称。")
}
