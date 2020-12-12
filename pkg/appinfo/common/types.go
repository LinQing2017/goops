package common

import (
	"fmt"
	"github.com/fatih/color"
	"goops/pkg/appinfo/db_tools/types"
	"goops/pkg/appinfo/ews_client"
	"strings"
)

type AppInformation struct {
	NAME               string
	APPID              string
	PortalInfo         types.AppPortalInfo
	URL                string
	EWSClusterInfo     []ews_client.EWSCluster
	K8SClusterInfo     []types.K8SCluster
	ClusterBindDomains map[string][]*types.RMDomains
}

func (c *AppInformation) GetK8SClusterInfo(clusterId string) *types.K8SCluster {

	for _, cluster := range c.K8SClusterInfo {
		if strings.EqualFold(cluster.ID, clusterId) {
			return &cluster
		}
	}
	return nil
}

func (c *AppInformation) GetEWSClusterInfo(clusterId string) *ews_client.EWSCluster {

	for _, cluster := range c.EWSClusterInfo {
		if strings.EqualFold(cluster.ID, clusterId) {
			return &cluster
		}
	}
	return nil
}

func (c *AppInformation) GetClusterString(clusterId string) string {

	if cluster := c.GetEWSClusterInfo(clusterId); cluster != nil {
		return fmt.Sprintf("%s EWS集群", clusterId)
	}
	if cluster := c.GetK8SClusterInfo(clusterId); cluster != nil {
		return fmt.Sprintf("%s K8S集群(%s)", clusterId, color.BlueString(cluster.Area))
	}
	return fmt.Sprintf("%s %s", clusterId, color.HiRedString("未找到到集群信息"))
}
