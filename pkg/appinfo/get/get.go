package get

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	sys_tool "goops/pkg/util/sys"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get",
		Short:                 "获取应用基本信息",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			Main(cmd, args)
		},
	}
	AddGetFlag(cmd.Flags())
	return cmd
}

func Main(cmd *cobra.Command, args []string) {

	if len(args) < 1 || envType == 0 {
		logrus.Error("请指定环境变量和APP名称")
		return
	}

	db_tools.InitDBClient()

	if appInfo, err := common.GetAppInfo(args[0], envType); err == nil {
		switch outputType {
		case "raw":
			sys_tool.PrintJSON(appInfo)
		case "brief":
			sys_tool.PrintJSON(appInfo.GetBerif())
		default:
			logrus.Error("输出格式异常")
		}

	} else {
		logrus.Error("没有查询到信息")
	}

	db_tools.CloseAllDBClient()
}
