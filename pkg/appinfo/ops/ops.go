package ops

import (
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/client/ews"
	pclient "goops/pkg/appinfo/client/portal"
	"goops/pkg/appinfo/db_tools"
)

func Main(cmd *cobra.Command, args []string) {

	db_tools.InitDBClient()

	switch action {
	case "sdomain":
		pclient.SwitchDomain(args[0], args[1], args[2], envType)
	default:
		logrus.Error("操作不支持")
	}

	db_tools.CloseAllDBClient()
}

var Example = dedent.Dedent(`
		# 切换集群域名
		sdp-app-cli -e <环境编号> ops -a sdomain <颗粒名称>  <域名>  <集群名称>

		# 编辑Kong Service后端
		开发中~~~！！`)

func NewOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ops",
		Short:                 "相关运维操作",
		DisableFlagsInUseLine: true,
		Example:               Example,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	ews.AddEWSFlag(cmd.Flags())
	pclient.AddProtalFlag(cmd.Flags())
	AddOpsFlag(cmd.Flags())
	return cmd
}
