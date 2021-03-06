package ews

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

type ClusterConfigs struct {
	Cpu                     int    `json:"cpu"`
	Memory                  int    `json:"memory"`
	Size                    int    `json:"size"`
	jdkVersion              string `json:"jdk_version"`
	TomcatVersion           string `json:"tomcat_version"`
	Pinpoint                bool   `json:"pinpoint"`
	PinpointVersion         string `json:"pinpoint_version"`
	PinpointGroupName       string `json:"pinpoint_group_name"`
	Dump                    bool   `json:"dump"`
	GcLog                   bool   `json:"gc_log"`
	MaxThread               int    `json:"max_thread"`
	DeployMode              string `json:"deploy_mode"`
	AutoSmoothDeploy        bool   `json:"auto_smooth_deploy"`
	RollbackWhenDeployAbort bool   `json:"rollback_when_deploy_abort"`
	HealthCheck             string `json:"health_check"`
	DisallowSpider          bool   `json:"disallow_spider"`
	Color                   string `json:"color"`
	SensitiveWord           bool   `json:"sensitive_word"`
	GcAlgorithm             string `json:"gc_algorithm"`
}

type Instances struct {
	ID         string `json:"id"`
	Server     string `json:"server"`
	ServerIp   string `json:"server_ip"`
	Port       int    `json:"port"`
	PackageUrl string `json:"package_url"`
	Version    string `json:"version"`
}

type Domains struct {
	Domain   string `json:"domain"`
	Protocol string `json:"https"`
}

type EWSCluster struct {
	ID           string         `json:"id"`
	RollBackable bool           `json:"rollbackable"`
	Configs      ClusterConfigs `json:"configs"`
	Instances    []Instances    `json:"instances"`
	MajorDomains []Domains      `json"major_domains"`
}

func (c *EWSCluster) ShortClusterName() string {
	if (!strings.EqualFold(c.ID, "")) && len(c.ID) > 6 {
		return color.HiGreenString("%s-%s", "ews", c.ID[len(c.ID)-6:])
	} else {
		return color.HiRedString("ERROR")
	}
}

func (c *EWSCluster) ShortClusterNameNoColor() string {
	if (!strings.EqualFold(c.ID, "")) && len(c.ID) > 6 {
		return fmt.Sprintf("%s-%s", "ews", c.ID[len(c.ID)-6:])
	} else {
		return ""
	}
}

func (c *EWSCluster) IsPackageNotFound() string {
	if len(c.Instances) > 0 && !strings.EqualFold(c.Instances[0].PackageUrl, "") {
		return ""
	} else {
		return color.HiRedString("Not Found")
	}
}

func (c *EWSCluster) GetPackageVersion() (string, string) {
	if len(c.Instances) > 0 && !strings.EqualFold(c.Instances[0].PackageUrl, "") {
		return c.Instances[0].PackageUrl, c.Instances[0].Version
	} else {
		return color.HiRedString("Not Found"), ""
	}
}
