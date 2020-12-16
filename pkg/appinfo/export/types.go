package export

type ClusterInfo struct {
	AppName     string
	ClusterName string
	MaxReplica  int
	MemRequest  string
	MemLimit    string
	AppUrl      string
}
