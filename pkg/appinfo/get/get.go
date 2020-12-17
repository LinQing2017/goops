package get

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	sys_tool "goops/pkg/util/sys"
	systools "goops/pkg/util/sys"
	"os"
	"strings"
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
		case "appname":
			ConvertID2Name()
		default:
			logrus.Error("输出格式异常")
		}

	} else {
		logrus.Error("没有查询到信息")
	}

	db_tools.CloseAllDBClient()
}

func ConvertID2Name() {

	if !strings.EqualFold(intputFile, "") {
		logrus.Error("请指定输入文件")
		os.Exit(-1)
	}
	appnames := make([]string, 0)
	for _, appid := range systools.ReadLine(intputFile) {
		var app types.App
		if err := db_tools.GetOne(db_tools.PortalMongoDB, "app", bson.M{"_id": primitive.ObjectIDFromHex(appid)}, db_tools.NdpPortalClient, app); err != nil {
			logrus.Error(err.Error())
		}
		appnames = append(appnames, app.Name)
	}
	for _, name := range appnames {
		fmt.Println(name)
	}

}
