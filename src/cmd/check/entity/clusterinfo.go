package entity

import "go.mongodb.org/mongo-driver/bson"

type ClusterInfo struct {
	ClusterId       string
	AppId           string
	AppName         string
	EnvType         string
	DeploymentId    string
	Area            string
	Cluster         string
	Environment     string
	KongServiceName []string
	ServiceNodePort int64
	Domains         []string
	Modules         []string
	Targets         []string
	CreateTime      int64
	UpdateTime      int64
}

func ParseClusterInfo(b *bson.Raw) ClusterInfo {

	clusterInfo := ClusterInfo{
		//ClusterId:       b.Lookup("cluster_id").StringValue(),
		AppId:   b.Lookup("app_id").StringValue(),
		AppName: b.Lookup("app_name").StringValue(),
		//EnvType:         b.Lookup("env_type").StringValue(),
		//DeploymentId:    b.Lookup("deployment_id").StringValue(),
		//Area:            b.Lookup("area").StringValue(),
		//Cluster:         b.Lookup("cluster").StringValue(),
		//Environment:     b.Lookup("environment").StringValue(),
		//KongServiceName: make([]string, 0),
		//ServiceNodePort: 0,
		//Domains:         make([]string, 0),
		//Modules:         make([]string, 0),
		//Targets:         make([]string, 0),
	}

	//kongServiceNameValues, _ := b.Lookup("kong_service_name").Array().Values()
	//for _, v := range kongServiceNameValues {
	//	clusterInfo.KongServiceName = append(clusterInfo.KongServiceName, v.String())
	//}
	//
	//domainsValues, _ := b.Lookup("domains").Array().Values()
	//for _, v := range domainsValues {
	//	clusterInfo.Domains = append(clusterInfo.Domains, v.String())
	//}
	//modulesValues, _ := b.Lookup("modules").Array().Values()
	//for _, v := range modulesValues {
	//	clusterInfo.Modules = append(clusterInfo.Modules, v.String())
	//}
	//targetsValues, _ := b.Lookup("targets").Array().Values()
	//for _, v := range targetsValues {
	//	clusterInfo.Modules = append(clusterInfo.Modules, v.String())
	//}
	return clusterInfo
}
