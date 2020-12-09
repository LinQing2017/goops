package get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/common"
	"goops/pkg/appinfo/db_tools"
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

	for _, appname := range args {
		appInfo, err := common.GetAppInfo(appname, envType)
		if err != nil {
			continue
		}
		if jsonByte, err := json.Marshal(appInfo); err == nil {
			var jsonFormate bytes.Buffer
			if err = json.Indent(&jsonFormate, jsonByte, "", "    "); err != nil {
				logrus.Error("解析json字符传异常")
			} else {
				fmt.Println(jsonFormate.String())
			}
		} else {
			logrus.Error("解析json字符传异常")
		}
	}
	db_tools.CloseAllDBClient()
}
