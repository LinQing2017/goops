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

	if (!strings.EqualFold(outputType, "appname")) && (len(args) < 1 || envType == 0) {
		logrus.Error("请指定环境变量和APP名称")
		return
	}

	db_tools.InitDBClient()

	switch outputType {
	case "raw":
		appInfo := Get(args[0])
		sys_tool.PrintJSON(*appInfo)
	case "brief":
		appInfo := Get(args[0])
		sys_tool.PrintJSON(*appInfo.GetBerif())
	case "appname":
		ConvertID2Name()
	default:
		logrus.Error("输出格式异常")
	}

	db_tools.CloseAllDBClient()
}

func Get(appname string) *common.AppInformation {

	if appInfo, err := common.GetAppInfo(appname, envType); err == nil {
		return &appInfo
	} else {
		logrus.Error("没有查询到信息")
		os.Exit(-1)
	}
	return nil
}

func ConvertID2Name() {

	if strings.EqualFold(intputFile, "") {
		logrus.Error("请指定输入文件")
		os.Exit(-1)
	}
	appnames := make([]string, 0)
	for _, appid := range systools.ReadLine(intputFile) {
		var app types.App
		objectId, _ := primitive.ObjectIDFromHex(strings.TrimSpace(appid))
		if err := db_tools.GetOne(db_tools.PortalMongoDB, "app", bson.M{"_id": objectId}, db_tools.NdpPortalClient, &app); err != nil {
			logrus.Error(objectId.Hex(), err.Error())
		}
		if !strings.EqualFold(app.Name, "") {
			appnames = append(appnames, app.Name)
		}
	}
	for _, name := range appnames {
		fmt.Println(name)
	}

}
