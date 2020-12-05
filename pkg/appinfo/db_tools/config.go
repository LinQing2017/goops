package db_tools

import (
	"github.com/spf13/pflag"
	"os"
)

var (
	PortalMongoUser   string
	PortalMongoPasswd string
	PortalMongoUrl    string
	PortalMongoDB     string
)

func AddDBFlag(flags *pflag.FlagSet) {
	flags.StringVar(&PortalMongoUser, "mongo-user", os.Getenv("PORTAL_MONGODB_USER"), "Mongodb 用户名。")
	flags.StringVar(&PortalMongoPasswd, "mongo-password", os.Getenv("PORTAL_MONGODB_PASSWORD"), "Mongodb 密码。")
	flags.StringVar(&PortalMongoDB, "mongo-database", os.Getenv("PORTAL_MONGODB_DATABASE"), "Mongodb 数据库名称。")
	flags.StringVar(&PortalMongoUrl, "mongo-url", os.Getenv("PORTAL_MONGODB_URL"), "Mongodb 地址。")
}
