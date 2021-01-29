package ops

import (
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
		pclient.SwitchDomain(args[1], args[2], args[3], envType)
	default:
		logrus.Error("操作不支持")
	}

	db_tools.CloseAllDBClient()
}

func NewOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ops",
		Short:                 "相关运维操作",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	ews.AddEWSFlag(cmd.Flags())
	pclient.AddProtalFlag(cmd.Flags())
	AddOpsFlag(cmd.Flags())
	return cmd
}
