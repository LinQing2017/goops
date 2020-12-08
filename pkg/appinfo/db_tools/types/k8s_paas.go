package types

type K8SCluster struct {
	ID      string           `bson:"_id"`
	EnvType string           `bson:"env_type"`
	Area    string           `bson:"area"`
	AppId   string           `bson:"app_id"`
	AppName string           `bson:"app_name"`
	Config  K8SClusterConfig `bson:"config"`
}

type K8SClusterConfig struct {
	NodeSelector map[string]string `bson:"node_selector"`
	NodeType     string            `bson:"node_type"`
	NodeName     string            `bson:"node_name"`
}
