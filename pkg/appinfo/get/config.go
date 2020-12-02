package get

import (
	"github.com/spf13/pflag"
	"os"
)

var (
	portalMongoUser   string
	portalMongoPasswd string
	portalMongoUrl    string
	portalMongoDB     string
	nameFromFile      string
)

func AddGetFlag(flags *pflag.FlagSet) {
	flags.StringVar(&portalMongoUser, "mongo-user", os.Getenv("PORTAL_MONGODB_USER"), "Mongodb 用户名。")
	flags.StringVar(&portalMongoPasswd, "mongo-password", os.Getenv("PORTAL_MONGODB_PASSWORD"), "Mongodb 密码。")
	flags.StringVar(&portalMongoDB, "mongo-database", os.Getenv("PORTAL_MONGODB_DATABASE"), "Mongodb 数据库名称。")
	flags.StringVar(&portalMongoUrl, "mongo-url", os.Getenv("PORTAL_MONGODB_URL"), "Mongodb 地址。")
	flags.StringVar(&nameFromFile, "name-from-file", "apps-list", "从文件中获取APP信息。")

}
