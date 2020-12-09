package common

import (
	"goops/pkg/appinfo/db_tools/types"
	"goops/pkg/appinfo/ews_client"
)

type AppInformation struct {
	NAME           string
	APPID          string
	PortalInfo     types.AppPortalInfo
	URL            string
	EWSClusterInfo []ews_client.EWSCluster
	K8SClusterInfo []types.K8SCluster
}
