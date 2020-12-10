package types

type RMDomains struct {
	ID          string           `bson:"_id"`
	OwnerId     string           `bson:"owner_id"`
	Area        string           `bson:"area"`
	Environment string           `bson:"environment"`
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
