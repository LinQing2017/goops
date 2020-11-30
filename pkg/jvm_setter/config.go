package jvm_setter

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/pflag"
)

var (
	templateFile         string
	outputFile           string
	catalinaBase         string
	reserveConnectorArgs = mapset.NewSet("port", "protocol", "redirectPort", "connectionTimeout")
)

func AddJVMSetter(flags *pflag.FlagSet) {
	flags.StringVar(&catalinaBase, "catalina-base", "/usr/local/tomcat", "tomcat默认配置目录。")
	flags.StringVar(&templateFile, "template", "server-template.xml", "配置模板文件。")
	flags.StringVar(&outputFile, "output", "server.xml", "配置输出文件。")
}
