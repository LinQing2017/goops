package entity

type ClusterInfo struct {
	Id         string `bson:"_id"`
	AppId      string `bson:"app_id"`
	AppName    string `bson:"app_name"`
	Type       string `bson:"type"`
	Group      string `bson:"group"`
	MaxReplica int    `bson:"max_replica"`
	Replica    int    `bson:"replica"`
}
