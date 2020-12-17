package export

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/util/table"
)

func NewCmdExport() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "export",
		Short:                 "导出信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	AddExportFlag(cmd.Flags())
	return cmd
}

func Main(cmd *cobra.Command, args []string) {
	db_tools.InitDBClient()
	if envType == 0 {
		logrus.Error("请指定环境变量和APP名称")
		return
	}

	switch outputFormat {

	case "k8s-cluster":
		// 导出K8S集群信息
		output := ExportK8sCluster()
		table.Output(output)
	default:
		logrus.Error("请指定正确的导出类型")
	}
	db_tools.CloseAllDBClient()
}
