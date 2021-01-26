package ews

import (
	"github.com/spf13/pflag"
	"os"
)

var (
	EWSURL   string
	EWSToken string
)

func AddEWSFlag(flags *pflag.FlagSet) {
	flags.StringVar(&EWSURL, "ews-url", os.Getenv("EWS_URL"), "生产环境弹性WEB地址")
	flags.StringVar(&EWSToken, "ews-token", os.Getenv("EWS_TOKEN"), "认证的Token地址")
}
