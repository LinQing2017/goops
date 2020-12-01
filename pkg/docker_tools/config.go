package docker_tools

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/spf13/pflag"
)

var (
	templateFile         string
	outputFile           string
	reserveConnectorArgs = mapset.NewSet("port", "protocol", "redirectPort", "connectionTimeout")
)

func AddTomcatSetter(flags *pflag.FlagSet) {
	flags.StringVar(&templateFile, "template", "/usr/local/tomcat/conf/server-template.xml", "配置模板文件。")
	flags.StringVar(&outputFile, "output", "/usr/local/tomcat/conf/server.xml", "配置输出文件。")
}
