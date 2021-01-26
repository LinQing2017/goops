package portal

import (
	"github.com/spf13/pflag"
	"os"
)

var (
	UcProxyURL string
	ServerUrl  string
	XLoginName string
	XPassword  string
)

func AddProtalFlag(flags *pflag.FlagSet) {
	flags.StringVar(&UcProxyURL, "uc-proxy-url", os.Getenv("UC_PROXY_URL"), "")
	flags.StringVar(&ServerUrl, "server-url", os.Getenv("SERVER_URL"), "")
	flags.StringVar(&XLoginName, "x-login-name", os.Getenv("XLOGINNAME"), "")
	flags.StringVar(&XPassword, "x-password", os.Getenv("XPASSWORD"), "")

}
