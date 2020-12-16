package export

type ClusterInfo struct {
	AppName     string
	ClusterName string
	Replica     int
	MaxReplica  int
	MemRequest  string
	MemLimit    string
	AppUrl      string
}
