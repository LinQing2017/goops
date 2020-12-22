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

type BerifDomains struct {
	Domain    string
	Protocols string
	Bind      string
}

type BerifAppInformation struct {
	NAME       string
	APPID      string
	Single     bool
	URL        string
	ClusterIDs []string
	Domains    []BerifDomains
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
		return fmt.Sprintf("%s EWS集群(%s)", clusterId, color.YellowString("ews"))
	}
	if cluster := c.GetK8SClusterInfo(clusterId); cluster != nil {
		return fmt.Sprintf("%s K8S集群(%s)", clusterId, color.BlueString(cluster.Area))
	}
	return fmt.Sprintf("%s %s", clusterId, color.HiRedString("未找到到集群信息"))
}

// 返回应用的简略信息
func (c *AppInformation) GetBerif() *BerifAppInformation {

	berif := BerifAppInformation{
		NAME:       c.NAME,
		APPID:      c.APPID,
		Single:     c.PortalInfo.APP.SingleInstance,
		URL:        c.URL,
		ClusterIDs: make([]string, 0),
		Domains:    make([]BerifDomains, 0),
	}

	for _, cluster := range c.K8SClusterInfo {
		berif.ClusterIDs = append(berif.ClusterIDs, fmt.Sprintf("%s(%s)", cluster.ID, cluster.Area))
	}
	for _, cluster := range c.EWSClusterInfo {
		berif.ClusterIDs = append(berif.ClusterIDs, fmt.Sprintf("%s(%s)", cluster.ID, "ews"))
	}

	for _, v := range c.ClusterBindDomains {
		for _, b := range v {
			berifDomain := BerifDomains{
				Domain:    b.Domain,
				Protocols: strings.Join(b.Protocols, ","),
				Bind:      b.ShortClusterId(),
			}
			berif.Domains = append(berif.Domains, berifDomain)
		}

	}
	return &berif

}

// 弹性实例是否找不到包
func (c *AppInformation) IsPackageNotFound() bool {

	packageNum := 0 // 能够找到包路径的弹性web集群数目
	for _, ewsCluster := range c.EWSClusterInfo {
		if len(ewsCluster.Instances) > 0 && !strings.EqualFold(ewsCluster.Instances[0].PackageUrl, "") {
			packageNum++
		}
	}
	// 所有集群都能找到一个对应实例，有包地址。否则认为应用缺少war包
	return packageNum < len(c.PortalInfo.EWSServiceList)
}
