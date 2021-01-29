package types

import (
	"fmt"
	"strings"
)

type RMDomains struct {
	ID          string           `bson:"_id"`
	OwnerId     string           `bson:"owner_id"`
	Area        string           `bson:"area"`
	Environment string           `bson:"environment"`
	Component   string           `bson:"component"`
	Cluster     string           `bson:"cluster"`
	Domain      string           `bson:"domain"`
	Protocols   []string         `bson:"protocols"`
	ClusterId   string           `bson:"cluster_id"`
	Labels      []RMDomainsLabel `bson:"labels"`
	CreateTime  int              `bson:"create_time"`
}

type RMDomainsLabel struct {
	Key   string `bson:"k"`
	Value string `bson:"v"`
}

// 该域名是否是测试域名
func (c *RMDomains) IsTestDomain() bool {
	isTestDomain := false
	for _, label := range c.Labels {
		if strings.EqualFold(label.Key, "role") && strings.EqualFold(label.Value, "test") {

			isTestDomain = true
		}
	}
	return isTestDomain
}

// 改域名是否是组件域名
func (c *RMDomains) IsCmptDomain() bool {
	return strings.EqualFold(c.GetLabel("type"), "component")
}

func (c *RMDomains) GetLabel(key string) string {
	for _, label := range c.Labels {
		if strings.EqualFold(label.Key, key) {
			return label.Value
		}
	}
	return ""
}

func (c *RMDomains) ShortClusterId() string {

	var serviceType string
	for _, label := range c.Labels {
		if strings.EqualFold(label.Key, "service_type") {
			serviceType = label.Value
			break
		}
	}

	return fmt.Sprintf("%s-%s", serviceType, c.ClusterId[len(c.ClusterId)-6:])

}
