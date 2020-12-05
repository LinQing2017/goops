package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"goops/pkg/appinfo/db_tools"
	"goops/pkg/appinfo/db_tools/types"
	mongotools "goops/pkg/util/mongo"
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
	portalDBURI := "mongodb://" + db_tools.PortalMongoUser + ":" + db_tools.PortalMongoPasswd + "@" + db_tools.PortalMongoUrl + "/" + db_tools.PortalMongoDB + "?autoConnectRetry=true"
	ndpPortalClient := mongotools.MongoClient(portalDBURI)
	portalAppInfo := db_tools.GetPortalInfo(args[0], db_tools.PortalMongoDB, ndpPortalClient)
	PrintApplicationInfo(portalAppInfo)
	mongotools.MongoDisconnect(ndpPortalClient)

}

func PrintApplicationInfo(portalAppInfo types.AppPortalInfo) {

	fmt.Println("*******************************************************")
	fmt.Println("应用名称：", portalAppInfo.APP.Name)
	fmt.Println("应用ID：", portalAppInfo.APP.ID.Hex())
	fmt.Println("单实例：", portalAppInfo.APP.SingleInstance)
	fmt.Println("*******************************************************")
	for _, service := range portalAppInfo.K8SServiceList {
		fmt.Printf("%2d_K8S_Cluster : %s\n", service.Env, service.ClusterId)
	}
	for _, service := range portalAppInfo.EWSServiceList {
		fmt.Printf("%2d_EWS_Cluster : %s\n", service.Env, service.ClusterId)
	}
}
