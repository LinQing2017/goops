package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	mongotools "goops/pkg/util/mongo"
	systools "goops/pkg/util/sys"
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
	portalDBURI := "mongodb://" + portalMongoUser + ":" + portalMongoPasswd + "@" + portalMongoUrl + "/" + portalMongoDB + "?autoConnectRetry=true"
	ndpPortalClient := mongotools.MongoClient(portalDBURI)

	if !strings.EqualFold(nameFromFile, "") {
		portalAppInfoList := make([]types.AppPortalInfo, 0)
		for _, appname := range systools.ReadLine(nameFromFile) {
			appname = strings.TrimSpace(appname)
			portalAppInfo := db_tools.GetPortalInfo(appname, portalMongoDB, ndpPortalClient)
			portalAppInfoList = append(portalAppInfoList, portalAppInfo)
			fmt.Println("*******************************************************")
			fmt.Println(portalAppInfo.APP.Name, portalAppInfo.APP.ID.String())
			fmt.Println("*******************************************************")
			for _, service := range portalAppInfo.K8SServiceList {
				if service.Env == 18 {
					fmt.Println("K8S集群：", service.ClusterId, service.Env, service.PreN)
				}

			}
		}
	} else {
		portalAppInfo := db_tools.GetPortalInfo(args[0], portalMongoDB, ndpPortalClient)
		PrintApplicationInfo(portalAppInfo)
	}

}

func PrintApplicationInfo(portalAppInfo types.AppPortalInfo) {

	fmt.Println("*******************************************************")
	fmt.Println("应用名称	：" + portalAppInfo.APP.Name)
	fmt.Println("应用ID		：" + portalAppInfo.APP.ID.Hex())

	for _, service := range portalAppInfo.K8SServiceList {
		fmt.Println("K8S集群：", service.ClusterId, service.Env, service.PreN)
	}
	//for _, service := range portalAppInfo.EWSServiceList {
	//	fmt.Println("EWS集群：", service.ClusterId, service.Env, service.PreN)
	//}
}
