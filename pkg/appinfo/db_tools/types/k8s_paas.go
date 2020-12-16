package types

type K8SCluster struct {
	ID         string `bson:"_id"`
	EnvType    string `bson:"env_type"`
	Area       string `bson:"area"`
	AppId      string `bson:"app_id"`
	AppName    string `bson:"app_name"`
	MaxReplica int    `bson:"max_replica"`

	Config K8SClusterConfig `bson:"config"`
}

type K8SClusterConfig struct {
	NodeSelector      map[string]string    `bson:"node_selector"`
	NodeType          string               `bson:"node_type"`
	NodeName          string               `bson:"node_name"`
	ContainerTemplate K8SContainerTemplate `bson:"container_template"`
}

type K8SContainerTemplate struct {
	Resources K8SResources `bson:"resources"`
}

type K8SResources struct {
	Limits   K8SCal `bson:"limits"`
	Requests K8SCal `bson:"requests"`
}

type K8SCal struct {
	Cpu    map[string]string `bson:"cpu"`
	Memory map[string]string `bson:"memory"`
}
