package types

import "fmt"

type K8SCluster struct {
	ID         string `bson:"_id"`
	EnvType    string `bson:"env_type"`
	Area       string `bson:"area"`
	AppId      string `bson:"app_id"`
	AppName    string `bson:"app_name"`
	MaxReplica int    `bson:"max_replica"`
	Replica    int    `bson:"replica"`

	Config K8SClusterConfig `bson:"config"`
}

type K8SClusterConfig struct {
	NodeSelector      map[string]string    `bson:"node_selector"`
	NodeType          string               `bson:"node_type"`
	NodeName          string               `bson:"node_name"`
	ContainerTemplate K8SContainerTemplate `bson:"container_template"`
	Hpa               HPAConfig            `bson:"hpa"`
}

type K8SContainerTemplate struct {
	Resources K8SResources `bson:"resources"`
}

type HPAConfig struct {
	Enable     bool `bson:"enable"`
	MinReplica int  `bson:"minReplica"`
	MaxReplica int  `bson:"maxReplica"`
}

type K8SResources struct {
	Limits   K8SCal `bson:"limits"`
	Requests K8SCal `bson:"requests"`
}

type K8SCal struct {
	Cpu    map[string]string `bson:"cpu"`
	Memory map[string]string `bson:"memory"`
}

func (c *K8SCluster) ShortClusterName() string {
	return fmt.Sprintf("%s-%s", "ews", c.ID[len(c.ID)-6:])
}
