package eureka

import "github.com/spf13/pflag"

var (
	eurekaNetFilter  string
	eurekaPortFilter string
	outputFilename   string
)

func AddEurekaFlag(flags *pflag.FlagSet) {
	flags.StringVar(&eurekaNetFilter, "net", "", "过滤对应网段。")
	flags.StringVar(&eurekaPortFilter, "port", "", "过滤对应端口。")
	flags.StringVar(&outputFilename, "file", "", "输出结果数据到Excel文件。")
}
